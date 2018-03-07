package main

import (
	"fmt"
	"net/http"

	"gostats/controllers/v3"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No environment file found.  Please add .env file to root with required API keys")
	}
	// Instantiate gorilla mux
	r := mux.NewRouter()
	// Define routes
	r.HandleFunc("/", controllers.PostsHandler)
	// Start server with gorilla mux
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
