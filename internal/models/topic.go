package models

import (
	"time"
)

// models: structs shared across layers
// adminID can be null (there might not be an admin for the topic)

type Topic struct {
	TopicID 	int			`json:"topic_id"`
	TopicName 	string 		`json:"topic_name"`
	AdminID   	int    		`json:"admin_id"`
	CreatorID 	int 		`json:"creator_id"`
	DateCreated time.Time	`json:"date_created"`
}
