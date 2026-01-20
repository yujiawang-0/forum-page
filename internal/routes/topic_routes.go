package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/yujiawang-0/forum-page/internal/database"
	"github.com/yujiawang-0/forum-page/internal/handlers/topics"
)

// topic routes includes topic routes
// as well as post routes that

func TopicRoutes(db *database.Database) func(r chi.Router) {
	return func(r chi.Router) {
		handler := &topics.TopicHandler{
			DB: db,
		}
		// already the link is /topics
		r.Get("/", handler.HandleGetAllTopics)
		r.Post("/", handler.HandleCreateTopic)

		// topics/{topic_id}
		r.Route("/{topic_id}", func (r chi.Router) {
			r.Get("/", handler.HandleGetTopicByID)
			r.Put("/", handler.HandleUpdateTopic)
			r.Delete("/", handler.HandleDeleteTopic)
			
			// topics/{topic_id}/posts
			r.Route("/posts", NestedPostRoutes(db))
		})
		
	}
}