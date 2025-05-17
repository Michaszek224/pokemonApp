package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbname,
	)

	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Could not open db: %v", err)
	}
	defer DB.Close()

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Could not ping db: %v", err)
	}
	log.Println("Connected to Postgresql")
}
