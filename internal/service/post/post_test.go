package post

import (
	"OzonTestTask/internal/model"
	"OzonTestTask/internal/storage/in-memory"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	postService *PostService
	ctx         context.Context
)

// в слое бизнес-логики тестирую только с in-memory хранилищем, т.к.
// методы бизнес-логики не зависят от типа хранилища
// также тестирую только работу валидаций,
// т.к. в позитивных случаях - просто проброс в слой работы с хранилищем (уже протестировано в /storage)
func TestMain(m *testing.M) {
	store := in_memory.NewInMemoryStorage()
	postService = NewPostService(store)
	ctx = context.Background()
	m.Run()
}

func TestCreatePost_EmptyTitle(t *testing.T) {
	post := &model.Post{
		Title:              "",
		Content:            "Содержимое",
		Author:             "Автор",
		AreCommentsAllowed: true,
	}
	err := postService.CreatePost(ctx, post)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "заголовок поста не может быть пустым")
}

func TestCreatePost_EmptyContent(t *testing.T) {
	post := &model.Post{
		Title:              "Заголовок",
		Content:            "",
		Author:             "Автор",
		AreCommentsAllowed: true,
	}
	err := postService.CreatePost(ctx, post)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "пост не может быть пустым")
}

func TestCreatePost_EmptyAuthor(t *testing.T) {
	post := &model.Post{
		Title:              "Заголовок",
		Content:            "Контент",
		Author:             "",
		AreCommentsAllowed: true,
	}
	err := postService.CreatePost(ctx, post)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "имя автора не может быть пустым")
}
