package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Y0sh1dk/kubebuilder-external-resource-demo/internal/todo"
	"github.com/gorilla/mux"
)

var todos []todo.Todo

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/todos", getTodos).Methods("GET")
	router.HandleFunc("/todos/{id}", getTodo).Methods("GET")
	router.HandleFunc("/todos", createTodo).Methods("POST")
	router.HandleFunc("/todos/{id}", updateTodo).Methods("PUT")
	router.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")

	log.Println("Starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting all todos")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(todos)
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	log.Println("Getting todo with id", id)

	for _, todo := range todos {
		if todo.ID == id {
			_ = json.NewEncoder(w).Encode(todo)
			return
		}
	}
	log.Println("Todo not found")

	http.NotFound(w, r)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todo todo.Todo
	_ = json.NewDecoder(r.Body).Decode(&todo)
	todo.ID = len(todos) + 1

	log.Println("Creating todo with id", todo.ID)

	todos = append(todos, todo)
	_ = json.NewEncoder(w).Encode(todo)
}

// Update todo
func updateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	log.Println("Updating todo with id", id)

	var updatedTodo todo.Todo
	_ = json.NewDecoder(r.Body).Decode(&updatedTodo)

	for index, todo := range todos {
		if todo.ID == id {
			todos[index].Title = updatedTodo.Title
			updatedTodo = todos[index]
			break
		}
	}

	_ = json.NewEncoder(w).Encode(updatedTodo)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	log.Println("Deleting todo with id", id)

	for index, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:index], todos[index+1:]...)
			break
		}
	}

	_ = json.NewEncoder(w).Encode(todos)
}
