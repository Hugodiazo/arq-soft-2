package courses

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/hugodiazo/arq-soft-2/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

	var enrollment Enrollment
	if err := json.NewDecoder(r.Body).Decode(&enrollment); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	_, err := db.MongoDB.Collection("enrollments").InsertOne(context.TODO(), enrollment)
	if err != nil {
		http.Error(w, "Error al inscribir usuario", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Usuario inscrito con éxito"})
}

// GetEnrollments obtiene todas las inscripciones
func GetEnrollments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	cursor, err := db.MongoDB.Collection("enrollments").Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Error al obtener inscripciones", http.StatusInternalServerError)
		return
	}

	var enrollments []Enrollment
	if err = cursor.All(context.TODO(), &enrollments); err != nil {
		http.Error(w, "Error al decodificar inscripciones", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(enrollments)
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
