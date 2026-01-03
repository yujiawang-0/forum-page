package dataaccess

import (
	"context"
	"time"

	"github.com/yujiawang-0/forum-page/internal/database"
	"github.com/yujiawang-0/forum-page/internal/models"
)

//other notes:
// dataaccess is used for DB queries (CRUD)
// like a services folder, talks to the database directly
// returns go structs
// references: https://github.com/arturfil/yt-go-coffee-api-v2/blob/main/services/coffee.go

// type User struct {
// 	ID   		int    			`json:"user_id"`
// 	Username 	string 		`json:"username"`
// 	Password 	string		`json:"-"` // not sent to the client
// 	Role 		string		`json:"role"`
// 	DateCreated time.Time	`json:"date_created"`
// }

const dbTimeout = time.Second * 3

// Reading
func GetAllUsers(db *database.Database) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT user_id, username, role, date_created FROM users`
	rows, err := db.Conn.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Role,
			&user.DateCreated,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)

	}
	// remember that it returns a pointer to the slice, and not the slice itself
	return users, nil
}

func GetUserByID(db *database.Database, id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT user_id, username, role, date_created FROM users WHERE user_id = $1`
	row := db.Conn.QueryRow(ctx, query, id)

	var user models.User
	
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Role,
		&user.DateCreated,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Writing

func CreateUser(db *database.Database, user models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `INSERT INTO users (username, password, role, date_created) 
	VALUES ($1, $2, $3, $4) returning * `

	var returnedUser models.User

	err := db.Conn.QueryRow(
		ctx, 
		query,
		user.Username,
		user.Password,
		user.Role,
		time.Now(), // date_created is now
	).Scan(
		&returnedUser.ID,
		&returnedUser.Username,
		&returnedUser.Password, //TODO: remember to remove this and update SQL returning statement 
		&returnedUser.Role,	
		&returnedUser.DateCreated,
	)
	if err != nil {
		return nil, err
	}

	return &returnedUser, nil
}

func UpdateUser( db *database.Database, id int, body models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
        UPDATE coffees
        SET
            username = $1,
            password = $2,
            role = $3,
        WHERE user_id = $4
        returning *
		`
	var returnedUser models.User

	err := db.Conn.QueryRow(
		ctx, 
		query,
		body.Username,
		body.Password,
		body.Role,
		id,
	).Scan(
		&returnedUser.ID,
		&returnedUser.Username,
		&returnedUser.Password, //TODO: remember to remove this and update SQL returning statement 
		&returnedUser.Role,
		&returnedUser.DateCreated,
	)
	if err != nil {
		return nil, err
	}

	return &returnedUser, nil

}

func DeleteUser(db *database.Database, id int) (error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM users WHERE user_id = $1`
	_, err := db.Conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}