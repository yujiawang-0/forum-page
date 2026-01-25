package auth

// handles login, logout

import (
	//"encoding/json"
	//"fmt"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/yujiawang-0/forum-page/internal/api"

	// "github.com/yujiawang-0/forum-page/internal/handlers/auth"

	userService "github.com/yujiawang-0/forum-page/internal/dataaccess"
	"github.com/yujiawang-0/forum-page/internal/database"
)

type AuthHandler struct {
	DB *database.Database
}

// POST/auth/login
func (a *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var input Request

	// read JSON
	if err := api.ReadJSON(w, r, &input); err != nil {
		api.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate if relevant fields are filled 
	if input.Username == "" || input.Password == "" {
		api.ErrorJSON(w, errors.New("missing required fields"), http.StatusBadRequest)
		return
	}

	// fetch user by username
	user, err := userService.GetUserByUsernameForAuth(a.DB, input.Username)
	if err != nil {
		api.ErrorJSON(w, errors.New("invalid username"), http.StatusUnauthorized)
		return
	}

	fmt.Println("username:" + input.Username)
	fmt.Println("input password:" + input.Password)
	fmt.Println("DB password hash:" + user.Password)

	// check the password
	if !VerifyPassword(user.Password, input.Password) {
		api.ErrorJSON(w, errors.New("invalid password"), http.StatusUnauthorized)
		return
	}

	// create JWT
	token, err := createToken(user.ID, user.Role)
	if err != nil {
		api.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	api.WriteJSON(w, http.StatusOK, api.Envelop{
		"token": token,
		"user": user})

}

// handleLogout