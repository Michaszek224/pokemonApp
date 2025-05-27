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

	var err error
	DB, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatalf("Could not open db: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Could not ping db: %v", err)
	}
	log.Println("Connected to Postgresql")

	createPokemonTable := `
	CREATE TABLE IF NOT EXISTS pokemon (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		type VARCHAR(50) NOT NULL,
		level INT NOT NULL
	);`
	_, err = DB.Exec(createPokemonTable)
	if err != nil {
		log.Fatalf("Could not create pokemon table: %v", err)
	}
	log.Println("Pokemon table is ready")
}

func AddPokemon(name, pokemonType string, level int) error {
	query := `INSERT INTO pokemon (name, type, level ) VALUES ($1, $2, $3)`
	_, err := DB.Exec(query, name, pokemonType, level)
	if err != nil {
		log.Printf("Error adding pokemon: %v", err)
		return err
	}
	log.Printf("Pokemon %s added successfully", name)
	return nil
}
