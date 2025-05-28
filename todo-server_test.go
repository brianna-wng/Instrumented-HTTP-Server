package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTodos(t *testing.T) {
	// reset state and populate with set todos
	todos = []Todo{
		{ID: 1, Title: "Test Task", Completed: false},
		{ID: 2, Title: "Test Task 2", Completed: true},
		{ID: 3, Title: "Test Task 3", Completed: false},
	}

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()

	getTodos(w, req)

	res := w.Result()
	defer res.Body.Close()

	var received []Todo
	err := json.NewDecoder(res.Body).Decode(&received)
	assert.NoError(t, err)
	assert.Len(t, received, 3)
	assert.Equal(t, todos, received)
}

func TestAddTodo(t *testing.T) {
	// reset shared state to ensure tests are isolated
	todos = []Todo{}
	nextID = 1

	// prepare request body
	body := []byte(`{"title":"Test Task"}`)
	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	addTodo(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var todo Todo
	err := json.NewDecoder(res.Body).Decode(&todo)
	assert.NoError(t, err)
	assert.Equal(t, "Test Task", todo.Title)
	assert.Equal(t, 1, todo.ID)
	assert.False(t, todo.Completed)
}

func TestMarkCompleted(t *testing.T) {
	// reset shared state to ensure tests are isolated
	todos = []Todo{
		{ID: 1, Title: "Test Task", Completed: false},
		{ID: 2, Title: "Test Task 2", Completed: true},
		{ID: 3, Title: "Test Task 3", Completed: false},
	}

	req := httptest.NewRequest(http.MethodPut, "/todos/1", nil)
	w := httptest.NewRecorder()

	markCompleted(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var updated Todo
	err := json.NewDecoder(res.Body).Decode(&updated)
	assert.NoError(t, err)
	assert.Equal(t, "Test Task", updated.Title)
	assert.Equal(t, 1, updated.ID)
	assert.True(t, updated.Completed)
}