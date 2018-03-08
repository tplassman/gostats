package main

import (
	"net/http"

	"gostats/controllers/final"
)

func main() {
	// Define routes
	http.HandleFunc("/", controllers.PostsHandler)
	// Listen
	http.ListenAndServe(":8080", nil)
}
