package in_memory

import (
	"OzonTestTask/internal/model"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	storage *InMemoryStorage
	ctx     context.Context
)

func conf() {
	storage = NewInMemoryStorage()
	ctx = context.Background()
}

func TestCreateAndGetPost(t *testing.T) {
	conf()
	post := &model.Post{
		Title:              "Тестовый пост",
		Content:            "Содержимое",
		Author:             "Даша",
		AreCommentsAllowed: true,
	}
	require.NoError(t, storage.CreatePost(ctx, post), "пост не создан")
	created, err := storage.GetPostByID(ctx, post.ID)
	require.NoError(t, err, "пост не найден")
	assert.Equal(t, post.Title, created.Title)
	assert.Equal(t, post.Author, created.Author)
}

func TestGetAllPosts(t *testing.T) {
	conf()
	post1 := &model.Post{
		Title:   "Первый",
		Content: "Текст",
		Author:  "Даша",
	}
	post2 := &model.Post{
		Title:   "Второй",
		Content: "Текст",
		Author:  "Аня",
	}

	require.NoError(t, storage.CreatePost(ctx, post1), "пост не создан")
	time.Sleep(time.Millisecond)
	require.NoError(t, storage.CreatePost(ctx, post2), "пост не создан")

	posts, err := storage.GetAllPosts(ctx)
	require.NoError(t, err)
	assert.Len(t, posts, 2)
	assert.GreaterOrEqualf(t, posts[0].CreatedAt.UnixMilli(), posts[1].CreatedAt.UnixMilli(), "должен сортироваться по дате создания")
}

func TestGetPostByID_WrongID(t *testing.T) {
	conf()
	post, err := storage.GetPostByID(ctx, -1)
	assert.Nil(t, post)
	assert.Error(t, err)
}

func TestCreateComment(t *testing.T) {
	conf()
	post := &model.Post{
		Title:              "Пост",
		AreCommentsAllowed: true,
	}
	require.NoError(t, storage.CreatePost(ctx, post), "пост не создан")

	comment := &model.Comment{
		PostID:  post.ID,
		Author:  "Автор",
		Content: "Комментарий",
	}
	require.NoError(t, storage.CreateComment(ctx, comment), "комментарий не создан")
	assert.NotZero(t, comment.ID)
}

func TestCreateComment_WrongPostID(t *testing.T) {
	conf()

	comment := &model.Comment{
		PostID:  -1,
		Author:  "Дарья",
		Content: "ААА",
	}
	err := storage.CreateComment(ctx, comment)
	assert.Error(t, err)
}

func TestGetCommentsByPost(t *testing.T) {
	conf()

	post := &model.Post{
		Title:              "Пост",
		AreCommentsAllowed: true,
	}
	require.NoError(t, storage.CreatePost(ctx, post), "пост не создан")

	root := &model.Comment{
		PostID:  post.ID,
		Content: "Корневой",
	}
	require.NoError(t, storage.CreateComment(ctx, root), "комментарий не создан")

	reply := &model.Comment{
		PostID:          post.ID,
		ParentCommentID: &root.ID,
		Content:         "Ответ",
	}
	require.NoError(t, storage.CreateComment(ctx, reply), "комментарий не создан")

	rootComments, total, err := storage.GetCommentsByPost(ctx, post.ID, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Equal(t, root.ID, rootComments[0].ID)
}

func TestGetReplies(t *testing.T) {
	conf()
	post := &model.Post{
		Title:              "Пост",
		AreCommentsAllowed: true,
	}
	require.NoError(t, storage.CreatePost(ctx, post), "пост не создан")

	root := &model.Comment{
		PostID: post.ID,
	}
	require.NoError(t, storage.CreateComment(ctx, root), "комментарий не создан")

	c1 := &model.Comment{
		PostID:          post.ID,
		ParentCommentID: &root.ID,
	}
	require.NoError(t, storage.CreateComment(ctx, c1), "комментарий не создан")

	c2 := &model.Comment{PostID: post.ID, ParentCommentID: &c1.ID}
	require.NoError(t, storage.CreateComment(ctx, c2))

	replies, err := storage.GetReplies(ctx, root.ID)
	require.NoError(t, err)
	expected := []int{c1.ID, c2.ID}

	assert.Equal(t, len(expected), len(replies))
	for i, r := range replies {
		assert.Equal(t, expected[i], r.ID)
	}
}

func TestGetRepliesDeep(t *testing.T) {
	conf()
	post := &model.Post{
		Title:              "Пост",
		AreCommentsAllowed: true,
	}
	require.NoError(t, storage.CreatePost(ctx, post), "пост не создан")

	parentID := 0
	var ids []int
	for i := 1; i <= 5; i++ {
		c := &model.Comment{PostID: post.ID}
		if parentID != 0 {
			c.ParentCommentID = &parentID
		}
		require.NoError(t, storage.CreateComment(ctx, c), "комментарий не создан")
		parentID = c.ID
		ids = append(ids, c.ID)
	}

	replies, err := storage.GetReplies(ctx, ids[0])
	require.NoError(t, err)
	assert.Len(t, replies, len(ids)-1)
	for i, r := range replies {
		assert.Equal(t, ids[i+1], r.ID)
	}
}

func TestGetRepliesDeepAndBranching(t *testing.T) {
	conf()
	post := &model.Post{Title: "Комменты с ветвлением", AreCommentsAllowed: true}
	require.NoError(t, storage.CreatePost(ctx, post), "пост не создан")

	root := &model.Comment{PostID: post.ID}
	require.NoError(t, storage.CreateComment(ctx, root), "комментарий не создан")

	// ветка ответов на корень
	c1 := &model.Comment{PostID: post.ID, ParentCommentID: &root.ID} // 1.2
	c2 := &model.Comment{PostID: post.ID, ParentCommentID: &c1.ID}   // 1.2.3
	c3 := &model.Comment{PostID: post.ID, ParentCommentID: &c2.ID}   // 1.2.3.4
	c4 := &model.Comment{PostID: post.ID, ParentCommentID: &c3.ID}   // 1.2.3.4.5
	c5 := &model.Comment{PostID: post.ID, ParentCommentID: &c4.ID}   // 1.2.3.4.5.6

	for _, c := range []*model.Comment{c1, c2, c3, c4, c5} {
		require.NoError(t, storage.CreateComment(ctx, c), "комментарий не создан")
	}
	// ветвления
	branch1 := &model.Comment{PostID: post.ID, ParentCommentID: &c1.ID}   // 1.2.7
	branch2 := &model.Comment{PostID: post.ID, ParentCommentID: &c3.ID}   // 1.2.3.4.8
	branch3 := &model.Comment{PostID: post.ID, ParentCommentID: &root.ID} // 1.9

	for _, b := range []*model.Comment{branch1, branch2, branch3} {
		require.NoError(t, storage.CreateComment(ctx, b), "комментарий не создан")
	}

	replies, err := storage.GetReplies(ctx, root.ID)
	require.NoError(t, err)

	expectedOrder := []int{
		c1.ID, c2.ID, c3.ID, c4.ID, c5.ID,
		branch2.ID, branch1.ID, branch3.ID,
	}

	// логика такая: сначала выводим всю ветку ответов на 1, т.е.
	// 1, 1.2, 1.3, ... , 1.6
	// после как бы на одном уровне начинаем выводить ответы на ответы, т. е.
	// визуально 1.2.3.4.5 должен быть на одном уровне с 1.2.3.4.8
	// под 1.2 должен быть 1.9 (но не физически ПРЯМО под ним, т.к. сначала идет вся ветка ответов,
	// а как бы на том же визуальном уровне вложенности)
	// понятнее будет, если нарисую на листочке, но вот схематично:
	// 1
	//  | 2
	//  |  | 3
	//  |  |  | 4
	//  |  |    | 5
	//  |  |    |  | 6
	//  |  |    | 8
	//  |  | 7
	//  | 9

	assert.Len(t, replies, len(expectedOrder))
	for i, r := range replies {
		assert.Equal(t, expectedOrder[i], r.ID,
			"неверный порядок комментариев на позиции %d", i)
	}
}

func TestPagination(t *testing.T) {
	conf()
	post := &model.Post{Title: "Пост для пагинации", AreCommentsAllowed: true}
	require.NoError(t, storage.CreatePost(ctx, post), "пост не создан")

	for i := 1; i <= 5; i++ {
		c := &model.Comment{PostID: post.ID, Content: fmt.Sprintf("Коммент %d", i)}
		require.NoError(t, storage.CreateComment(ctx, c))
	}

	limit, offset := 2, 1
	comments, total, err := storage.GetCommentsByPost(ctx, post.ID, limit, offset)
	require.NoError(t, err)
	assert.Equal(t, 5, total)
	assert.Len(t, comments, limit)
}
