package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
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
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// GetUserByID obtiene un usuario desde la base de datos por su ID
func GetUserByID(userID int) (User, error) {
	var user User
	err := db.DB.QueryRow("SELECT id, name, email, role FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Role)
	if err == sql.ErrNoRows {
		return user, err // Usuario no encontrado
	} else if err != nil {
		log.Println("Error al obtener usuario:", err)
		return user, err // Otro error
	}

	return user, nil
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	// Validar que todos los campos estén completos
	if user.Name == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	// Validar el formato del correo electrónico
	emailRegex := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	if !regexp.MustCompile(emailRegex).MatchString(user.Email) {
		http.Error(w, "Formato de correo electrónico inválido", http.StatusBadRequest)
		return
	}

	// Asignar un rol por defecto si no se proporciona
	if user.Role == "" {
		user.Role = "user"
	}

	// Encriptar la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error al encriptar la contraseña", http.StatusInternalServerError)
		return
	}

	// Guarda la contraseña encriptada en la base de datos
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

	var userID int
	var storedPassword, userRole string // Agrega userRole aquí para obtener el rol
	err := db.DB.QueryRow("SELECT id, password, role FROM users WHERE email = ?", creds.Email).Scan(&userID, &storedPassword, &userRole)
	if err == sql.ErrNoRows {
		http.Error(w, "Usuario no encontrado", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password)); err != nil {
		http.Error(w, "Credenciales incorrectas", http.StatusUnauthorized)
		return
	}

	// Generar el token JWT con el userID y el rol
	expirationTime := time.Now().Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID, // Incluye el user_id en las reclamaciones
		"email":   creds.Email,
		"role":    userRole, // Ahora incluye el rol del usuario
		"exp":     expirationTime.Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error al generar token", http.StatusInternalServerError)
		return
	}

	log.Println("Token generado:", tokenString) // Imprime el token para depuración

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
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

func GetUserIDFromToken(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("token no proporcionado")
	}

	// Extraer el token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parsear el token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return 0, fmt.Errorf("token inválido")
	}

	// Extraer las reclamaciones
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return 0, fmt.Errorf("ID de usuario no encontrado en el token")
		}
		return int(userID), nil
	}
	return 0, fmt.Errorf("no se pudieron obtener las reclamaciones del token")
}

// Obtener un usuario por ID desde la base de datos
func GetUserByIDFromDB(userID int) (User, error) {
	var user User
	err := db.DB.QueryRow("SELECT id, name, email, role FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Role)
	if err == sql.ErrNoRows {
		return user, fmt.Errorf("usuario no encontrado")
	} else if err != nil {
		return user, err
	}
	return user, nil
}
