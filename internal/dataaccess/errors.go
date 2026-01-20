package dataaccess

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrUsernameTaken = errors.New("username already exists")

var ErrUserNotFound = errors.New("referenced user does not exist")

var ErrTopicNameTaken = errors.New("topic already exists")

var ErrTopicNotFound = errors.New("topic not found")


func UniqueUsernameViolation(err error) error {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" { // UniqueViolation
			return ErrUsernameTaken
		}
	}
	return err
}

func TranslateTopicError(err error) error { // check if username exists for adminID and creatorID
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == "23503" { // ForeignKeyViolation
			return ErrUserNotFound
		}
		if pgErr.Code == "23505" { // UniqueViolation
			return ErrTopicNameTaken
		}
	}
	return err
}

func TranslatePostError(err error) error { // check if username exists for adminID and creatorID
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == "23503" { // ForeignKeyViolation
			return ErrTopicNotFound
		}
	}
	return err
}

