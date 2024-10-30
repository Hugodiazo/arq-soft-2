package main

import (
	"log"
	"net/http"

	"github.com/hugodiazo/arq-soft-2/api/courses"
	"github.com/hugodiazo/arq-soft-2/api/users"
	"github.com/hugodiazo/arq-soft-2/db"
)

func main() {
	db.ConnectDB()
	db.ConnectMongoDB()

	// Endpoints de usuarios
	http.HandleFunc("/users", users.GetAllUsers)           // GET /users para listar usuarios
	http.HandleFunc("/users/", users.GetUserByID)          // GET /users/{id} para obtener usuario por ID
	http.HandleFunc("/users/login", users.Login)           // POST /users/login para login
	http.HandleFunc("/users/register", users.RegisterUser) // POST /users/register para registro
	http.HandleFunc("/users/update", users.UpdateUser)     // PUT /users para actualizar usuario

	// Endpoints de cursos
	// Endpoints de cursos
	http.HandleFunc("/courses", courses.GetCourses)           // GET /courses (listar cursos)
	http.HandleFunc("/courses/create", courses.CreateCourse)  // POST /courses/create (crear curso)
	http.HandleFunc("/courses/", courses.GetCourseByID)       // GET /courses/{id} (obtener curso por ID)
	http.HandleFunc("/courses/update/", courses.UpdateCourse) // PUT /courses/update/{id}
	http.HandleFunc("/courses/enroll", courses.EnrollUser)    // POST /courses/enroll (inscribir usuario)
	http.HandleFunc("/enrollments", courses.GetEnrollments)   // GET /enrollments

	log.Println("Servidor iniciado en http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
