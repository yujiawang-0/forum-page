package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/yujiawang-0/forum-page/internal/database"
	"github.com/yujiawang-0/forum-page/internal/handlers/topics"
	"github.com/yujiawang-0/forum-page/internal/handlers/auth"
)

// TopicRoutes includes topic routes as well as post routes that go through the topic first

func TopicRoutes(db *database.Database) func(r chi.Router) {
	return func(r chi.Router) {
		handler := &topics.TopicHandler{
			DB: db,
		}
		// already the link is /topics: GET /topics
		r.Get("/", handler.HandleGetAllTopics)
		
		r.Group(func (r chi.Router) { 
			r.Use(auth.RequireAuth)
			// POST /topics
			r.Post("/", handler.HandleCreateTopic)
		})
			

		// topics/{topic_id}
		r.Route("/{topic_id}", func (r chi.Router) {
			// GET topics/{topic_id}
			r.Get("/", handler.HandleGetTopicByID)

			r.Group(func (r chi.Router) { 
				r.Use(auth.RequireAuth)
				// PUT /topics/{topic_id}
				r.Put("/", handler.HandleUpdateTopic)
				// DELETE /topics/{topic_id}
				r.Delete("/", handler.HandleDeleteTopic)
			})
			
			// topics/{topic_id}/posts
			r.Route("/posts", NestedPostRoutes(db))
		})
		
	}
}