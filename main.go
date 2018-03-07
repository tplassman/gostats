package main

import (
	"fmt"
	"net/http"

	"cabstats/controllers"

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
<<<<<<< HEAD
	r.HandleFunc("/", controllers.PostsHandler)
=======
	r.HandleFunc("/", controllers.IndexHandler)
	r.HandleFunc("/posts", controllers.PostsHandler)
>>>>>>> 8cb75d53f9121fb6112f68d45deeae5572ccd1ee
	// Start server with gorilla mux
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
