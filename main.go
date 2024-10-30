package main

import (
	"log"
	"net/http"

	"github.com/hugodiazo/arq-soft-2/api/courses"
	"github.com/hugodiazo/arq-soft-2/api/search"
	"github.com/hugodiazo/arq-soft-2/api/users"
	"github.com/hugodiazo/arq-soft-2/db"
)

func main() {
	// Conectar a MySQL y MongoDB
	db.ConnectDB()
	db.ConnectMongoDB()

	// Endpoints de usuarios
	http.HandleFunc("/users", users.GetAllUsers)           // GET /users para listar usuarios
	http.HandleFunc("/users/", users.GetUserByID)          // GET /users/{id}
	http.HandleFunc("/users/login", users.Login)           // POST /users/login
	http.HandleFunc("/users/register", users.RegisterUser) // POST /users/register
	http.HandleFunc("/users/update", users.UpdateUser)     // PUT /users para actualizar usuario

	// Endpoints de cursos
	http.HandleFunc("/courses", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			courses.GetCourses(w, r) // GET /courses
		case http.MethodPost:
			courses.CreateCourse(w, r) // POST /courses
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/courses/", courses.GetCourseByID)       // GET /courses/{id}
	http.HandleFunc("/courses/update/", courses.UpdateCourse) // PUT /courses/update/{id}
	http.HandleFunc("/courses/enroll", courses.EnrollUser)    // POST /courses/enroll
	http.HandleFunc("/enrollments", courses.GetEnrollments)   // GET /enrollments

	// Endpoint de búsqueda
	http.HandleFunc("/search", search.SearchCourses) // GET /search?q=<query>

	// Iniciar el servidor
	log.Println("Servidor iniciado en http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
