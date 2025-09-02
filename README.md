# Go To-Do API

This project is a simple REST API for managing a to-do list, built using Go. It allows users to perform basic CRUD operations: create, read, update, and delete tasks.

## Project Structure

```
go-todo-api
├── src
│   ├── main.go         # Entry point of the application
│   ├── handlers.go     # HTTP handler functions for CRUD operations
│   ├── models.go       # Data structures for tasks
│   └── utils.go        # Utility functions for JSON handling and error management
├── go.mod              # Go module configuration
└── README.md           # Project documentation
```

## Getting Started

### Prerequisites

- Go installed on your machine (version 1.16 or later)
- A terminal or command prompt

### Installation

1. Clone the repository:

   ```
   git clone https://github.com/yourusername/go-todo-api.git
   cd go-todo-api
   ```

2. Navigate to the `src` directory:

   ```
   cd src
   ```

3. Initialize Go modules (if not already done):

   ```
   go mod init go-todo-api
   ```

### Running the API

1. Start the server:

   ```
   go run main.go
   ```

2. The API will be running at `http://localhost:8080`.

### API Endpoints

- **Add Task**: `POST /tasks`
  - Request Body: JSON object with `title` and `completed` fields.
  
- **Get Tasks**: `GET /tasks`
  - Response: JSON array of tasks.

- **Update Task**: `PUT /tasks/{id}`
  - Request Body: JSON object with `title` and `completed` fields.

- **Delete Task**: `DELETE /tasks/{id}`

### Example Requests

- **Add Task**:

  ```
  curl -X POST http://localhost:8080/tasks -d '{"title": "Learn Go", "completed": false}' -H "Content-Type: application/json"
  ```

- **Get Tasks**:

  ```
  curl http://localhost:8080/tasks
  ```

- **Update Task**:

  ```
  curl -X PUT http://localhost:8080/tasks/1 -d '{"title": "Learn Go", "completed": true}' -H "Content-Type: application/json"
  ```

- **Delete Task**:

  ```
  curl -X DELETE http://localhost:8080/tasks/1
  ```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author
Alejandro Mercado Peña
