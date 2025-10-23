package model

import "time"

type Comment struct {
	ID              string    `json:"id"`
	PostID          string    `json:"postId"`
	ParentCommentID *string   `json:"parentCommentId,omitempty"`
	Path            string    `json:"path"`
	Author          string    `json:"author"`
	Content         string    `json:"content"`
	CreatedAt       time.Time `json:"createdAt"`
}
