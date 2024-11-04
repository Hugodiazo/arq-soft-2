package courses

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hugodiazo/arq-soft-2/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtKey = []byte("my_secret_key")

// getUserIDFromToken extrae el ID del usuario del token JWT
func getUserIDFromToken(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	log.Println("Encabezado Authorization:", authHeader) // Depurar el encabezado

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
		log.Println("Error al parsear o token inválido:", err)
		return 0, fmt.Errorf("token inválido")
	}

	// Extraer las reclamaciones del token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return 0, fmt.Errorf("ID de usuario no encontrado en el token")
		}
		log.Println("ID de usuario extraído:", int(userID))
		return int(userID), nil
	}

	return 0, fmt.Errorf("no se pudieron obtener las reclamaciones del token")
}

func indexCourseInSolr(course Course, id string) {
	course.ID = primitive.ObjectID{} // Limpiamos el ID para evitar conflictos

	// Construimos la URL de Solr
	url := "http://localhost:8983/solr/courses/update?commit=true"
	body, _ := json.Marshal(struct {
		ID           string `json:"id"`
		Title        string `json:"title"`
		Description  string `json:"description"`
		Instructor   string `json:"instructor"`
		Duration     int    `json:"duration"`
		Level        string `json:"level"`
		Availability bool   `json:"availability"`
	}{
		ID:           id,
		Title:        course.Title,
		Description:  course.Description,
		Instructor:   course.Instructor,
		Duration:     course.Duration,
		Level:        course.Level,
		Availability: course.Availability,
	})

	resp, err := http.Post(url, "application/json", strings.NewReader(string(body)))
	if err != nil {
		log.Println("Error al indexar curso en Solr:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Error en la respuesta de Solr:", resp.Status)
	}
}

// Course representa un curso en la base de datos
type Course struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title        string             `json:"title"`
	Description  string             `json:"description"`
	Instructor   string             `json:"instructor"`
	Duration     int                `json:"duration"`
	Level        string             `json:"level"`
	Availability bool               `json:"availability"`
}

// CreateCourse maneja la creación de un curso
func CreateCourse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var course Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	// Crear un nuevo ObjectID
	course.ID = primitive.NewObjectID()

	// Insertar el curso en MongoDB
	_, err := db.MongoDB.Collection("courses").InsertOne(context.TODO(), course)
	if err != nil {
		http.Error(w, "Error al crear el curso", http.StatusInternalServerError)
		return
	}

	// Convertir el ID a string para Solr
	stringID := course.ID.Hex()
	course.ID = primitive.ObjectID{} // Limpiamos el ObjectID si es necesario para Solr

	// Indexar el curso en Solr
	indexCourseInSolr(course, stringID)

	json.NewEncoder(w).Encode(map[string]string{"message": "Curso creado con éxito"})
}

// GetCourses maneja la obtención de todos los cursos
func GetCourses(w http.ResponseWriter, r *http.Request) {
	cursor, err := db.MongoDB.Collection("courses").Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Error al obtener cursos", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var courses []Course
	for cursor.Next(context.TODO()) {
		var course Course
		if err := cursor.Decode(&course); err != nil {
			continue
		}
		courses = append(courses, course)
	}

	json.NewEncoder(w).Encode(courses)
}

// GetCourseByID maneja la obtención de un curso por ID
func GetCourseByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/courses/")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var course Course
	filter := bson.M{"_id": objectID}

	err = db.MongoDB.Collection("courses").FindOne(context.TODO(), filter).Decode(&course)
	if err != nil {
		http.Error(w, "Curso no encontrado", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(course)
}

// UpdateCourse maneja la actualización de un curso
func UpdateCourse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/courses/update/")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var course Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": course}

	result, err := db.MongoDB.Collection("courses").UpdateOne(context.TODO(), filter, update)
	if err != nil || result.MatchedCount == 0 {
		http.Error(w, "Error al actualizar curso", http.StatusInternalServerError)
		return
	}

	// Convertir el ID a string para Solr
	stringID := objectID.Hex()

	// Actualizar el curso en Solr
	indexCourseInSolr(course, stringID)

	json.NewEncoder(w).Encode(map[string]string{"message": "Curso actualizado con éxito"})
}

// Enrollment representa la inscripción de un usuario en un curso
type Enrollment struct {
	UserID   int    `json:"user_id" bson:"user_id"`
	CourseID string `json:"course_id" bson:"course_id"`
	Status   string `json:"status" bson:"status"`
}

// EnrollUser maneja la inscripción de un usuario en un curso
func EnrollUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "No se pudo obtener el ID del usuario", http.StatusUnauthorized)
		return
	}

	var enrollment Enrollment
	if err := json.NewDecoder(r.Body).Decode(&enrollment); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	enrollment.UserID = userID
	enrollment.Status = "active"

	_, err = db.MongoDB.Collection("enrollments").InsertOne(context.TODO(), enrollment)
	if err != nil {
		http.Error(w, "Error al inscribir usuario", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Usuario inscrito con éxito"})
}

// GetEnrollments obtiene todas las inscripciones
func GetEnrollments(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "No se pudo obtener el ID del usuario", http.StatusUnauthorized)
		return
	}

	// Buscar las inscripciones del usuario en la base de datos
	cursor, err := db.MongoDB.Collection("enrollments").Find(context.TODO(), bson.M{"user_id": userID})
	if err != nil {
		http.Error(w, "Error al obtener inscripciones", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var enrolledCourses []Course
	for cursor.Next(context.TODO()) {
		var enrollment Enrollment
		if err := cursor.Decode(&enrollment); err != nil {
			continue
		}

		// Convertir el course_id de string a ObjectID
		objectID, err := primitive.ObjectIDFromHex(enrollment.CourseID)
		if err != nil {
			continue
		}

		// Buscar los detalles del curso usando el ObjectID
		var course Course
		err = db.MongoDB.Collection("courses").FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&course)
		if err != nil {
			continue
		}
		enrolledCourses = append(enrolledCourses, course)
	}

	// Verificar si no se encontraron cursos inscritos
	if len(enrolledCourses) == 0 {
		json.NewEncoder(w).Encode(map[string]string{"message": "No estás inscrito en ningún curso."})
		return
	}

	json.NewEncoder(w).Encode(enrolledCourses)
}

// DeleteCourse maneja la eliminación de un curso por ID
func DeleteCourse(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/courses/delete/")

	// Convertir el ID de string a ObjectId
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Intentar eliminar el curso de la base de datos
	filter := bson.M{"_id": objectID}
	result, err := db.MongoDB.Collection("courses").DeleteOne(context.TODO(), filter)
	if err != nil || result.DeletedCount == 0 {
		http.Error(w, "Curso no encontrado o no pudo ser eliminado", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Curso eliminado con éxito"})
}

// UnenrollUser maneja la desinscripción de un usuario de un curso
func UnenrollUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extraer el ID del usuario desde el token
	userID, err := getUserIDFromToken(r)
	if err != nil {
		log.Println("Error al obtener el ID del usuario:", err)
		http.Error(w, "No se pudo obtener el ID del usuario", http.StatusUnauthorized)
		return
	}
	log.Println("ID de usuario extraído:", userID)

	// Obtener el `course_id` de los parámetros de la URL
	courseID := r.URL.Query().Get("course_id")
	if courseID == "" {
		log.Println("ID del curso no proporcionado")
		http.Error(w, "ID del curso no proporcionado", http.StatusBadRequest)
		return
	}
	log.Println("ID del curso proporcionado:", courseID)

	// Convertir `course_id` de cadena a `ObjectID`
	objectID, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		log.Println("ID del curso inválido:", err)
		http.Error(w, "ID del curso inválido", http.StatusBadRequest)
		return
	}
	log.Println("ObjectID del curso convertido correctamente:", objectID)

	// Eliminar la inscripción de la base de datos
	filter := bson.M{
		"user_id":   userID,
		"course_id": courseID, // Usa el `courseID` como una cadena si está almacenado como tal
	}
	result, err := db.MongoDB.Collection("enrollments").DeleteOne(context.TODO(), filter)
	if err != nil || result.DeletedCount == 0 {
		log.Println("Error al desinscribirse o inscripción no encontrada:", err)
		http.Error(w, "Error al desinscribirse o inscripción no encontrada", http.StatusInternalServerError)
		return
	}

	log.Println("Desinscripción exitosa para userID:", userID, "y courseID:", courseID)
	json.NewEncoder(w).Encode(map[string]string{"message": "Desinscripción exitosa"})
}
