package model

import "time"

type Comment struct {
	ID              int       `json:"id"`
	PostID          int       `json:"postId"`
	ParentCommentID *int      `json:"parentCommentId,omitempty"`
	Path            string    `json:"path"`
	Author          string    `json:"author"`
	Content         string    `json:"content"`
	CreatedAt       time.Time `json:"createdAt"`
}
