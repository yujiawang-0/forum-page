package models

import (
	//"fmt"
	"time"
)

// models: structs shared across layers

type User struct {
	ID   		int    			`json:"user_id"`
	Username 	string 		`json:"username"`
	Password 	string		`json:"-"` // not sent to the client
	Role 		string		`json:"role"`
	IsActive 	bool 		`json:"is_active"`
	DateCreated time.Time	`json:"date_created"`
}

