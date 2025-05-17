package main

import (
	"fmt"
	"net/http"
	"pokemon/db"
)

func main() {
	db.InitDB()

	http.HandleFunc("/", handlerIndex)
	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")
}
