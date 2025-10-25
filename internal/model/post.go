package model

import "time"

type Post struct {
	ID                 int       `json:"id" db:"id"`
	Title              string    `json:"title" db:"title"`
	Content            string    `json:"content" db:"content"`
	Author             string    `json:"author" db:"author"`
	AreCommentsAllowed bool      `json:"areCommentsAllowed" db:"are_comments_allowed"`
	CreatedAt          time.Time `json:"createdAt" db:"created_at"`
}
