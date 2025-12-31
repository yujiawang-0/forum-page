package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yujiawang-0/forum-page/internal/routes"
)

// HTTP routing

func Setup() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	setUpRoutes(r)
	return r
}

func setUpRoutes(r chi.Router) {
	r.Group(routes.GetRoutes())
}
