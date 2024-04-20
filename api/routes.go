package api

import (
	"github.com/bogdanbratsky/golang-itego-api/internal/handlers"
	"github.com/gorilla/mux"
)

// !USERS!
//
// POST /api/register
// POST /api/login
// GET /api/users
// GET /api/users/{id}
//
// !EMAIL!
//
// POST /api/send-email
//
// !ARTICLES!
//
//
//
// !CATEGORIES!

func Router() *mux.Router {
	r := mux.NewRouter()

	// USER ROUTES
	r.HandleFunc("/api/register", handlers.CreateUserHandler).Methods("POST")
	r.HandleFunc("/api/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/api/users", handlers.GetUsersHandler).Methods("GET")

	return r
}
