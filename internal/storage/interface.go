package storage

import (
	"OzonTestTask/internal/model"
	"context"
)

type PostStorage interface {
	CreatePost(ctx context.Context, post *model.Post) error
	GetAllPosts(ctx context.Context) ([]model.Post, error)
	GetPostByID(ctx context.Context, id string) (*model.Post, error)
}

type CommentStorage interface {
	CreateComment(ctx context.Context, comment *model.Comment) error
	GetCommentsByPost(ctx context.Context, postID string, limit, offset int) ([]model.Comment, error)
	GetReplies(ctx context.Context, parentCommentID string) ([]model.Comment, error)
}
