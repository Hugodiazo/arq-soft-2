package main

import (
	"log"
	"net/http"

	"github.com/hugodiazo/arq-soft-2/api/courses"
	"github.com/hugodiazo/arq-soft-2/api/search"
	"github.com/hugodiazo/arq-soft-2/api/users"
	"github.com/hugodiazo/arq-soft-2/db"
)

// Middleware para habilitar CORS
func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Conexión a la base de datos
	db.ConnectDB()
	db.ConnectMongoDB()

	// Llamar a la función para indexar todos los cursos en Solr
	courses.IndexAllCoursesInSolr()

	// Crear un nuevo mux
	mux := http.NewServeMux()

	// Rutas del backend
	mux.HandleFunc("/users", users.GetAllUsers)           // GET /users
	mux.HandleFunc("/users/", users.GetUserByID)          // GET /users/{id}
	mux.HandleFunc("/users/login", users.Login)           // POST /users/login
	mux.HandleFunc("/users/register", users.RegisterUser) // POST /users/register
	mux.HandleFunc("/users/update", users.UpdateUser)     // PUT /users

	// Manejo de rutas para cursos
	mux.HandleFunc("/courses", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			courses.GetCourses(w, r)
		case http.MethodPost:
			courses.CreateCourse(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/courses/", courses.GetCourseByID)        // GET /courses/{id}
	mux.HandleFunc("/courses/update/", courses.UpdateCourse)  // PUT /courses/update/{id}
	mux.HandleFunc("/courses/enroll", courses.EnrollUser)     // POST /courses/enroll
	mux.HandleFunc("/enrollments", courses.GetEnrollments)    // GET /enrollments
	mux.HandleFunc("/search", search.SearchCourses)           // GET /search?q=<query>
	mux.HandleFunc("/courses/unenroll", courses.UnenrollUser) // DELETE /courses/Unenroll

	// Usar el middleware para habilitar CORS
	handler := enableCors(mux)

	// Iniciar el servidor
	log.Println("Servidor iniciado en http://localhost:8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
