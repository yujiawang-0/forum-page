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

// CRUD for users

const dbTimeout = time.Second * 3

// Reading
func GetAllUsers(db *database.Database) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT user_id, username, role, date_created, is_active FROM users`
	rows, err := db.Conn.Query(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Role,
			&user.DateCreated,
			&user.IsActive,
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

	query := `SELECT user_id, username, role, date_created, is_active FROM users WHERE user_id = $1`
	row := db.Conn.QueryRow(ctx, query, id)

	var user models.User
	
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Role,
		&user.DateCreated,
		&user.IsActive,
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
	query := `INSERT INTO users (username, password)
	VALUES ($1, $2) returning user_id, username, role, date_created, is_active`

	var returnedUser models.User

	// TODO: remember to hash password 
	err := db.Conn.QueryRow(
		ctx, 
		query,
		user.Username,
		user.Password, 
	).Scan(
		&returnedUser.ID,
		&returnedUser.Username,
		&returnedUser.Role,	
		&returnedUser.DateCreated,
		&returnedUser.IsActive,
	)
	if err != nil {
		return nil, UniqueUsernameViolation(err)
	}

	return &returnedUser, nil
}

func UpdateUser(db *database.Database, id int, body models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// coalesce: take the first non-null value out of all arguments
	query := `
        UPDATE users
        SET
            username = COALESCE($1, username), 
            password = COALESCE($2, password),
            role = COALESCE($3, role)
        WHERE user_id = $4
        returning user_id, username, role, date_created, is_active
		`
	var returnedUser models.User

	// TODO: remember to hash password
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
		&returnedUser.Role,
		&returnedUser.DateCreated,
		&returnedUser.IsActive,
	)
	if err != nil {
		return nil, UniqueUsernameViolation(err)
	}

	return &returnedUser, nil

}

// TODO: Edit DeleteUser such that it only makes the is_active flag FALSE. users should not be able to be deleted by the client, only from the DMS
func DeleteUser(db *database.Database, id int) (error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `UPDATE users
		SET is_active = FALSE
		WHERE user_id = $1`
	_, err := db.Conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}