package post

import (
	"OzonTestTask/internal/model"
	"OzonTestTask/internal/storage"
	"context"
)

type PostService struct {
	store storage.PostStorage
}

func NewPostService(s storage.PostStorage) *PostService {
	return &PostService{store: s}
}

func (s *PostService) CreatePost(ctx context.Context, post *model.Post) error {
	return s.store.CreatePost(ctx, post)
}

func (s *PostService) GetAllPosts(ctx context.Context) ([]model.Post, error) {
	return s.store.GetAllPosts(ctx)
}

func (s *PostService) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	return s.store.GetPostByID(ctx, id)
}
