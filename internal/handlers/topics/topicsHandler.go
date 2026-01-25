package topics

import (
	//"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/yujiawang-0/forum-page/internal/api"
	"github.com/yujiawang-0/forum-page/internal/dataaccess"
	"github.com/yujiawang-0/forum-page/internal/models"

	topicService "github.com/yujiawang-0/forum-page/internal/dataaccess"
	"github.com/yujiawang-0/forum-page/internal/database"
	"github.com/yujiawang-0/forum-page/internal/handlers/auth"
)


type TopicHandler struct {
	DB *database.Database
}

// GET/topics
func (u *TopicHandler) HandleGetAllTopics(w http.ResponseWriter, r *http.Request) {

	topics, err := topicService.GetAllTopics(u.DB)
	if err != nil {
		api.MessageLogs.ErrorLog.Println(err)
		api.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	api.WriteJSON(w, http.StatusOK, api.Envelop{"topics": topics})
}


// GET/topics/{topic_id}
func (u *TopicHandler) HandleGetTopicByID(w http.ResponseWriter, r *http.Request) {
	//parse the id
	id := chi.URLParam(r, "topic_id")
	id_int, err := strconv.Atoi(id)
	if err != nil {
		api.MessageLogs.ErrorLog.Println(err)
		api.ErrorJSON(w, errors.New("invalid topic id"), http.StatusBadRequest)
		return 
	}
	topic, err := topicService.GetTopicByID(u.DB, id_int)
	if err != nil {
		api.MessageLogs.ErrorLog.Println(err)
		api.ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	api.WriteJSON(w, http.StatusOK, api.Envelop{"topic": topic})

}

// POST/topics
func (u *TopicHandler) HandleCreateTopic(w http.ResponseWriter, r *http.Request) {
	// cannot create topic if topic_name is already in database
	// cannot assign admin when creating topic, can only do so when updating the topic
	// creator of topic is the default admin of the topic
	
	type CreateTopicRequest struct {
		TopicName string `json:"topic_name"`
		// CreatorID int `json:"creator_id"`
	}

	var input CreateTopicRequest

	if err := api.ReadJSON(w, r, &input); err != nil {
		api.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate if relevant fields are filled 
	if input.TopicName == "" {
		api.ErrorJSON(w, errors.New("topic name required"), http.StatusBadRequest)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		api.ErrorJSON(w, errors.New("unauthenticated"), http.StatusUnauthorized)
		return
	}

	newTopic:= models.Topic{
		TopicName: input.TopicName,
		AdminID: userID,
		CreatorID: userID,
	}

	topic, err := topicService.CreateTopic(u.DB, newTopic)
	if err != nil {
		if errors.Is(err, dataaccess.ErrTopicNameTaken) {
			api.ErrorJSON(w, err, http.StatusConflict)
			return
		}
		if errors.Is(err, dataaccess.ErrUserNotFound) {
			api.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		api.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	api.WriteJSON(w, http.StatusCreated, api.Envelop{"topic": topic})

}

// PUT/topics/{topic_id}
func (u *TopicHandler) HandleUpdateTopic(w http.ResponseWriter, r *http.Request) {
	// you can update admin herer
	// you can update the name of the topic, provided it is not already in use in the db 
	
	type UpdateTopicRequest struct {
		TopicName string `json:"topic_name"`
		AdminID int `json:"admin_id"`
	}
	
	//parse the id
	id := chi.URLParam(r, "topic_id")
	id_int, err := strconv.Atoi(id)
	if err != nil {
		api.MessageLogs.ErrorLog.Println(err)
		api.ErrorJSON(w, errors.New("invalid topic id"), http.StatusBadRequest)
		return 
	}

	var input UpdateTopicRequest

	if err := api.ReadJSON(w, r, &input); err != nil {
		api.ErrorJSON(w, err, http.StatusBadRequest)
		return
	} 
	
	topicNameTrimmed := strings.TrimSpace(input.TopicName)
	// validate if relevant fields are filled 
	if input.TopicName == "" {
		api.ErrorJSON(w, errors.New("topic name cannot be empty"), http.StatusBadRequest)
		return
	}

	updateTopic:= models.Topic{
		TopicName: topicNameTrimmed,
		AdminID: input.AdminID,
	}

	topic, err := topicService.UpdateTopic(u.DB, id_int, updateTopic)
	if err != nil {
		if errors.Is(err, dataaccess.ErrTopicNameTaken) {
			api.ErrorJSON(w, err, http.StatusConflict)
			return
		}
		if errors.Is(err, dataaccess.ErrUserNotFound) {
			api.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}
		
		api.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	api.WriteJSON(w, http.StatusOK, api.Envelop{"topic": topic})

}

// DELETE/topics/{topic_id}
func (u *TopicHandler) HandleDeleteTopic(w http.ResponseWriter, r *http.Request) {
	//parse the id
	id := chi.URLParam(r, "topic_id")
	id_int, err := strconv.Atoi(id)
	if err != nil {
		api.MessageLogs.ErrorLog.Println(err)
		api.ErrorJSON(w, errors.New("invalid topic id"), http.StatusBadRequest)
		return 
	}

	err = topicService.DeleteTopic(u.DB, id_int)
	if err != nil {
		api.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	api.WriteJSON(w, http.StatusOK, api.Envelop{"message": "topic deleted successfully"})
}
