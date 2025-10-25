package postgreSQL

import (
	"OzonTestTask/internal/model"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var storage *Storage
var db *sqlx.DB

func TestMain(m *testing.M) {
	const TEST_POSTGRES_DSN = "postgres://postgres:password@localhost:5432/posts-comments-test-db?sslmode=disable"

	var err error
	db, err = NewDBConnection(TEST_POSTGRES_DSN)
	if err != nil {
		log.Fatalf("не удалось подключиться к БД: %v", err)
	}
	fmt.Println("Подключено тестовое хранилище")
	storage = NewStorage(db)
	db.Exec("TRUNCATE TABLE comments CASCADE")
	db.Exec("TRUNCATE TABLE posts CASCADE")
	code := m.Run()
	defer db.Close()

	os.Exit(code)
}

func TestCreateAndGetPost(t *testing.T) {
	ctx := context.Background()

	post := &model.Post{
		Title:              "Тестовый пост 1",
		Content:            "Содержимое поста",
		Author:             "Даша",
		AreCommentsAllowed: true,
	}

	err := storage.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("ошибка при создании поста: %v", err)
	}

	createdPost, err := storage.GetPostByID(ctx, post.ID)
	if err != nil {
		t.Fatalf("ошибка при получении поста: %v", err)
	}
	assert.Equal(t, post.Title, createdPost.Title)
}

func TestGetAllPosts(t *testing.T) {
	ctx := context.Background()

	post1 := &model.Post{
		Title:   "Пост 2",
		Content: "Текст",
		Author:  "Дарья",
	}

	err := storage.CreatePost(ctx, post1)
	if err != nil {
		t.Errorf(" не удалось создать пост: %v", err)
	}
	posts, err := storage.GetAllPosts(ctx)
	if err != nil {
		t.Fatalf("ошибка при получении всех постов: %v", err)
	}
	expectedMinLen := 1
	if expectedMinLen > len(posts) {
		t.Errorf("ожидался минимум 1 пост, получено 0")
	}
}

func TestCreateCommentsAndGetComments(t *testing.T) {
	ctx := context.Background()

	post := &model.Post{
		Title:              "Пост для комментов",
		Content:            "Текст поста",
		Author:             "Василий",
		AreCommentsAllowed: true,
	}
	err := storage.CreatePost(ctx, post)
	if err != nil {
		t.Errorf("не удалось создать пост: %v", err)
	}

	comment1 := &model.Comment{
		PostID:  post.ID,
		Author:  "Анна",
		Content: "Корневой комментарий",
	}
	err = storage.CreateComment(ctx, comment1)
	if err != nil {
		t.Fatalf("ошибка при создании комментария: %v", err)
	}

	comment2 := &model.Comment{
		PostID:          post.ID,
		Author:          "Олег",
		Content:         "Ответ на первый коммент",
		ParentCommentID: &comment1.ID,
	}
	err = storage.CreateComment(ctx, comment2)
	if err != nil {
		t.Fatalf("ошибка при создании вложенного комментария: %v", err)
	}

	rootComments, total, err := storage.GetCommentsByPost(ctx, post.ID, 10, 0)
	if err != nil {
		t.Fatalf("ошибка при получении корневых комментариев: %v", err)
	}
	if total != 1 {
		t.Errorf("ожидался 1 корневой комментарий, получено %d", total)
	}
	if rootComments[0].ID != comment1.ID {
		t.Errorf("ожидалось ID=%d, получено ID=%d", comment1.ID, rootComments[0].ID)
	}

	replies, err := storage.GetReplies(ctx, comment1.ID)
	if err != nil {
		t.Fatalf("ошибка при получении вложенных комментариев: %v", err)
	}
	if len(replies) != 1 {
		t.Errorf("ожидался 1 вложенный комментарий, получено %d", len(replies))
	}
	if replies[0].ID != comment2.ID {
		t.Errorf("ожидалось ID=%d,получено ID=%d", comment2.ID, replies[0].ID)
	}
}

func TestGetRepliesDeep(t *testing.T) {
	ctx := context.Background()

	post := &model.Post{
		Title:              "Постик",
		Content:            "Текст",
		AreCommentsAllowed: true,
	}
	if err := storage.CreatePost(ctx, post); err != nil {
		t.Fatalf("не удалось создать пост: %v", err)
	}

	parentID := 0
	var expectedOrder []int
	for i := 1; i <= 6; i++ {
		comment := &model.Comment{
			PostID: post.ID,
		}
		if parentID != 0 {
			comment.ParentCommentID = &parentID
		}
		if err := storage.CreateComment(ctx, comment); err != nil {
			t.Fatalf("не удалось создать коммент глубины %d: %v", i, err)
		}
		parentID = comment.ID
		expectedOrder = append(expectedOrder, comment.ID)
	}

	replies, err := storage.GetReplies(ctx, expectedOrder[0])
	if err != nil {
		t.Fatalf("не удалось получить вложенные комменты: %v", err)
	}

	if len(replies) != len(expectedOrder)-1 {
		t.Fatalf("ожидались %d комм., получено %d", len(expectedOrder)-1, len(replies))
	}

	for i, reply := range replies {
		if reply.ID != expectedOrder[i+1] {
			t.Errorf("на позиции %d: ожидалось ID %d, получено %d", i, expectedOrder[i+1], reply.ID)
		}
	}
}

func TestGetRepliesDeepAndBranching_Postgres(t *testing.T) {
	ctx := context.Background()

	post := &model.Post{
		Title:              "Жуткий тест",
		Content:            "Текст",
		AreCommentsAllowed: true,
	}
	if err := storage.CreatePost(ctx, post); err != nil {
		t.Fatalf("не удалось создать пост: %v", err)
	}

	root := &model.Comment{PostID: post.ID}
	if err := storage.CreateComment(ctx, root); err != nil {
		t.Fatalf("не удалось создать корневой коммент: %v", err)
	}

	c1 := &model.Comment{PostID: post.ID, ParentCommentID: &root.ID}
	storage.CreateComment(ctx, c1)

	c2 := &model.Comment{PostID: post.ID, ParentCommentID: &c1.ID}
	storage.CreateComment(ctx, c2)

	c3 := &model.Comment{PostID: post.ID, ParentCommentID: &c2.ID}
	storage.CreateComment(ctx, c3)

	c4 := &model.Comment{PostID: post.ID, ParentCommentID: &c3.ID}
	storage.CreateComment(ctx, c4)

	c5 := &model.Comment{PostID: post.ID, ParentCommentID: &c4.ID}
	storage.CreateComment(ctx, c5)

	branch1 := &model.Comment{PostID: post.ID, ParentCommentID: &c1.ID}
	storage.CreateComment(ctx, branch1)

	branch2 := &model.Comment{PostID: post.ID, ParentCommentID: &c3.ID}
	storage.CreateComment(ctx, branch2)

	branch3 := &model.Comment{PostID: post.ID, ParentCommentID: &root.ID}
	storage.CreateComment(ctx, branch3)

	replies, err := storage.GetReplies(ctx, root.ID)
	if err != nil {
		t.Fatalf("не удалось получить вложенные комментарии: %v", err)
	}

	expectedOrder := []int{
		c1.ID,
		c2.ID,
		c3.ID,
		c4.ID,
		c5.ID,
		branch2.ID,
		branch1.ID,
		branch3.ID,
	}

	if len(replies) != len(expectedOrder) {
		t.Fatalf("ожидалось %d комм., получено %d", len(expectedOrder), len(replies))
	}

	for i, reply := range replies {
		if reply.ID != expectedOrder[i] {
			t.Errorf("на позиции %d: ожидалось ID %d, получено %d", i, expectedOrder[i], reply.ID)
		}
	}
}

func TestGetCommentsByPostPagination(t *testing.T) {
	ctx := context.Background()

	post := &model.Post{
		Title:              "Тест пагинации",
		Content:            "Какой-то текст",
		AreCommentsAllowed: true,
	}
	if err := storage.CreatePost(ctx, post); err != nil {
		t.Fatalf("не удалось создать пост: %v", err)
	}

	var rootComments []*model.Comment
	for i := 1; i <= 5; i++ {
		comment := &model.Comment{
			PostID:  post.ID,
			Content: "Комментарий номер " + fmt.Sprint(i),
		}
		if err := storage.CreateComment(ctx, comment); err != nil {
			t.Fatalf("не удалось создать корневой комментарий %d: %v", i, err)
		}
		rootComments = append(rootComments, comment)
	}

	offset := 1
	limit := 2

	comments, total, err := storage.GetCommentsByPost(ctx, post.ID, limit, offset)
	if err != nil {
		t.Fatalf("ошибка при получении корневых комментариев: %v", err)
	}

	if total != 5 {
		t.Errorf("ожидалось всего 5 корневых комментов, получено %d", total)
	}

	if len(comments) != 2 {
		t.Errorf("ожидалось получить %d комментария, получено %d", limit, len(comments))
	}

	for i, c := range comments {
		if c.ID != rootComments[i+offset].ID {
			t.Errorf("на позиции %d: ожидалось ID %d, получено %d", i, rootComments[i].ID, c.ID)
		}
	}
}
