package main

import (
	"encoding/json"
	"net/http"
	"strconv" // Import strconv for int-to-string conversion
	"sync"
)

// Global variables for in-memory storage and concurrency control
var (
	tasks  = make(map[int]Task) // Stores tasks using int IDs as keys
	nextID = 1                  // Next available unique ID for new tasks
	mu     sync.Mutex           // Mutex to ensure thread-safe access to the tasks map
)

// AddTask handles POST requests to create a new task
func AddTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	// Try to decode the JSON request body into a Task struct
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		// Error handling: If decoding fails, respond with 400 Bad Request and the error message
		HandleError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}
	// Additional error handling example: Check if the title is empty
	if task.Title == "" {
		HandleError(w, http.StatusBadRequest, "Task title cannot be empty")
		return
	}
	mu.Lock()                                    // Lock the mutex to safely modify shared data
	task.ID = strconv.Itoa(nextID)               // Convert int nextID to string for Task.ID
	tasks[nextID] = task                         // Store the new task in the map using int key
	nextID++                                     // Increment the ID counter for the next task
	mu.Unlock()                                  // Unlock the mutex after modification
	RespondWithJSON(w, http.StatusCreated, task) // Respond with the created task as JSON
}

// GetTasks handles GET requests to retrieve all tasks
func GetTasks(w http.ResponseWriter, r *http.Request) {
	mu.Lock()                                // Lock the mutex to safely read shared data
	defer mu.Unlock()                        // Ensure the mutex is unlocked after the function returns
	RespondWithJSON(w, http.StatusOK, tasks) // Respond with all tasks as JSON
}

// UpdateTask handles PUT requests to update an existing task
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		HandleError(w, http.StatusBadRequest, err.Error()) // Respond with 400 Bad Request if decoding fails
		return
	}
	// Convert string ID to int for map lookup
	id, err := strconv.Atoi(task.ID)
	if err != nil {
		HandleError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}
	mu.Lock()         // Lock the mutex to safely modify shared data
	defer mu.Unlock() // Ensure the mutex is unlocked after the function returns
	if _, exists := tasks[id]; !exists {
		http.Error(w, "Task not found", http.StatusNotFound) // Respond with 404 if the task doesn't exist
		return
	}
	tasks[id] = task                        // Update the task in the map
	RespondWithJSON(w, http.StatusOK, task) // Respond with the updated task as JSON
}

// DeleteTask handles DELETE requests to remove a specific task by ID
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id") // Get the task ID from the query parameters (?id=)
	id, err := strconv.Atoi(idStr)   // Convert id from string to int for comparison
	if err != nil {
		HandleError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}
	mu.Lock()         // Lock the mutex to safely modify shared data
	defer mu.Unlock() // Ensure the mutex is unlocked after the function returns
	if _, exists := tasks[id]; exists {
		delete(tasks, id)                             // Delete the task from the map
		RespondWithJSON(w, http.StatusNoContent, nil) // Respond with 204 No Content
		return
	}
	http.Error(w, "Task not found", http.StatusNotFound) // Respond with 404 if the task doesn't exist
}
