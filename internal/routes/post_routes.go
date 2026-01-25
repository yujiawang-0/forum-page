package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/yujiawang-0/forum-page/internal/database"
	"github.com/yujiawang-0/forum-page/internal/handlers/auth"
	"github.com/yujiawang-0/forum-page/internal/handlers/posts"
)


func MainPostRoutes(db *database.Database) func(r chi.Router) {
	return func(r chi.Router) {
		handler := &posts.PostHandler{
			DB: db,
		}
		// GET /posts
		r.Get("/", handler.HandleGetAllPosts) // shows all the posts regardless of topic
	}
}

func NestedPostRoutes(db *database.Database) func(r chi.Router) {
	return func(r chi.Router) {
		handler := &posts.PostHandler{
			DB: db,
		}
		// GET topics/{topic_id}/posts
		r.Get("/", handler.HandleGetPostsByTopicID) // shows all the posts about that topic

		// GET topics/{topic_id}/posts/{post_id}
		r.Get("/{post_id}", handler.HandleGetPostByID) // right now this disregards the topic_id entirely (query only by postid)

		r.Group(func (r chi.Router) {
			r.Use(auth.RequireAuth)

			// POST /topics/{topic_id}/posts
			r.Post("/", handler.HandleCreatePost)
			// PUT /topics/{topic_id}/posts/{post_id}
			r.Put("/{post_id}", handler.HandleUpdatePost)
			// DELETE /topics/{topic_id}/posts/{post_id}
			r.Delete("/{post_id}", handler.HandleDeletePost)
		})
		
	}
}