package api

import (
	"github.com/bogdanbratsky/golang-itego-api/internal/handlers"
	"github.com/gorilla/mux"
)

// LIST OF ENDPOINTS:
//
// !USERS!
//
// POST /api/v1/register
// POST /api/v1/login
// GET /api/v1/users
// GET /api/v1/users/{id}
//
// !ARTICLES!
//
// POST /api/v1/articles
// DELETE /api/v1/articles/{id}
// PATCH /api/v1/articles/{id}
// GET /api/v1/articles
// GET /api/v1/articles/{id}
//
// !CATEGORIES!
//
// POST /api/v1/categories
// DELETE /api/v1/categories/{id}
// PATCH /api/v1/categories/{id}
// GET /api/v1/categories
// GET /api/v1/articles/{id}
//
// !EMAIL!
//
// POST /api/v1/send-email
//

func Router() *mux.Router {
	r := mux.NewRouter()

	// USER ROUTES
	r.HandleFunc("/api/v1/register", handlers.CreateUserHandler).Methods("POST")
	r.HandleFunc("/api/v1/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/api/v1/users", handlers.GetUsersHandler).Methods("GET")
	r.HandleFunc("/api/v1/users/{id}", handlers.GetUserByIdHandler).Methods("GET")
	r.HandleFunc("/api/v1/checktoken", handlers.CheckToken).Methods("GET")
	r.HandleFunc("/api/v1/users/{id}", handlers.DeleteUserHandler).Methods("DELETE")

	// ARTICLE ROUTES
	r.HandleFunc("/api/v1/articles", handlers.CreateArticleHandler).Methods("POST")
	r.HandleFunc("/api/v1/articles", handlers.GetArticlesHandler).Methods("GET")
	r.HandleFunc("/api/v1/articles/{id}", handlers.DeleteArticleHandler).Methods("DELETE")
	r.HandleFunc("/api/v1/articles/{id}", handlers.UpdateArticleHandler).Methods("PATCH")
	r.HandleFunc("/api/v1/articles/{id}", handlers.GetArticleByIdHandler).Methods("GET")

	r.HandleFunc("/api/v1/categories/{id}/articles", handlers.GetArticlesByCategoryHandler).Methods("GET")

	// CATEGORY ROUTES
	r.HandleFunc("/api/v1/categories", handlers.CreateCategoryHandler).Methods("POST")
	r.HandleFunc("/api/v1/categories", handlers.GetCategoriesHandler).Methods("GET")
	r.HandleFunc("/api/v1/categories/{id}", handlers.DeleteCategoryHandler).Methods("DELETE")
	r.HandleFunc("/api/v1/categories/{id}", handlers.UpdateCategoryHandler).Methods("PATCH")
	r.HandleFunc("/api/v1/categories/{id}", handlers.GetCategoryByIdHandler).Methods("GET")

	// EMAIL
	r.HandleFunc("/api/v1/send-email", handlers.SendEmailHandler).Methods("POST")

	return r
}
