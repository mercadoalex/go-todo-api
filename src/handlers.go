package main

import (
	"encoding/json"
	"net/http"
	"strconv" // For converting between string and int

	"github.com/gorilla/mux"
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
	task.ID = int(id)                            // Assign int64 to int for Task.ID
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
		if err := rows.Scan(&task.ID, &task.Title, &task.Completed); err != nil {
			HandleError(w, http.StatusInternalServerError, err.Error())
			return
		}
		tasks = append(tasks, task)
	}
	RespondWithJSON(w, http.StatusOK, tasks) // Respond with all tasks as JSON
}

// UpdateTask handles PUT requests to update an existing task in the SQLite database
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		HandleError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		HandleError(w, http.StatusBadRequest, err.Error())
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
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	task.ID = id // Set the correct ID in the response
	RespondWithJSON(w, http.StatusOK, task)
}

// DeleteTask handles DELETE requests to remove a specific task by ID from the SQLite database
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
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
