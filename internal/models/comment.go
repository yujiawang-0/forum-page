package models

import (
	"time"
)

// models: structs shared across layers

type Comment struct {
	CommentID 	int 		`json:"comment_id"`
	PostID   	int    		`json:"user_id"`
	CreatorID 	int 		`json:"creator_id"`
	TopicID 	int			`json:"topic_id"`
	ReplyCommentID int		`json:"reply_comment_id"`
	Content 	string		`json:"content"`
	DateCreated time.Time	`json:"date_created"`
	DateUpdated time.Time	`json:"date_updated"`
}
