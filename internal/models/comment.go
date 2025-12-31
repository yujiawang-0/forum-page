package models

import (
	"time"
)

// models: structs shared across layers

type Comment struct {
	CommentID 	int 		`json:"commentid"`
	PostID   	int    		`json:"userid"`
	CreatorID 	int 		`json:"creatorid"`
	TopicID 	int			`json:"topicid"`
	ReplyCommentID int		`json:"reply_comment_id"`
	Content 	string		`json:"content"`
	DateCreated time.Time	`json:"date_created"`
	DateUpdated time.Time	`json:"date_updated"`
}
