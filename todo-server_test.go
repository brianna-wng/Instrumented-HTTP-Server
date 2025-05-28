package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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