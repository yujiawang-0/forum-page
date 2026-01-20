package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yujiawang-0/forum-page/internal/database"
	"github.com/yujiawang-0/forum-page/internal/routes"
)

// HTTP routing

func Setup(db *database.Database) chi.Router {
	r := chi.NewRouter()
	
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	setUpUserRoutes(r, db)
	return r
}

func setUpUserRoutes(r chi.Router, db *database.Database) {
	r.Route("/users", routes.UserRoutes(db))
	r.Route("/topics", routes.TopicRoutes(db))
	r.Route("/posts", routes.MainPostRoutes(db))
	// r.Group(routes.GetRoutes())
}
