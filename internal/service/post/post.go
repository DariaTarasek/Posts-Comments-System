package post

import (
	"OzonTestTask/internal/model"
	"OzonTestTask/internal/storage"
	"context"
	"fmt"
)

type PostService struct {
	store storage.PostStorage
}

func NewPostService(s storage.PostStorage) *PostService {
	return &PostService{store: s}
}

func (s *PostService) CreatePost(ctx context.Context, post *model.Post) error {
	if post.Title == "" {
		return fmt.Errorf("заголовок поста не может быть пустым")
	}
	if post.Content == "" {
		return fmt.Errorf("пост не может быть пустым")
	}
	if post.Author == "" {
		return fmt.Errorf("имя автора не может быть пустым")
	}
	err := s.store.CreatePost(ctx, post)
	if err != nil {
		return fmt.Errorf("не удалось создать пост: %v", err)
	}
	return nil
}

func (s *PostService) GetAllPosts(ctx context.Context) ([]model.Post, error) {
	posts, err := s.store.GetAllPosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список постов: %v", err)
	}
	return posts, nil
}

func (s *PostService) GetPostByID(ctx context.Context, id int) (*model.Post, error) {
	post, err := s.store.GetPostByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить пост: %v", err)
	}
	return post, nil
}
