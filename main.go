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

	// Create router:
	router := mux.NewRouter()
	router.HandleFunc("/", indexPage)
	router.HandleFunc("GET /users", getUsers)
	router.HandleFunc("GET /users/{id}", getUser)
	router.HandleFunc("POST /users/{id}", createUser)
	router.HandleFunc("PUT /users/{id}", updateUser)
	router.HandleFunc("DELETE /users/{id}", deleteUser)

	// Start server (port, router):
	err := http.ListenAndServe(":8080", jsonContentTypeMiddleware(router)) // sets header type for all routes
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
func indexPage(w http.ResponseWriter, r *http.Request)  {}
func getUsers(w http.ResponseWriter, r *http.Request)   {}
func getUser(w http.ResponseWriter, r *http.Request)    {}
func createUser(w http.ResponseWriter, r *http.Request) {}
func updateUser(w http.ResponseWriter, r *http.Request) {}
func deleteUser(w http.ResponseWriter, r *http.Request) {}
