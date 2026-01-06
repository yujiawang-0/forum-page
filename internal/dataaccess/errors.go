package dataaccess

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrUsernameTaken = errors.New("username already exists")


func UniqueUsernameViolation(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return ErrUsernameTaken
		}
	}
	return err
}
