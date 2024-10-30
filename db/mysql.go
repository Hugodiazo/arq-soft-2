package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // Driver para MySQL
)

var DB *sql.DB

// ConnectDB establece la conexión a la base de datos MySQL
func ConnectDB() {
	dsn := "root:Pirata02@tcp(127.0.0.1:3306)/arqsoft2"
	var err error

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error al conectar a MySQL:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Error al hacer ping a la base de datos:", err)
	}

	log.Println("Conexión a MySQL exitosa")
}
