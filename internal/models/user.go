package models

import (
	"fmt"
	"time"
)

// models: structs shared across layers

type User struct {
	ID   		int    			`json:"id"`
	Username 	string 		`json:"username"`
	Password 	string		`json:"-"` // not sent to the client
	Role 		string		`json:"role"`
	DateCreated time.Time	`json:"date_created"`
}

func (user *User) Greet() string {
	return fmt.Sprintf("Hello, I am %s", user.Username)
}
