package users

import (
	//"encoding/json"
	//"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/yujiawang-0/forum-page/internal/api"
	"github.com/yujiawang-0/forum-page/internal/dataaccess"
	"github.com/yujiawang-0/forum-page/internal/models"

	//"github.com/yujiawang-0/forum-page/internal/dataaccess"
	userService "github.com/yujiawang-0/forum-page/internal/dataaccess"
	"github.com/yujiawang-0/forum-page/internal/database"
)

// handlers: HTTP layer (parse JSON, URL, validate, return JSON)
// Should not contain SQL, pgx and do routing
// receives a call from routes and talks to services/dataaccess
// basically controllers in MVC
// what does the client need to do (to the database)?
// a handler is a function that takes a http request and writes a http response

// WITH REFERENCE TO:
// https://github.com/arturfil/yt-go-coffee-api-v2/blob/main/controllers/coffee.go
// https://github.com/dreamsofcode-io/golang-microservice-course-nn/blob/main/code/010-configuration-completion/handler/order.go


type UserHandler struct {
	DB *database.Database
}

// GET/users
func (u *UserHandler) HandleGetAllUsers(w http.ResponseWriter, r *http.Request) {

	users, err := userService.GetAllUsers(u.DB)
	if err != nil {
		api.MessageLogs.ErrorLog.Println(err)
		api.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	api.WriteJSON(w, http.StatusOK, api.Envelop{"users": users})
}


// GET/users/user/{id}
func (u *UserHandler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	//parse the id
	id := chi.URLParam(r, "id")
	id_int, err := strconv.Atoi(id)
	if err != nil {
		api.MessageLogs.ErrorLog.Println(err)
		api.ErrorJSON(w, errors.New("invalid user id"), http.StatusBadRequest)
		return 
	}
	user, err := userService.GetUserByID(u.DB, id_int)
	if err != nil {
		api.MessageLogs.ErrorLog.Println(err)
		api.ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	api.WriteJSON(w, http.StatusOK, api.Envelop{"user": user})

}

// POST/users/user
func (u *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	//TODO: cannot create user if username is already in database
	
	// creating a new user will always default them to the "user" role. 
	// updating to "admin" role requires the use of UpdateUser
	
	type CreateUserRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var input CreateUserRequest

	if err := api.ReadJSON(w, r, &input); err != nil {
		api.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate if relevant fields are filled 
	if input.Username == "" || input.Password == "" {
		api.ErrorJSON(w, errors.New("missing required fields"), http.StatusBadRequest)
		return
	}

	newUser:= models.User{
		Username: input.Username,
		Password: input.Password,
	}

	user, err := userService.CreateUser(u.DB, newUser)
	if err != nil {
		if errors.Is(err, dataaccess.ErrUsernameTaken) {
			api.ErrorJSON(w, err, http.StatusConflict)
			return
		}

		api.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	api.WriteJSON(w, http.StatusCreated, api.Envelop{"user": user})

}

// PUT/users/user/{id}
func (u *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	// when updating, frontend should send all fields, which will be sent to the database again. 
	
	type UpdateUserRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role string `json:"role"`
	}
	
	//parse the id
	id := chi.URLParam(r, "id")
	id_int, err := strconv.Atoi(id)
	if err != nil {
		api.MessageLogs.ErrorLog.Println(err)
		api.ErrorJSON(w, errors.New("invalid user id"), http.StatusBadRequest)
		return 
	}

	var input UpdateUserRequest

	if err := api.ReadJSON(w, r, &input); err != nil {
		api.ErrorJSON(w, err, http.StatusBadRequest)
		return
	} 

	// validate if relevant fields are filled 
	if input.Username == "" || input.Password == "" ||input.Role == "" {
		api.ErrorJSON(w, errors.New("username, password and role are required fields"), http.StatusBadRequest)
		return
	}

	

	updateUser:= models.User{
		Username: input.Username,
		Password: input.Password,
		Role: input.Role,
	}

	user, err := userService.UpdateUser(u.DB, id_int, updateUser)
	if err != nil {
		if errors.Is(err, dataaccess.ErrUsernameTaken) {
			api.ErrorJSON(w, err, http.StatusConflict)
			return
		}
		
		api.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	api.WriteJSON(w, http.StatusOK, api.Envelop{"user": user})

}

// DELETE/users/user/{id}
func (u *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	//parse the id
	id := chi.URLParam(r, "id")
	id_int, err := strconv.Atoi(id)
	if err != nil {
		api.MessageLogs.ErrorLog.Println(err)
		api.ErrorJSON(w, errors.New("invalid user id"), http.StatusBadRequest)
		return 
	}

	err = userService.DeleteUser(u.DB, id_int)
	if err != nil {
		api.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	api.WriteJSON(w, http.StatusOK, api.Envelop{"message": "user deleted successfully"})
}
