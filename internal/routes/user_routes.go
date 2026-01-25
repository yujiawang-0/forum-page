package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/yujiawang-0/forum-page/internal/database"
	"github.com/yujiawang-0/forum-page/internal/handlers/users"
)

func UserRoutes(db *database.Database) func(r chi.Router) {
	return func(r chi.Router) {
		handler := &users.UserHandler{
			DB: db,
		}
		r.Get("/", handler.HandleGetAllUsers)
		r.Get("/user/{id}", handler.HandleGetUserByID)
		//r.Get("/user/{username}", handler.HandleGetUserByUsername)
		r.Post("/user", handler.HandleCreateUser)
		r.Put("/user/{id}", handler.HandleUpdateUser)
		r.Delete("/user/{id}", handler.HandleDeleteUser)
	}
}
