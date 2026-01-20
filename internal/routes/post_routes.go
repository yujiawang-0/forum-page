package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/yujiawang-0/forum-page/internal/database"
	"github.com/yujiawang-0/forum-page/internal/handlers/posts"
)


func MainPostRoutes(db *database.Database) func(r chi.Router) {
	return func(r chi.Router) {
		handler := &posts.PostHandler{
			DB: db,
		}
		// /posts
		r.Get("/", handler.HandleGetAllPosts) // shows all the posts regardless of topic
	}
}

func NestedPostRoutes(db *database.Database) func(r chi.Router) {
	return func(r chi.Router) {
		handler := &posts.PostHandler{
			DB: db,
		}
		// topics/{topic_id}/posts
		r.Get("/", handler.HandleGetPostsByTopicID) // shows all the posts regardless of topic
		r.Get("/{post_id}", handler.HandleGetPostByID) // right now this disregards the topic_id entirely (query only by postid)
		r.Post("/", handler.HandleCreatePost)
		r.Put("/{post_id}", handler.HandleUpdatePost)
		r.Delete("/{post_id}", handler.HandleDeletePost)
	}
}