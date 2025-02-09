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
	// Connection to DB:
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Routing and Handler functions:

}
