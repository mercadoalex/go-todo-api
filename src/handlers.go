package main

import (
	"encoding/json"
	"net/http"
	"strconv" // For converting between string and int
	"sync"
)

// Global variables for in-memory storage and concurrency control
var (
	tasks  = make(map[int]Task) // Stores tasks using int IDs as keys
	nextID = 1                  // Next available unique ID for new tasks
	mu     sync.Mutex           // Mutex to ensure thread-safe access to the tasks map
)

// AddTask handles POST requests to create a new task in the SQLite database
func AddTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	// Decode the JSON request body into a Task struct
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		HandleError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}
	// Validate that the title is not empty
	if task.Title == "" {
		HandleError(w, http.StatusBadRequest, "Task title cannot be empty")
		return
	}
	// Insert the new task into the database
	res, err := db.Exec("INSERT INTO tasks (title, completed) VALUES (?, ?)", task.Title, task.Completed)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Get the ID of the newly inserted task
	id, _ := res.LastInsertId()
	task.ID = strconv.FormatInt(id, 10)          // Convert int64 to string for Task.ID
	RespondWithJSON(w, http.StatusCreated, task) // Respond with the created task as JSON
}

// GetTasks handles GET requests to retrieve all tasks from the SQLite database
func GetTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, completed FROM tasks")
	if err != nil {
		HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()
	var tasks []Task
	for rows.Next() {
		var task Task
		var completed int
		// Scan each row into a Task struct
		if err := rows.Scan(&task.ID, &task.Title, &completed); err != nil {
			HandleError(w, http.StatusInternalServerError, err.Error())
			return
		}
		task.Completed = completed == 1 // Convert SQLite integer to Go bool
		tasks = append(tasks, task)
	}
	RespondWithJSON(w, http.StatusOK, tasks) // Respond with all tasks as JSON
}

// UpdateTask handles PUT requests to update an existing task in the SQLite database
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	// Decode the JSON request body into a Task struct
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		HandleError(w, http.StatusBadRequest, err.Error())
		return
	}
	// Convert string ID to int for database query
	id, err := strconv.Atoi(task.ID)
	if err != nil {
		HandleError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}
	// Update the task in the database
	res, err := db.Exec("UPDATE tasks SET title = ?, completed = ? WHERE id = ?", task.Title, task.Completed, id)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound) // Respond with 404 if the task doesn't exist
		return
	}
	RespondWithJSON(w, http.StatusOK, task) // Respond with the updated task as JSON
}

// DeleteTask handles DELETE requests to remove a specific task by ID from the SQLite database
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id") // Get the task ID from the query parameters (?id=)
	id, err := strconv.Atoi(idStr)   // Convert id from string to int for database query
	if err != nil {
		HandleError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}
	// Delete the task from the database
	res, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound) // Respond with 404 if the task doesn't exist
		return
	}
	RespondWithJSON(w, http.StatusNoContent, nil) // Respond with 204 No Content if successful
}
