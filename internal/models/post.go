package models

import (
	"time"
)

// models: structs shared across layers

type Post struct {
	PostID   	int    		`json:"userid"`
	CreatorID 	int 		`json:"creatorid"`
	TopicID 	int			`json:"topicid"`
	Title 	string 			`json:"title"`
	Content 	string		`json:"content"`
	DateCreated time.Time	`json:"date_created"`
	DateUpdated time.Time	`json:"date_updated"`
}


