package model

import "time"

type Post struct {
	ID                 int       `json:"id"`
	Title              string    `json:"title"`
	Content            string    `json:"content"`
	Author             string    `json:"author"`
	AreCommentsAllowed bool      `json:"areCommentsAllowed"`
	CreatedAt          time.Time `json:"createdAt"`
}
