package resolvers

import (
	"OzonTestTask/internal/service/comment"
	"OzonTestTask/internal/service/post"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PostService    *post.PostService
	CommentService *comment.CommentService
}
