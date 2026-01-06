package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/yujiawang-0/forum-page/internal/api"

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



const (
	GetAllUsers = "users.GetAllUsers"
	GetUserByID = "users.GetUserByID"

	SuccessfulListUsersMessage = "Successfully listed users"
	ErrParseStrToInt = "unable to parse userid string to int"
	// ErrRetrieveDatabase        = "Failed to retrieve database in %s"
	ErrRetrieveUsers           = "Failed to retrieve users in %s"
	ErrEncodeView              = "Failed to retrieve users in %s"
)

type UserHandler struct {
	DB *database.Database
}

// GET/users
func (u *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) (*api.Response, error) {

	users, err := userService.GetAllUsers(u.DB)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrRetrieveUsers, GetAllUsers))
	}

	data, err := json.Marshal(users)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrEncodeView, GetAllUsers))
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{SuccessfulListUsersMessage},
	}, nil
}


// func (u *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	
// }

// GET/users/{id}
func (u *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	id := chi.URLParam(r, "id")
	id_int, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrParseStrToInt, GetAllUsers))
	}
	users, err := userService.GetUserByID(u.DB, id_int)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrRetrieveUsers, GetAllUsers))
	}

	data, err := json.Marshal(users)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrEncodeView, GetAllUsers))
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{SuccessfulListUsersMessage},
	}, nil

}
func (u *UserHandler) CreateBook(w http.ResponseWriter, r *http.Request) {}
func (u *UserHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {}
func (u *UserHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {}
