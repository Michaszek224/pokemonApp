package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"pokemon/db"
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
	fmt.Println("Server is listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	err := templ.ExecuteTemplate(w, "index.html", nil)
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
