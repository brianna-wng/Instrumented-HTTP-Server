package main

import (
	"encoding/json" // encode/decode JSON data
	"fmt"           // formatted I/O
	"log"           // logging errors
	"net/http"      // building HTTP servers
	"strconv"       // convert strings to integers
	"strings"       // string manipulation
	"sync"          // synchronization primitives, ex. mutexes for safe concurrent access
	"time"
	"bytes"
	"os"

	"github.com/DataDog/datadog-go/v5/statsd"
)

// creating struct for each todo
type Todo struct {
	ID        int    `json:"id"` // specifies how fields are named when serialized/deserialized in JSON, capitalized letters represent exported/public fields
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var (
	todos   = []Todo{} // dynamic array holding Todo structs
	nextID  = 1
	todoMux = sync.Mutex{} // mutex/lock to ensure safe concurrent access to todos slice
	statsdClient *statsd.Client
	logger *log.Logger
)

func init() {
	var err error
	statsdClient, err = statsd.New("127.0.0.1:8125")
	if err != nil {
		log.Fatalf("Error creating statsd client: %v", err)
	}

	file, err := os.OpenFile("todo-server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error creating log file: %v", err)
	}
	logger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// handles GET /todos and writes todos to response
func getTodos(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	
	w.Header().Set("Content-Type", "application/json") // set response header so clients know they're getting a JSON
	todoMux.Lock()                                     // lock mutex to safely access todos
	json.NewEncoder(w).Encode(todos)                   // serializes todos and writes to response
	todoMux.Unlock()

	statsdClient.Timing("get_todos.duration", time.Since(start), nil, 1)
	statsdClient.Count("get_todos.count", 1, nil, 1)
	logger.Printf("GET /todos: %d todos returned", len(todos))

}

// handles POST /todos and responds with newly created todo as JSON
func addTodo(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	if req.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed) // 405 Method Not Allowed
		return
	}

	var newTodo Todo
	err := json.NewDecoder(req.Body).Decode(&newTodo) // decodes request body into newTodo
	if err != nil || strings.TrimSpace(newTodo.Title) == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest) // 400 Bad Request
		return
	}

	// adds new todo to todos slice
	todoMux.Lock()
	newTodo.ID = nextID
	nextID++
	todos = append(todos, newTodo)
	todoMux.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(newTodo)
}

func markCompleted(w http.ResponseWriter, req *http.Request) {
	if req.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(req.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest) // 400
		return
	}

	todoMux.Lock()
	defer todoMux.Unlock()
	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Completed = true
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(todos[i])
			return
		}
	}

	http.Error(w, "Todo not found", http.StatusNotFound) // 404 Not Found
}

func main() {
	http.HandleFunc("/todos", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			getTodos(w, req)
		} else if req.Method == "POST" {
			addTodo(w, req)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/todos/", markCompleted)

	fmt.Println("Server started on http://localhost:8090")
	logger.Println("Server started on http://localhost:8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
