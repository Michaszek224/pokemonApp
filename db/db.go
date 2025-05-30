package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"pokemon/models"

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

func GetAllPokemons() ([]models.Pokemon, error) {
	rows, err := DB.Query("SELECT id, name, type, level FROM pokemon ORDER BY id ASC")
	if err != nil {
		return nil, fmt.Errorf("error querying pokemon: %w", err)
	}
	defer rows.Close()

	var pokemons []models.Pokemon
	for rows.Next() {
		var p models.Pokemon
		if err := rows.Scan(&p.ID, &p.Name, &p.Type, &p.Level); err != nil {
			return nil, fmt.Errorf("error scanning pokemon row: %w", err)
		}
		pokemons = append(pokemons, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through pokemon rows: %w", err)
	}

	return pokemons, nil
}

func GetPokemonById(id int) (models.Pokemon, error) {
	var p models.Pokemon
	query := `SELECT id, name, type, level FROM POKEMON WHERE id = $1`
	row := DB.QueryRow(query, id)
	err := row.Scan(&p.ID, &p.Name, &p.Type, &p.Level)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, fmt.Errorf("pokemon with ID %d not found", id)
		}
		return p, fmt.Errorf("error getting pokemon by ID: %w", err)
	}
	return p, nil
}

func UpdatePokemon(p models.Pokemon) error {
	query := `UPDATE POKEMON SET name = $1, type = $2, level = $3 WHERE id = $4`
	_, err := DB.Exec(query, p.Name, p.Type, p.Level, p.ID)
	if err != nil {
		log.Printf("Error updating pokemon: %v", err)
		return err
	}
	log.Printf("POKEMON with id= %d updated succesfully", p.ID)
	return nil
}

func DeletePokemon(id int) error {
	query := `DELETE FROM pokemon WHERE id = $1`
	res, err := DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting pokemon: %v", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no pokemon found with id= %v to delete", id)
	}

	log.Printf("Deleted pokemon with id= %d", id)
	return nil
}
