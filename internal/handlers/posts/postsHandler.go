package posts

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/yujiawang-0/forum-page/internal/api"
	"github.com/yujiawang-0/forum-page/internal/dataaccess"
	"github.com/yujiawang-0/forum-page/internal/models"

	postService "github.com/yujiawang-0/forum-page/internal/dataaccess"
	"github.com/yujiawang-0/forum-page/internal/database"
	"github.com/yujiawang-0/forum-page/internal/handlers/auth"
)
// TODO:
// enforcing that some posts belong to some topics (for every handler)


type PostHandler struct {
	DB *database.Database
}

// GET/posts
func (u *PostHandler) HandleGetAllPosts(w http.ResponseWriter, r *http.Request) {

	posts, err := postService.GetAllPosts(u.DB)
	if err != nil {
		api.MessageLogs.ErrorLog.Println(err)
		api.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	api.WriteJSON(w, http.StatusOK, api.Envelop{"posts": posts})
}

// GET/topics/{topic_id}/posts
func (u *PostHandler) HandleGetPostsByTopicID(w http.ResponseWriter, r *http.Request) {
	//parse the id
	topicId := chi.URLParam(r, "topic_id")
	topicIdInt, err := strconv.Atoi(topicId)
	if err != nil {
		api.MessageLogs.ErrorLog.Println(err)
		api.ErrorJSON(w, errors.New("invalid topic id"), http.StatusBadRequest)
		return 
	}
	posts, err := postService.GetPostsByTopicId(u.DB, topicIdInt)
	if err != nil {
		api.MessageLogs.ErrorLog.Println(err)
		api.ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	api.WriteJSON(w, http.StatusOK, api.Envelop{"posts": posts})

}

// GET/topics/{topic_id}/posts/{post_id}
func (u *PostHandler) HandleGetPostByID(w http.ResponseWriter, r *http.Request) {
	//parse the id
	postId := chi.URLParam(r, "post_id")
	postIdInt, err := strconv.Atoi(postId)
	if err != nil {
		api.MessageLogs.ErrorLog.Println(err)
		api.ErrorJSON(w, errors.New("invalid post id"), http.StatusBadRequest)
		return 
	}
	post, err := postService.GetPostByID(u.DB, postIdInt)
	if err != nil {
		api.MessageLogs.ErrorLog.Println(err)
		api.ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	api.WriteJSON(w, http.StatusOK, api.Envelop{"post": post})

}

// POST/topics/{topic_id}/posts
func (u *PostHandler) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	// cannot create post if there is no topic attached to it
	
	type CreatePostRequest struct {
		Title string  `json:"title"`
		Content string `json:"content"`
		// CreatorID int `json:"creator_id"`
	}
	
	topicID, err:= strconv.Atoi(chi.URLParam(r, "topic_id"))
	// should not normally hit this?
	if err != nil {
		api.MessageLogs.ErrorLog.Println(err)
		api.ErrorJSON(w, errors.New("invalid topic id"), http.StatusBadRequest)
		return 
	}
	
	var input CreatePostRequest

	if err := api.ReadJSON(w, r, &input); err != nil {
		api.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate if relevant fields are filled 
	if input.Title == "" || input.Content == "" {
		api.ErrorJSON(w, errors.New("title and content fields are required"), http.StatusBadRequest)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		api.ErrorJSON(w, errors.New("unauthenticated"), http.StatusUnauthorized)
		return
	}

	newPost:= models.Post{
		Title: input.Title,
		Content: input.Content,
		CreatorID: userID,
		TopicID: topicID,
	}

	post, err := postService.CreatePost(u.DB, newPost)
	if err != nil {
		if errors.Is(err, dataaccess.ErrTopicNotFound) {
			api.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		api.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	api.WriteJSON(w, http.StatusCreated, api.Envelop{"post": post})

}

// PUT/topics/{topic_id}/posts/{post_id}
func (u *PostHandler) HandleUpdatePost(w http.ResponseWriter, r *http.Request) {
	// you can update the title and content of the post without restriction 
	
	type UpdatePostRequest struct {
		Title string  `json:"title"`
		Content string `json:"content"`
	}
	
	//parse the topic id
	// TODO: later might need to make sure (for all post handlers) that every post is under the correct topic
	// need to make sure that there is no url tampering
	_, err := strconv.Atoi(chi.URLParam(r, "topic_id"))
	if err != nil {
		api.MessageLogs.ErrorLog.Println(err)
		api.ErrorJSON(w, errors.New("invalid topic id"), http.StatusBadRequest)
		return 
	}

	//parse the post id
	postIdInt, err := strconv.Atoi(chi.URLParam(r, "post_id"))
	if err != nil {
		api.MessageLogs.ErrorLog.Println(err)
		api.ErrorJSON(w, errors.New("invalid post id"), http.StatusBadRequest)
		return 
	}

	var input UpdatePostRequest

	if err := api.ReadJSON(w, r, &input); err != nil {
		api.ErrorJSON(w, err, http.StatusBadRequest)
		return
	} 
	
	titleTrimmed := strings.TrimSpace(input.Title)
	// validate if relevant fields are filled 
	if titleTrimmed == "" || input.Content == "" {
		api.ErrorJSON(w, errors.New("title and content cannot be empty"), http.StatusBadRequest)
		return
	}

	updatePost:= models.Post{
		Title: input.Title,
		Content: input.Content,
	}
	

	post, err := postService.UpdatePost(u.DB, postIdInt, updatePost)
	if err != nil {
		if errors.Is(err, dataaccess.ErrTopicNotFound) {
			api.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}
		
		api.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	api.WriteJSON(w, http.StatusOK, api.Envelop{"post": post})

}

// DELETE/topics/{topic_id}/posts/{post_id}
func (u *PostHandler) HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	//parse the id
	postId := chi.URLParam(r, "post_id")
	postIdInt, err := strconv.Atoi(postId)
	if err != nil {
		api.MessageLogs.ErrorLog.Println(err)
		api.ErrorJSON(w, errors.New("invalid post id"), http.StatusBadRequest)
		return 
	}

	err = postService.DeletePost(u.DB, postIdInt)
	if err != nil {
		api.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	api.WriteJSON(w, http.StatusOK, api.Envelop{"message": "post deleted successfully"})
}
