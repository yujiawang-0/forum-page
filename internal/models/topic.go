package models

import (
	"time"
)

// models: structs shared across layers

type Topic struct {
	TopicID 	int			`json:"topicid"`
	TopicName 	string 		`json:"topicname"`
	AdminID   	int    		`json:"adminid"`
	CreatorID 	int 		`json:"creatorid"`
	DateCreated time.Time	`json:"date_created"`
}
