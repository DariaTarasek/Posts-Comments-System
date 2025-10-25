package in_memory

import (
	"OzonTestTask/internal/model"
	"context"
	"fmt"
	_ "github.com/stretchr/testify"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreatePost(t *testing.T) {
	store := NewInMemoryStorage()
	post := &model.Post{
		Title:              "Новый пост",
		Content:            "Какой-то текст",
		Author:             "Даша",
		AreCommentsAllowed: true,
	}
	var expected error
	ctx := context.Background()
	err := store.CreatePost(ctx, post)
	assert.Equal(t, expected, err)
}

func TestGetAllPosts(t *testing.T) {
	store := NewInMemoryStorage()
	post := &model.Post{
		Title:              "Новый пост 1",
		Content:            "Какой-то текст",
		Author:             "Даша",
		AreCommentsAllowed: true,
	}
	ctx := context.Background()
	err := store.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("не удалось создать пост: %v ", err)
	}
	time.Sleep(time.Second)
	post = &model.Post{
		Title:              "Новый пост 2",
		Content:            "Какой-то текст",
		Author:             "Даша",
		AreCommentsAllowed: false,
	}
	err = store.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("не удалось создать пост: %v ", err)
	}
	result, err := store.GetAllPosts(ctx)
	fmt.Println(result)
	expectedIds := []int{2, 1}
	actualIds := make([]int, 0)
	for _, v := range result {
		actualIds = append(actualIds, v.ID)
	}

	assert.Equal(t, expectedIds, actualIds)

}

func TestGetPostByID(t *testing.T) {
	store := NewInMemoryStorage()
	post := &model.Post{
		Title:              "Новый пост",
		Content:            "Какой-то текст",
		Author:             "Даша",
		AreCommentsAllowed: true,
	}
	ctx := context.Background()
	err := store.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("не удалось создать пост: %v", err)
	}
	result, err := store.GetPostByID(ctx, 1)
	fmt.Println(*result)
	expectedTitle := "Новый пост"
	assert.Equal(t, expectedTitle, result.Title)

}

func TestCreateComment(t *testing.T) {
	ctx := context.Background()
	store := NewInMemoryStorage()

	post := &model.Post{
		Title:              "Мой новый пост",
		Content:            "Текст",
		AreCommentsAllowed: true,
	}
	if err := store.CreatePost(ctx, post); err != nil {
		t.Fatalf("не получилось создать пост: %v", err)
	}

	comment1 := &model.Comment{PostID: post.ID}
	if err := store.CreateComment(ctx, comment1); err != nil {
		t.Fatalf("не получилось создать коммент: %v", err)
	}
}

func TestCreateReply(t *testing.T) {
	ctx := context.Background()
	store := NewInMemoryStorage()

	post := &model.Post{
		Title:              "Мой новый пост",
		Content:            "Текст",
		AreCommentsAllowed: true,
	}
	if err := store.CreatePost(ctx, post); err != nil {
		t.Fatalf("не получилось создать пост: %v", err)
	}

	comment1 := &model.Comment{PostID: post.ID}
	if err := store.CreateComment(ctx, comment1); err != nil {
		t.Fatalf("не получилось создать коммент: %v", err)
	}
	parentCommentID := 1
	comment2 := &model.Comment{PostID: post.ID, ParentCommentID: &parentCommentID}
	if err := store.CreateComment(ctx, comment2); err != nil {
		t.Fatalf("не получилось создать коммент-ответ: %v", err)
	}
}

func TestGetCommentsByPost(t *testing.T) {
	ctx := context.Background()
	store := NewInMemoryStorage()

	post := &model.Post{
		Title:              "Мой новый пост",
		Content:            "Текст",
		AreCommentsAllowed: true,
	}
	if err := store.CreatePost(ctx, post); err != nil {
		t.Fatalf("не получилось создать пост: %v", err)
	}

	comment1 := &model.Comment{PostID: post.ID}
	if err := store.CreateComment(ctx, comment1); err != nil {
		t.Fatalf("не получилось создать коммент: %v", err)
	}
	parentCommentID := 1
	comment2 := &model.Comment{PostID: post.ID, ParentCommentID: &parentCommentID}
	if err := store.CreateComment(ctx, comment2); err != nil {
		t.Fatalf("не получилось создать коммент-ответ: %v", err)
	}
	result, _, err := store.GetCommentsByPost(ctx, post.ID, 5, 0)
	if err != nil {
		t.Fatalf("не получилось получить корневые комменты к посту: %v", err)
	}
	expectedLen := 1
	assert.Equal(t, expectedLen, len(result))
}

func TestGetReplies(t *testing.T) {
	ctx := context.Background()
	store := NewInMemoryStorage()

	post := &model.Post{
		Title:              "Новый пост",
		Content:            "Какой-то текст",
		AreCommentsAllowed: true,
	}
	if err := store.CreatePost(ctx, post); err != nil {
		t.Fatalf("не получилось создать пост: %v", err)
	}

	comment1 := &model.Comment{PostID: post.ID}
	if err := store.CreateComment(ctx, comment1); err != nil {
		t.Fatalf("не получилось создать коммент: %v", err)
	}

	comment2 := &model.Comment{PostID: post.ID, ParentCommentID: &comment1.ID} // ответ на 1
	if err := store.CreateComment(ctx, comment2); err != nil {
		t.Fatalf("не получилось создать коммент: %v", err)
	}

	comment3 := &model.Comment{PostID: post.ID, ParentCommentID: &comment1.ID} // ответ на 1
	if err := store.CreateComment(ctx, comment3); err != nil {
		t.Fatalf("не получилось создать коммент: %v", err)
	}

	comment4 := &model.Comment{PostID: post.ID, ParentCommentID: &comment2.ID} // ответ на 2
	if err := store.CreateComment(ctx, comment4); err != nil {
		t.Fatalf("не получилось создать коммент: %v", err)
	}

	comment5 := &model.Comment{PostID: post.ID, ParentCommentID: &comment2.ID} // ответ на 2
	if err := store.CreateComment(ctx, comment5); err != nil {
		t.Fatalf("не получилось создать коммент: %v", err)
	}

	replies, err := store.GetReplies(ctx, comment1.ID)
	if err != nil {
		t.Fatalf("не удалось получить ответы на комментарий: %v", err)
	}

	expected := []int{comment2.ID, comment4.ID, comment5.ID, comment3.ID}

	if len(replies) != len(expected) {
		t.Fatalf("expected %d replies, got %d", len(expected), len(replies))
	}

	fmt.Println(replies)

	for i, reply := range replies {
		if reply.ID != expected[i] {
			t.Errorf("at position %d: expected ID %d, got %d", i, expected[i], reply.ID)
		}
	}
}

func TestGetRepliesDeep(t *testing.T) {
	ctx := context.Background()
	store := NewInMemoryStorage()

	// Создаём пост
	post := &model.Post{
		Title:              "Постик",
		Content:            "Текст",
		AreCommentsAllowed: true,
	}
	if err := store.CreatePost(ctx, post); err != nil {
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
		if err := store.CreateComment(ctx, comment); err != nil {
			t.Fatalf("не удалось создать коммент глубины %d: %v", i, err)
		}
		parentID = comment.ID
		expectedOrder = append(expectedOrder, comment.ID)
	}

	replies, err := store.GetReplies(ctx, expectedOrder[0])
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

func TestGetRepliesDeepAndBranching(t *testing.T) {
	ctx := context.Background()
	store := NewInMemoryStorage()

	post := &model.Post{
		Title:              "Жуткий тест",
		Content:            "Текст",
		AreCommentsAllowed: true,
	}
	if err := store.CreatePost(ctx, post); err != nil {
		t.Fatalf("не удалось создать пост: %v", err)
	}

	root := &model.Comment{PostID: post.ID}
	if err := store.CreateComment(ctx, root); err != nil {
		t.Fatalf("не удалось создать коммент: %v", err)
	}

	c1 := &model.Comment{PostID: post.ID, ParentCommentID: &root.ID}
	if err := store.CreateComment(ctx, c1); err != nil {
		t.Fatalf("не удалось создать коммент 1/2: %v", err)
	}

	c2 := &model.Comment{PostID: post.ID, ParentCommentID: &c1.ID}
	if err := store.CreateComment(ctx, c2); err != nil {
		t.Fatalf("не удалось создать коммент 1/2/3: %v", err)
	}

	c3 := &model.Comment{PostID: post.ID, ParentCommentID: &c2.ID}
	if err := store.CreateComment(ctx, c3); err != nil {
		t.Fatalf("не удалось создать коммент 1/2/3/4: %v", err)
	}

	c4 := &model.Comment{PostID: post.ID, ParentCommentID: &c3.ID}
	if err := store.CreateComment(ctx, c4); err != nil {
		t.Fatalf("не удалось создать коммент 1/2/3/4/5: %v", err)
	}

	c5 := &model.Comment{PostID: post.ID, ParentCommentID: &c4.ID}
	if err := store.CreateComment(ctx, c5); err != nil {
		t.Fatalf("не удалось создать коммент 1/2/3/4/5/6: %v", err)
	}

	branch1 := &model.Comment{PostID: post.ID, ParentCommentID: &c1.ID}
	if err := store.CreateComment(ctx, branch1); err != nil {
		t.Fatalf("не удалось создать коммент 1/2/7: %v", err)
	}

	branch2 := &model.Comment{PostID: post.ID, ParentCommentID: &c3.ID}
	if err := store.CreateComment(ctx, branch2); err != nil {
		t.Fatalf("не удалось создать коммент 1/2/3/4/8: %v", err)
	}

	branch3 := &model.Comment{PostID: post.ID, ParentCommentID: &root.ID}
	if err := store.CreateComment(ctx, branch3); err != nil {
		t.Fatalf("не удалось создать коммент 1/9: %v", err)
	}

	replies, err := store.GetReplies(ctx, root.ID)
	if err != nil {
		t.Fatalf("GetReplies returned error: %v", err)
	}

	expectedOrder := []int{
		c1.ID,      // 1.2
		c2.ID,      // 1.2.3
		c3.ID,      // 1.2.3.4
		c4.ID,      // 1.2.3.4.5
		c5.ID,      // 1.2.3.4.5.6
		branch2.ID, // 1.2.3.4.7
		branch1.ID, // 1.2.8
		branch3.ID, // 1.9
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
