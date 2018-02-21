package main

import (
	"cabstats/controllers"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment file
	err := godotenv.Load()
	if err != nil {
		// TODO: Handle environment loading failure
	}

	// Instantiate gorilla mux
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/", controllers.IndexHandler)
	r.HandleFunc("/posts", controllers.PostsHandler)

	// Start server with gorilla mux
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
