package users

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hugodiazo/arq-soft-2/db"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User representa la estructura del usuario
type Users struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error al encriptar la contraseña", http.StatusInternalServerError)
		return
	}

	_, err = db.DB.Exec("INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)",
		user.Name, user.Email, string(hashedPassword), user.Role)

	if err != nil {
		log.Println("Error al registrar usuario:", err)
		http.Error(w, "Error al registrar usuario", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Usuario registrado con éxito"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	var storedPassword string
	err := db.DB.QueryRow("SELECT password FROM users WHERE email = ?", creds.Email).Scan(&storedPassword)
	if err == sql.ErrNoRows {
		http.Error(w, "Usuario no encontrado", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password)); err != nil {
		http.Error(w, "Contraseña incorrecta", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": creds.Email,
		"exp":   expirationTime.Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error al generar token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	json.NewEncoder(w).Encode(map[string]string{"message": "Inicio de sesión exitoso"})
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.DB.Query("SELECT id, name, email, role FROM users")
	if err != nil {
		log.Println("Error al obtener usuarios:", err)
		http.Error(w, "Error al obtener usuarios", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role); err != nil {
			log.Println("Error al escanear usuario:", err)
			continue
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		json.NewEncoder(w).Encode(map[string]string{"message": "No se encontraron usuarios"})
		return
	}

	json.NewEncoder(w).Encode(users)
}

// GetUserByID maneja la obtención de un usuario por ID
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extraer el ID del usuario desde la URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "ID de usuario no proporcionado", http.StatusBadRequest)
		return
	}
	userID := parts[2]

	var user User
	err := db.DB.QueryRow("SELECT id, name, email, role FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Role)
	if err == sql.ErrNoRows {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("Error al obtener usuario:", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// UpdateUser maneja la actualización de un usuario
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	_, err := db.DB.Exec("UPDATE users SET name = ?, email = ?, role = ? WHERE id = ?",
		user.Name, user.Email, user.Role, user.ID)
	if err != nil {
		log.Println("Error al actualizar usuario:", err)
		http.Error(w, "Error al actualizar usuario", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Usuario actualizado con éxito"})
}
