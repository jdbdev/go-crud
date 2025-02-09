package main

/*
tutorial:https://www.youtube.com/watch?v=aLVJY-1dKz8&t=368s
*/

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Variables / Structs / Arrays:

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	fmt.Println("go-crud")
	// Open connection to DB:
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create router and endpoints:
	router := mux.NewRouter()
	router.HandleFunc("/", indexPage)
	router.HandleFunc("GET /users", getUsers)
	router.HandleFunc("GET /users/{id}", getUser)
	router.HandleFunc("POST /users/{id}", createUser)
	router.HandleFunc("PUT /users/{id}", updateUser)
	router.HandleFunc("DELETE /users/{id}", deleteUser)

	// Start server (port, router):
	err = http.ListenAndServe(":8080", jsonContentTypeMiddleware(router)) // sets header type for all routes
	if err != nil {
		log.Fatal(err)
	}
}

// Middleware:
func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	})
}

// Handlers:

// Homepage:
func indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "User Portal")
}

// Get all users:
func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer rows.Close()

		users := []User{}
		for rows.Next() {
			var u User
			err = rows.Scan(&u.ID, &u.Name, &u.Email)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			users = append(users, u)
		}
		if err = rows.Err(); err != nil {
		}
		json.NewEncoder(w).Encode(users)
	}
}

// Get user by {id}:
func getUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var u User
		err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		json.NewEncoder(w).Encode(u)
	}
}

// Create new user:
func createUser(w http.ResponseWriter, r *http.Request) {}

// Update user:
func updateUser(w http.ResponseWriter, r *http.Request) {}

// Delete user:
func deleteUser(w http.ResponseWriter, r *http.Request) {}
