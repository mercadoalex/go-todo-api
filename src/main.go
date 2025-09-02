package main

import (
	"log"      // For logging server status and errors
	"net/http" // Provides HTTP server and client implementations

	"github.com/gorilla/mux" // Third-party package for advanced HTTP routing
)

func main() {
	// Create a new router using gorilla/mux, which helps manage URL routes
	router := mux.NewRouter()

	// Define the API endpoints and associate them with handler functions
	// POST /tasks: Calls AddTask to create a new task
	router.HandleFunc("/tasks", AddTask).Methods("POST")
	// GET /tasks: Calls GetTasks to retrieve all tasks
	router.HandleFunc("/tasks", GetTasks).Methods("GET")
	// PUT /tasks/{id}: Calls UpdateTask to modify a specific task by ID
	router.HandleFunc("/tasks/{id}", UpdateTask).Methods("PUT")
	// DELETE /tasks/{id}: Calls DeleteTask to remove a specific task by ID
	router.HandleFunc("/tasks/{id}", DeleteTask).Methods("DELETE")

	// Log a message to indicate the server is starting
	log.Println("Starting server on :8080")
	// Start the HTTP server on port 8080, using the router to handle requests
	// If the server fails to start, log the error and exit
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
