package model

import "time"

type Comment struct {
	ID              int       `json:"id" db:"id"`
	PostID          int       `json:"postId" db:"post_id"`
	ParentCommentID *int      `json:"parentCommentId,omitempty" db:"parent_comment_id,omitempty"`
	Path            string    `json:"path" db:"path"`
	Author          string    `json:"author" db:"author"`
	Content         string    `json:"content" db:"content"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
}
