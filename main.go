package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"pokemon/db"
	"pokemon/models"
	"strconv"
)

var templ *template.Template

func init() {
	var err error
	templ, err = template.ParseGlob("templates/*.html")

	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}
}

func main() {
	db.InitDB()

	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("POST /add", addPokemon)
	http.HandleFunc("DELETE /delete/{id}", deletePokemon)
	http.HandleFunc("GET /edit/{id}", editPokemonForm)
	http.HandleFunc("PUT /edit/{id}", editPokemon)
	fmt.Println("Server is listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	pokemons, err := db.GetAllPokemons()
	if err != nil {
		http.Error(w, "Error fetching pokemons: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error fetching pokemons: %v", err)
		return
	}
	err = templ.ExecuteTemplate(w, "index.html", pokemons)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
		return
	}
}

func addPokemon(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received POST request to add Pokemon\n")
	pokemonName := r.FormValue("pokemon_name")
	pokemonType := r.FormValue("pokemon_type")
	pokemonLvl, err := strconv.Atoi(r.FormValue("pokemon_level"))
	if err != nil {
		http.Error(w, "Invalid Pokemon level", http.StatusBadRequest)
		return
	}

	if pokemonName == "" || pokemonType == "" || pokemonLvl <= 0 {
		http.Error(w, "Invalid input: Name, Type and Level must be provided", http.StatusBadRequest)
		log.Printf("Invalid input: Name=%s, Type=%s, Level=%d", pokemonName, pokemonType, pokemonLvl)
		return
	}

	err = db.AddPokemon(pokemonName, pokemonType, pokemonLvl)
	if err != nil {
		http.Error(w, "Error adding Pokemon", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deletePokemon(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid pokemon ID", http.StatusBadRequest)
		return
	}
	err = db.DeletePokemon(id)
	if err != nil {
		log.Printf("Error deleting pokemon with id %d: %v", id, err)
		http.Error(w, "Error deleting Pokmeon: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func editPokemonForm(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid pokemon ID", http.StatusBadRequest)
		return
	}
	pokemon, err := db.GetPokemonById(id)
	if err != nil {
		log.Printf("Error getting Pokemon with ID %d for edit: %v", id, err)
		if err == sql.ErrNoRows {
			http.Error(w, "Pokemon not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching Pokemon for edit: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = templ.ExecuteTemplate(w, "edit-pokemon-row.html", pokemon)
	if err != nil {
		http.Error(w, "Error executing edit template: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error executing edit-pokemon-row template: %v", err)
		return
	}
}

func editPokemon(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid pokemon ID", http.StatusBadRequest)
		return
	}
	pokemonName := r.FormValue("pokemon_name")
	pokemonType := r.FormValue("pokemon_type")
	pokemonLvl, err := strconv.Atoi(r.FormValue("pokemon_level"))
	if err != nil {
		http.Error(w, "Invalid Pokemon level", http.StatusBadRequest)
		return
	}
	if pokemonName == "" || pokemonType == "" || pokemonLvl <= 0 {
		http.Error(w, "Invalid input: Name, Type and Level must be provided", http.StatusBadRequest)
		log.Printf("Invalid input: Name=%s, Type=%s, Level=%d", pokemonName, pokemonType, pokemonLvl)
		return
	}
	updatedPokemon := models.Pokemon{
		ID:    id,
		Name:  pokemonName,
		Type:  pokemonType,
		Level: pokemonLvl,
	}
	err = db.UpdatePokemon(updatedPokemon)
	if err != nil {
		log.Printf("Error updating Pokemon with ID %d: %v", id, err)
		http.Error(w, "Error updating Pokemon: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = templ.ExecuteTemplate(w, "pokemon-row.html", updatedPokemon)
	if err != nil {
		http.Error(w, "Error executing success template: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error executing edit-success template: %v", err)
		return
	}
	log.Printf("Pokemon with ID %d updated successfully", id)
}
