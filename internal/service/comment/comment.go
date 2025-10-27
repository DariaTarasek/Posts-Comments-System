package comment

import (
	"OzonTestTask/internal/model"
	"OzonTestTask/internal/storage"
	"OzonTestTask/internal/subscription"
	"context"
	"fmt"
)

type CommentService struct {
	store storage.CommentStorage
	sub   subscription.Subscription
}

func NewCommentService(store storage.CommentStorage, sub subscription.Subscription) *CommentService {
	return &CommentService{
		store: store,
		sub:   sub,
	}
}

func (s *CommentService) GetPostByID(ctx context.Context, id int) (*model.Post, error) {
	post, err := s.store.GetPostByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить пост: %v", err)
	}
	return post, nil
}

func (s *CommentService) CreateComment(ctx context.Context, comment *model.Comment) error {
	post, err := s.store.GetPostByID(ctx, comment.PostID)
	if err != nil {
		return fmt.Errorf("пост не найден: %v", err)
	}

	if !post.AreCommentsAllowed {
		return fmt.Errorf("этот пост запрещено комментировать")
	}

	if comment.Content == "" {
		return fmt.Errorf("комментарий не может быть пустым")
	}
	if comment.Author == "" {
		return fmt.Errorf("имя автора не может быть пустым")
	}

	if len(comment.Content) > 2000 {
		return fmt.Errorf("длина комментария не должна превышать 2000 символов")
	}

	err = s.store.CreateComment(ctx, comment)
	if err != nil {
		return fmt.Errorf("не удалось создать комментарий: %v", err)
	}

	if s.sub != nil {
		err = s.sub.Publish(comment.PostID, comment)
		if err != nil {
			return fmt.Errorf("не удалось отправить уведомление о новом комментарии: %v", err)
		}
	}

	return nil
}

func (s *CommentService) GetCommentsByPost(ctx context.Context, postID int, limit, offset int) ([]model.Comment, int, error) {
	comments, totalComments, err := s.store.GetCommentsByPost(ctx, postID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("не удалось получить корневые комментарии: %v", err)
	}
	// отдаю общее количество страниц - инфа для клиента
	totalPages := (totalComments + limit - 1) / limit
	return comments, totalPages, nil

}

func (s *CommentService) GetReplies(ctx context.Context, parentCommentID int) ([]model.Comment, error) {
	replies, err := s.store.GetReplies(ctx, parentCommentID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить вложенные комментарии: %v", err)
	}
	return replies, nil
}
