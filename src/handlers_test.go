package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestAddTask checks if the AddTask handler correctly creates a new task
func TestAddTask(t *testing.T) {
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
}

// TestGetTasks checks if the GetTasks handler returns all tasks successfully
func TestGetTasks(t *testing.T) {
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
}

// TestUpdateTask checks if the UpdateTask handler correctly updates an existing task
func TestUpdateTask(t *testing.T) {
	// Prepare updated task data (ID must exist in the DB for this test to pass)
	task := Task{ID: 1, Title: "Updated Task", Completed: true}
	body, _ := json.Marshal(task) // Convert the updated task to JSON

	// Create a new HTTP PUT request to /tasks/1 with the JSON body
	req, err := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err) // Fail the test if the request can't be created
	}
	req.Header.Set("Content-Type", "application/json") // Set the content type to JSON

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateTask) // Use the UpdateTask handler

	// Serve the HTTP request using the handler
	handler.ServeHTTP(rr, req)

	// Check if the status code is 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

// TestDeleteTask checks if the DeleteTask handler correctly deletes a task
func TestDeleteTask(t *testing.T) {
	// Create a new HTTP DELETE request to /tasks/1
	req, err := http.NewRequest("DELETE", "/tasks/1", nil)
	if err != nil {
		t.Fatal(err) // Fail the test if the request can't be created
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteTask) // Use the DeleteTask handler

	// Serve the HTTP request using the handler
	handler.ServeHTTP(rr, req)

	// Check if the status code is 204 No Content
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}
