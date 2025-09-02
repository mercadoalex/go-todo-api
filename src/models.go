package main

type Task struct {
	ID        string `json:"id"`        // Unique identifier for the task
	Title     string `json:"title"`     // Description or title of the task
	Completed bool   `json:"completed"` // Status: true if the task is done, false otherwise
}
