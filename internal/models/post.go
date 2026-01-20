package models

import (
	"time"
)

// models: structs shared across layers

type Post struct {
	PostID   	int    		`json:"post_id"`
	CreatorID 	int 		`json:"creator_id"`
	TopicID 	int			`json:"topic_id"`
	Title 		string 		`json:"title"`
	Content 	string		`json:"content"`
	DateCreated time.Time	`json:"date_created"`
	DateUpdated time.Time	`json:"date_updated"`
}


