package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Data model
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// In-memory data store
var users = map[string]User{
	"1": {ID: "1", Name: "Mahiru Sanketh", Email: "mahirusanketh21@gmail.com"},
	"2": {ID: "2", Name: "Shalika Damayanthi", Email: "shalika@gmail.com"},
	"3": {ID: "3", Name: "Thushantha Mallika Arachchi", Email: "thushantha@gmail.com"},
}

// writeJSON is a helper — we'll use this in every handler
func writeJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError sends a consistent error shape
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJson(w, status, msg)
}

// ListUsers handles GET /users
func ListUsers(w http.ResponseWriter, r *http.Request) {
	list := make([]User, 0, len(users))

	for _, u := range users {
		list = append(list, u)
	}

	writeJson(w, http.StatusOK, list)
}

// GetUser handles GET /users/{id}
func GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	user, ok := users[id]

	if !ok {
		writeError(w, http.StatusNotFound, fmt.Sprintf("user id %s is not found", id))
		return
	}

	writeJson(w, http.StatusOK, user)
}

// CreateUser handles POST /users
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var input User

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid json body")
		return
	}

	// Basic input validation
	if input.Name == "" {
		writeJson(w, http.StatusNotFound, "Name is not found")
		return
	}

	if input.Email == "" {
		writeJson(w, http.StatusNotFound, "Email is not found")
		return
	}

	// Generate basic id
	input.ID = fmt.Sprintf("%d", len(users)+1)

	users[input.ID] = input

	writeJson(w, http.StatusCreated, users[input.ID])

}

// DeleteUser handles DELETE /user/{id}
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if _, ok := users[id]; !ok {
		writeError(w, http.StatusNotFound, "User is not found")
		return
	}

	delete(users, id)
	w.WriteHeader(http.StatusNoContent) // 204 — success, no body
}

// UpdateUser handles PUT /user/{id}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var input User

	if _, ok := users[id]; !ok {
		writeError(w, http.StatusNotFound, "User is not found!")
		return
	}

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid json body")
		return
	}

	user := users[id]

	user.Email = input.Email
	user.Name = input.Name

	users[id] = user

	writeJson(w, http.StatusOK, user)

}
