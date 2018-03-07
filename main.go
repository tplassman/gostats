package main

import (
	"net/http"

	"gostats/controllers/v2"
)

func main() {
	// Define routes
	http.HandleFunc("/", controllers.PostsHandler)
	// Listen
	http.ListenAndServe(":8080", nil)
}
