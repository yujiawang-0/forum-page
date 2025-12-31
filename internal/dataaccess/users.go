package dataaccess

import (
	"github.com/yujiawang-0/forum-page/internal/database"
	"github.com/yujiawang-0/forum-page/internal/models"
)

//other notes:
// dataaccess is used for DB queries (CRUD)
// like a services folder, talks to the database directly


func ListUsers(db *database.Database) ([]models.User, error) {
	users := []models.User{
		{
			ID:   1,
			Username: "CVWO",
		},
	}
	return users, nil
}

func GetUserByID() {

}

func CreateUser() {
	
}