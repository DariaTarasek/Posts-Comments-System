package model

import "time"

type Comment struct {
	ID              int       `json:"id" db:"id"`
	PostID          int       `json:"post_id" db:"post_id"`
	ParentCommentID *int      `json:"parent_comment_id,omitempty" db:"parent_comment_id,omitempty"`
	Path            string    `json:"path" db:"path"`
	Author          string    `json:"author" db:"author"`
	Content         string    `json:"content" db:"content"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type PaginatedComments struct {
	Comments   []*Comment `json:"comments"`
	TotalPages int        `json:"totalPages"`
}
