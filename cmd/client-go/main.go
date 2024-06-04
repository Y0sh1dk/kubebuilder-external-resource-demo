package main

import (
	"log"

	todoClient "github.com/Y0sh1dk/kubebuilder-external-resource-demo/internal/clients/todo"
	"github.com/Y0sh1dk/kubebuilder-external-resource-demo/internal/todo"

	"github.com/k0kubun/pp/v3"
)

func main() {
	client := todoClient.NewClient("http://localhost:8080")

	log.Println("Creating a Todo")
	milkTodo, err := client.CreateTodo(todo.Todo{
		Title: "Buy milk",
	})
	if err != nil {
		panic(err)
	}

	log.Println("Getting single Todo")
	t, err := client.GetTodo(milkTodo.ID)
	if err != nil {
		panic(err)
	}

	pp.Println(t)

	log.Println("Getting all Todos")
	todos, err := client.GetTodos()
	if err != nil {
		panic(err)
	}

	pp.Println(todos)

	log.Println("Updating a Todo")
	_, err = client.UpdateTodo(todo.Todo{
		ID:    milkTodo.ID,
		Title: "Buy milk and eggs",
	})
	if err != nil {
		panic(err)
	}

	log.Println("Getting all Todos")
	todos, err = client.GetTodos()
	if err != nil {
		panic(err)
	}

	pp.Println(todos)

	log.Println("Deleting a Todo")
	err = client.DeleteTodo(milkTodo.ID)
	if err != nil {
		panic(err)
	}

	log.Println("Getting all Todos")
	todos, err = client.GetTodos()
	if err != nil {
		panic(err)
	}

	pp.Println(todos)
}
