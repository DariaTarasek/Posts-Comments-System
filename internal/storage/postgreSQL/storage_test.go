package postgreSQL

import (
	"OzonTestTask/internal/model"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

var (
	storage *Storage
	db      *sqlx.DB
	ctx     context.Context
)

func TestMain(m *testing.M) {
	const TEST_POSTGRES_DSN = "postgres://postgres:password@localhost:5432/posts-comments-test-db?sslmode=disable"

	var err error
	db, err = NewDBConnection(TEST_POSTGRES_DSN)
	if err != nil {
		log.Fatalf("не удалось подключиться к тестовой БД: %v", err)
	}
	storage = NewStorage(db)
	ctx = context.Background()

	_, _ = db.Exec("TRUNCATE TABLE comments CASCADE")
	_, _ = db.Exec("TRUNCATE TABLE posts CASCADE")

	code := m.Run()
	defer db.Close()
	os.Exit(code)
}

func TestCreateAndGetPost(t *testing.T) {
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
	post := &model.Post{
		Title:   "Пост",
		Content: "Текст",
		Author:  "Дарья",
	}
	require.NoError(t, storage.CreatePost(ctx, post), "пост не создан")
	posts, err := storage.GetAllPosts(ctx)
	require.NoError(t, err, "не удалось получить посты")
	assert.NotEmpty(t, posts, "должен быть хотя бы один пост")
}

func TestCreateAndGetComments(t *testing.T) {
	post := &model.Post{
		Title:              "Пост с комментами",
		Content:            "Текст",
		Author:             "Василий",
		AreCommentsAllowed: true,
	}
	require.NoError(t, storage.CreatePost(ctx, post), "пост не создан")

	root := &model.Comment{
		PostID:  post.ID,
		Author:  "Анна",
		Content: "Корневой",
	}
	require.NoError(t, storage.CreateComment(ctx, root), "комментарий не создан")

	reply := &model.Comment{
		PostID:          post.ID,
		ParentCommentID: &root.ID,
		Author:          "Олег",
		Content:         "Ответ",
	}
	require.NoError(t, storage.CreateComment(ctx, reply), "комментарий не создан")

	rootComments, total, err := storage.GetCommentsByPost(ctx, post.ID, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Equal(t, root.ID, rootComments[0].ID)

	replies, err := storage.GetReplies(ctx, root.ID)
	require.NoError(t, err)
	assert.Len(t, replies, 1)
	assert.Equal(t, reply.ID, replies[0].ID)
}

func TestGetPostByID_WrongID(t *testing.T) {
	post, err := storage.GetPostByID(ctx, -1)
	assert.Error(t, err)
	assert.Nil(t, post)
}

func TestCreateComment_WrongPostID(t *testing.T) {
	comment := &model.Comment{
		PostID:  -1,
		Author:  "Тест",
		Content: "Невалидный пост",
	}
	err := storage.CreateComment(ctx, comment)
	assert.Error(t, err)
}

func TestGetRepliesDeep(t *testing.T) {
	post := &model.Post{
		Title:              "Пост",
		Content:            "Текст",
		AreCommentsAllowed: true,
	}
	require.NoError(t, storage.CreatePost(ctx, post), "пост не создан")

	parentID := 0
	var expectedIDs []int
	for i := 1; i <= 5; i++ {
		c := &model.Comment{PostID: post.ID}
		if parentID != 0 {
			c.ParentCommentID = &parentID
		}
		require.NoError(t, storage.CreateComment(ctx, c), "комментарий не создан")
		parentID = c.ID
		expectedIDs = append(expectedIDs, c.ID)
	}

	replies, err := storage.GetReplies(ctx, expectedIDs[0])
	require.NoError(t, err)
	assert.Len(t, replies, len(expectedIDs)-1)

	for i, reply := range replies {
		assert.Equal(t, expectedIDs[i+1], reply.ID)
	}
}

func TestPagination(t *testing.T) {
	post := &model.Post{
		Title:              "Пост для пагинации",
		Content:            "Контент",
		AreCommentsAllowed: true,
	}
	require.NoError(t, storage.CreatePost(ctx, post), "пост не создан")

	for i := 1; i <= 5; i++ {
		c := &model.Comment{
			PostID:  post.ID,
			Content: fmt.Sprintf("Коммент %d", i),
		}
		require.NoError(t, storage.CreateComment(ctx, c), "комментарий не создан")
	}

	limit, offset := 2, 1
	comments, total, err := storage.GetCommentsByPost(ctx, post.ID, limit, offset)
	require.NoError(t, err)
	assert.Equal(t, 5, total)
	assert.Len(t, comments, limit)
}

func TestGetRepliesDeepAndBranching(t *testing.T) {
	post := &model.Post{
		Title:              "Комменты с ветвлениями",
		Content:            "Текст",
		AreCommentsAllowed: true,
	}
	require.NoError(t, storage.CreatePost(ctx, post), "пост не создан")

	// корневой коммент
	root := &model.Comment{
		PostID: post.ID,
	}
	require.NoError(t, storage.CreateComment(ctx, root))

	// ветка ответов на корень
	c1 := &model.Comment{ // 1.2
		PostID:          post.ID,
		ParentCommentID: &root.ID,
	}
	c2 := &model.Comment{ // 1.2.3
		PostID:          post.ID,
		ParentCommentID: &c1.ID,
	}
	c3 := &model.Comment{ // 1.2.3.4
		PostID:          post.ID,
		ParentCommentID: &c2.ID,
	}

	c4 := &model.Comment{ // 1.2.3.4.5
		PostID:          post.ID,
		ParentCommentID: &c3.ID,
	}

	c5 := &model.Comment{ // 1.2.3.4.5.6
		PostID:          post.ID,
		ParentCommentID: &c4.ID,
	}

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
	require.NoError(t, err, "ошибка при получении вложенных комментариев")

	expectedCount := 8
	assert.Len(t, replies, expectedCount, "неверное количество комментариев")
	expectedOrder := []int{
		c1.ID, c2.ID, c3.ID, c4.ID, c5.ID, branch2.ID, branch1.ID, branch3.ID,
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

	for i, reply := range replies {
		assert.Equal(t, expectedOrder[i], reply.ID,
			"неверный порядок комментариев на позиции %d", i)
	}
}
