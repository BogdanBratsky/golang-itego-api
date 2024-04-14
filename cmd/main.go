package main

import (
	"log"
	"net/http"

	"github.com/bogdanbratsky/golang-itego-api/db"
	"github.com/bogdanbratsky/golang-itego-api/internal/handlers"
)

func main() {
	db.InitDB()
	defer db.CloseDB()

	http.HandleFunc("/api/register", handlers.CreateUserHandler)
	http.HandleFunc("/api/login", handlers.LoginHandler)

	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Println("Error: ", err)
		return
	}
}
