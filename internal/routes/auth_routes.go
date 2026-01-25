package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/yujiawang-0/forum-page/internal/database"
	"github.com/yujiawang-0/forum-page/internal/handlers/auth"
)

func AuthRoutes(db *database.Database) func(r chi.Router) {
	return func(r chi.Router) {
		handler := &auth.AuthHandler{
			DB: db,
		}
		// POST /auth/login
		r.Post("/login", handler.HandleLogin) // handles login and provides valid JWT token
	}
}