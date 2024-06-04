package todo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Y0sh1dk/kubebuilder-external-resource-demo/internal/todo"
)

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
	}
}

func (c *Client) CreateTodo(t todo.Todo) (*todo.Todo, error) {
	createTodoURL := fmt.Sprintf("%s/todos", c.BaseURL)
	createTodoJSON, _ := json.Marshal(t)
	resp, err := http.Post(createTodoURL, "application/json", bytes.NewBuffer(createTodoJSON))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var createdTodo todo.Todo
	err = json.NewDecoder(resp.Body).Decode(&createdTodo)
	if err != nil {
		return nil, err
	}

	return &createdTodo, nil
}

func (c *Client) GetTodo(id int) (*todo.Todo, error) {
	getTodoURL := fmt.Sprintf("%s/todos/%d", c.BaseURL, id)
	resp, err := http.Get(getTodoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var todo todo.Todo
	err = json.NewDecoder(resp.Body).Decode(&todo)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (c *Client) GetTodos() ([]todo.Todo, error) {
	getTodosURL := fmt.Sprintf("%s/todos", c.BaseURL)
	resp, err := http.Get(getTodosURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var todos []todo.Todo
	err = json.NewDecoder(resp.Body).Decode(&todos)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (c *Client) UpdateTodo(t todo.Todo) (*todo.Todo, error) {
	updateTodoURL := fmt.Sprintf("%s/todos/%d", c.BaseURL, t.ID)
	updateTodoJSON, _ := json.Marshal(t)
	req, _ := http.NewRequest(http.MethodPut, updateTodoURL, bytes.NewBuffer(updateTodoJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updatedTodo todo.Todo
	err = json.NewDecoder(resp.Body).Decode(&updatedTodo)
	if err != nil {
		return nil, err
	}

	return &updatedTodo, nil
}

func (c *Client) DeleteTodo(id int) error {
	deleteTodoURL := fmt.Sprintf("%s/todos/%d", c.BaseURL, id)
	req, _ := http.NewRequest(http.MethodDelete, deleteTodoURL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
