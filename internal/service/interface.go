package service

import (
	"OzonTestTask/internal/model"
	"context"
)

type CommentService interface {
	CreateComment(ctx context.Context, comment *model.Comment) error
	GetCommentsByPost(ctx context.Context, postID int, limit, offset int) ([]model.Comment, int, error)
	GetReplies(ctx context.Context, parentCommentID int) ([]model.Comment, error)
	GetPostByID(ctx context.Context, id int) (*model.Post, error)
}

type PostService interface {
	CreatePost(ctx context.Context, post *model.Post) error
	GetAllPosts(ctx context.Context) ([]model.Post, error)
	GetPostByID(ctx context.Context, id int) (*model.Post, error)
}
