package comment

import (
	"OzonTestTask/internal/model"
	"OzonTestTask/internal/service/post"
	"OzonTestTask/internal/storage/in-memory"
	"OzonTestTask/internal/subscription"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	commentService *CommentService
	postService    *post.PostService
	subService     subscription.Subscription
	ctx            context.Context
)

// в слое бизнес-логики тестирую только с in-memory хранилищем, т.к.
// методы бизнес-логики не зависят от типа хранилища
// также тестирую только работу валидаций,
// т.к. в позитивных случаях - просто проброс в слой работы с хранилищем (уже протестировано в /storage)
func TestMain(m *testing.M) {
	store := in_memory.NewInMemoryStorage()
	commentService = NewCommentService(store, subService)
	postService = post.NewPostService(store)
	ctx = context.Background()
	m.Run()
}

func TestCreateComment_EmptyContent(t *testing.T) {
	p := &model.Post{
		Title:              "Название",
		Content:            "Содержимое",
		Author:             "Автор",
		AreCommentsAllowed: true,
	}
	err := postService.CreatePost(ctx, p)

	comment := &model.Comment{
		PostID:  1,
		Author:  "Новый автор",
		Content: "",
	}
	err = commentService.CreateComment(ctx, comment)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "комментарий не может быть пустым")
}

func TestCreateComment_EmptyAuthor(t *testing.T) {
	p := &model.Post{
		Title:              "Название",
		Content:            "Содержимое",
		Author:             "Автор",
		AreCommentsAllowed: true,
	}
	err := postService.CreatePost(ctx, p)

	comment := &model.Comment{
		PostID:  1,
		Author:  "",
		Content: "Текст",
	}
	err = commentService.CreateComment(ctx, comment)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "имя автора не может быть пустым")
}

func TestCreateComment_TooLongComment(t *testing.T) {
	p := &model.Post{
		Title:              "Название",
		Content:            "Содержимое",
		Author:             "Автор",
		AreCommentsAllowed: true,
	}
	err := postService.CreatePost(ctx, p)

	contentRunes := make([]rune, 2001)
	content := string(contentRunes)
	comment := &model.Comment{
		PostID:  1,
		Author:  "Новый автор",
		Content: content,
	}
	err = commentService.CreateComment(ctx, comment)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "длина комментария не должна превышать 2000 символов")
}

func TestCreateComment_CommentsNotAllowed(t *testing.T) {
	p := &model.Post{
		Title:              "Название",
		Content:            "Содержимое",
		Author:             "Автор",
		AreCommentsAllowed: false,
	}
	err := postService.CreatePost(ctx, p)

	comment := &model.Comment{
		PostID:  1,
		Author:  "Новый автор",
		Content: "Текст",
	}
	err = commentService.CreateComment(ctx, comment)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "этот пост запрещено комментировать")
}
