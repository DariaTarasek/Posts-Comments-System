package comment

import (
	"OzonTestTask/internal/model"
	"OzonTestTask/internal/storage"
	"context"
)

type CommentService struct {
	store storage.CommentStorage
}

func NewCommentService(s storage.CommentStorage) *CommentService {
	return &CommentService{store: s}
}

func (s *CommentService) CreateComment(ctx context.Context, comment *model.Comment) error {
	return s.store.CreateComment(ctx, comment)
}

func (s *CommentService) GetCommentsByPost(ctx context.Context, postID string, limit, offset int) ([]model.Comment, error) {
	return s.store.GetCommentsByPost(ctx, postID, limit, offset)
}

func (s *CommentService) GetReplies(ctx context.Context, parentCommentID string) ([]model.Comment, error) {
	return s.store.GetReplies(ctx, parentCommentID)
}
