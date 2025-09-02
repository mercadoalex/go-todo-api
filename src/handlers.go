package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

var (
	tasks  = make(map[int]Task)
	nextID = 1
	mu     sync.Mutex
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		HandleError(w, err, http.StatusBadRequest)
		return
	}
	mu.Lock()
	task.ID = nextID
	nextID++
	tasks[task.ID] = task
	mu.Unlock()
	RespondWithJSON(w, http.StatusCreated, task)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	RespondWithJSON(w, http.StatusOK, tasks)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		HandleError(w, err, http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	if _, exists := tasks[task.ID]; !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	tasks[task.ID] = task
	RespondWithJSON(w, http.StatusOK, task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	mu.Lock()
	defer mu.Unlock()
	for k := range tasks {
		if k == id {
			delete(tasks, k)
			RespondWithJSON(w, http.StatusNoContent, nil)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}