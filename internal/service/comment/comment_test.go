package comment

import (
	"OzonTestTask/internal/mocks"
	"OzonTestTask/internal/model"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var ctx = context.Background()

func TestCreateComment_EmptyContent(t *testing.T) {
	mockStorage := new(mocks.CommentStorage)
	commentService := NewCommentService(mockStorage, nil)

	comment := &model.Comment{
		PostID:  1,
		Author:  "Автор",
		Content: "",
	}

	post := &model.Post{
		ID:                 1,
		Title:              "Название",
		Content:            "Содержимое",
		Author:             "Автор",
		AreCommentsAllowed: true,
	}

	mockStorage.On("GetPostByID", mock.Anything, post.ID).Return(post, nil)

	err := commentService.CreateComment(ctx, comment)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "комментарий не может быть пустым")
}

func TestCreateComment_EmptyAuthor(t *testing.T) {
	mockStorage := new(mocks.CommentStorage)
	commentService := NewCommentService(mockStorage, nil)

	comment := &model.Comment{
		PostID:  1,
		Author:  "",
		Content: "Текст",
	}

	post := &model.Post{
		ID:                 1,
		Title:              "Название",
		Content:            "Содержимое",
		Author:             "Автор",
		AreCommentsAllowed: true,
	}

	mockStorage.On("GetPostByID", mock.Anything, post.ID).Return(post, nil)

	err := commentService.CreateComment(ctx, comment)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "имя автора не может быть пустым")
}

func TestCreateComment_TooLongComment(t *testing.T) {
	mockStorage := new(mocks.CommentStorage)
	commentService := NewCommentService(mockStorage, nil)

	contentRunes := make([]rune, 2001)
	comment := &model.Comment{
		PostID:  1,
		Author:  "Автор",
		Content: string(contentRunes),
	}

	post := &model.Post{
		ID:                 1,
		Title:              "Название",
		Content:            "Содержимое",
		Author:             "Автор",
		AreCommentsAllowed: true,
	}

	mockStorage.On("GetPostByID", mock.Anything, post.ID).Return(post, nil)

	err := commentService.CreateComment(ctx, comment)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "длина комментария не должна превышать 2000 символов")
}

func TestCreateComment_CommentsNotAllowed(t *testing.T) {
	mockStorage := new(mocks.CommentStorage)
	commentService := NewCommentService(mockStorage, nil)

	post := &model.Post{
		ID:                 1,
		Title:              "Название",
		Content:            "Содержимое",
		Author:             "Автор",
		AreCommentsAllowed: false,
	}

	mockStorage.On("GetPostByID", mock.Anything, post.ID).Return(post, nil)

	comment := &model.Comment{
		PostID:  1,
		Author:  "Автор",
		Content: "Текст",
	}

	err := commentService.CreateComment(ctx, comment)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "этот пост запрещено комментировать")

	mockStorage.AssertExpectations(t)
}
