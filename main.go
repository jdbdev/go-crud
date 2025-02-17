package main

/*
tutorial:https://www.youtube.com/watch?v=aLVJY-1dKz8&t=368s
*/

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
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

	// Create db tables if not exist:
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT, email TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	// Create router and endpoints:
	router := mux.NewRouter()
	router.HandleFunc("/", indexPage)
	router.HandleFunc("GET /users", getUsers(db))
	router.HandleFunc("GET /users/{id}", getUser(db))
	router.HandleFunc("POST /users/{id}", createUser(db))
	router.HandleFunc("PUT /users/{id}", updateUser(db))
	router.HandleFunc("DELETE /users/{id}", deleteUser(db))

	// Start server (port, router):
	err = http.ListenAndServe(":8000", jsonContentTypeMiddleware(router)) // sets header type for all routes
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

// Handler functions:

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
		err = db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		json.NewEncoder(w).Encode(u)
	}
}

// Create new user:
func createUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		json.NewDecoder(r.Body).Decode(&u)

		err := db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", u.Name, u.Email).Scan(&u.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(u)
	}
}

// Update user:
func updateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		json.NewDecoder(r.Body).Decode(&u)

		vars := mux.Vars(r)
		id := vars["id"]

		_, err := db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", u.Name, u.Email, id)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(u)
	}
}

// Delete user:
func deleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u User
		err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
			if err != nil {
				//todo : fix error handling
				w.WriteHeader(http.StatusNotFound)
				return
			}

			json.NewEncoder(w).Encode("User deleted")
		}
	}
}
