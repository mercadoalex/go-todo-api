package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// Test database file name
var testDBFile = "test_todo.db"

// TestMain sets up the database connection and schema before running tests
func TestMain(m *testing.M) {
	var err error
	// Open a new SQLite database for testing
	db, err = sql.Open("sqlite3", testDBFile)
	if err != nil {
		log.Fatal("Failed to open test database:", err)
	}
	// Create the tasks table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT 0
		)
	`)
	if err != nil {
		log.Fatal("Failed to create tasks table:", err)
	}
	// Run all tests
	code := m.Run()
	// Close the database and remove the test file after tests
	db.Close()
	os.Remove(testDBFile)
	os.Exit(code)
}

// TestAddTask checks if the AddTask handler correctly creates a new task
func TestAddTask(t *testing.T) {
	log.Println("Running TestAddTask...")
	// Create a sample task to send in the request body
	task := Task{Title: "Test Task", Completed: false}
	body, _ := json.Marshal(task) // Convert the task struct to JSON

	// Create a new HTTP POST request to /tasks with the JSON body
	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err) // Fail the test if the request can't be created
	}
	req.Header.Set("Content-Type", "application/json") // Set the content type to JSON

	// Create a ResponseRecorder to capture the handler's response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddTask) // Use the AddTask handler

	// Serve the HTTP request using the handler
	handler.ServeHTTP(rr, req)

	// Check if the status code is 201 Created
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
	log.Println("TestAddTask completed.")
}

// TestGetTasks checks if the GetTasks handler returns all tasks successfully
func TestGetTasks(t *testing.T) {
	log.Println("Running TestGetTasks...")
	// Create a new HTTP GET request to /tasks
	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err) // Fail the test if the request can't be created
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTasks) // Use the GetTasks handler

	// Serve the HTTP request using the handler
	handler.ServeHTTP(rr, req)

	// Check if the status code is 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	log.Println("TestGetTasks completed.")
}

// TestUpdateTask checks if the UpdateTask handler correctly updates an existing task
func TestUpdateTask(t *testing.T) {
	log.Println("Running TestUpdateTask...")
	// First, create a task to update
	task := Task{Title: "Initial Task", Completed: false}
	body, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	http.HandlerFunc(AddTask).ServeHTTP(rr, req)

	// Now, update the created task (ID will be 1)
	update := Task{ID: 1, Title: "Updated Task", Completed: true}
	updateBody, _ := json.Marshal(update)
	updateReq, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateRR := httptest.NewRecorder()

	// Use mux router to handle path variables for UpdateTask
	router := mux.NewRouter()
	router.HandleFunc("/tasks/{id}", UpdateTask).Methods("PUT")
	router.ServeHTTP(updateRR, updateReq)

	// Check if the status code is 200 OK
	if status := updateRR.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	log.Println("TestUpdateTask completed.")
}

// TestDeleteTask checks if the DeleteTask handler correctly deletes a task
func TestDeleteTask(t *testing.T) {
	log.Println("Running TestDeleteTask...")
	// First, create a task to delete
	task := Task{Title: "Task to Delete", Completed: false}
	body, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	http.HandlerFunc(AddTask).ServeHTTP(rr, req)

	// Now, delete the created task (ID will be 1)
	deleteReq, _ := http.NewRequest("DELETE", "/tasks/1", nil)
	deleteRR := httptest.NewRecorder()

	// Use mux router to handle path variables for DeleteTask
	router := mux.NewRouter()
	router.HandleFunc("/tasks/{id}", DeleteTask).Methods("DELETE")
	router.ServeHTTP(deleteRR, deleteReq)

	// Check if the status code is 204 No Content
	if status := deleteRR.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
	log.Println("TestDeleteTask completed.")
}
