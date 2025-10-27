package resolvers

import (
	"OzonTestTask/internal/mocks"
	"OzonTestTask/internal/model"
	"OzonTestTask/internal/subscription"
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

var ctx = context.Background()

func TestCreatePost(t *testing.T) {
	mockPostService := new(mocks.PostService)
	r := &Resolver{PostService: mockPostService}
	mutation := &mutationResolver{r}
	mockPostService.
		On("CreatePost", mock.Anything, mock.AnythingOfType("*model.Post")).
		Return(nil)

	title := "Тестовый пост"
	content := "Учусь работать с моками"
	author := "Даша"
	areCommentsAllowed := true

	post, err := mutation.CreatePost(ctx, title, content, author, areCommentsAllowed)
	require.NoError(t, err)
	require.Equal(t, title, post.Title)
	mockPostService.AssertExpectations(t)
}

func TestCreateComment(t *testing.T) {
	mockCommentService := new(mocks.CommentService)
	r := &Resolver{CommentService: mockCommentService}
	mutation := &mutationResolver{r}
	mockCommentService.
		On("CreateComment", mock.Anything, mock.AnythingOfType("*model.Comment")).
		Return(nil)

	postID := "1"
	content := "Тестовый комментарий"
	author := "Дарья"

	comment, err := mutation.CreateComment(ctx, postID, nil, author, content)

	require.NoError(t, err)
	require.Equal(t, content, comment.Content)

	mockCommentService.AssertExpectations(t)
}

func TestGetPosts(t *testing.T) {
	mockPostService := new(mocks.PostService)
	r := &Resolver{PostService: mockPostService}
	query := &queryResolver{r}
	mockPostService.
		On("GetAllPosts", mock.Anything).
		Return([]model.Post{
			{ID: 1},
			{ID: 2},
		}, nil)

	posts, err := query.Posts(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, posts[0].ID)
	require.Equal(t, 2, posts[1].ID)
	mockPostService.AssertExpectations(t)
}

func TestGetPost(t *testing.T) {
	mockCommentService := new(mocks.CommentService)
	r := &Resolver{CommentService: mockCommentService}
	query := &queryResolver{r}
	mockCommentService.On("GetPostByID", mock.Anything, mock.AnythingOfType("int")).
		Return(&model.Post{ID: 1, Author: "Дарья"}, nil)

	id := "5"
	post, err := query.Post(ctx, id)
	require.NoError(t, err)
	require.Equal(t, "Дарья", post.Author)
	require.Equal(t, 1, post.ID)
	mockCommentService.AssertExpectations(t)
}

func TestGetComments(t *testing.T) {
	mockCommentService := new(mocks.CommentService)
	r := &Resolver{CommentService: mockCommentService}
	query := &postResolver{Resolver: r}

	mockCommentService.
		On("GetCommentsByPost", mock.Anything, mock.AnythingOfType("int"),
			mock.AnythingOfType("int"), mock.AnythingOfType("int")).
		Return([]model.Comment{
			{ID: 1, Author: "Иван"},
		}, 1, nil)

	post := model.Post{
		ID:                 1,
		Title:              "пост",
		Content:            "текст",
		Author:             "Дарья Валерьевна",
		AreCommentsAllowed: true,
	}
	limit := 5
	offset := 0

	paginatedComments, err := query.Comments(ctx, &post, &limit, &offset)

	require.NoError(t, err)
	require.Len(t, paginatedComments.Comments, 1)
	require.Equal(t, "Иван", paginatedComments.Comments[0].Author)
	require.Equal(t, 1, paginatedComments.TotalPages)

	mockCommentService.AssertExpectations(t)
}

func TestGetReplies(t *testing.T) {
	mockCommentService := new(mocks.CommentService)
	r := &Resolver{CommentService: mockCommentService}
	query := &queryResolver{Resolver: r}

	mockCommentService.On("GetReplies", mock.Anything, mock.AnythingOfType("int")).
		Return([]model.Comment{
			{ID: 1, Author: "Сергей", Content: "Ответ 1"},
			{ID: 2, Author: "Александр", Content: "Ответ 2"}}, nil)

	replies, err := query.Replies(ctx, "1")
	require.NoError(t, err)
	require.Len(t, replies, 2)
	require.Equal(t, "Ответ 1", replies[0].Content)

	mockCommentService.AssertExpectations(t)
}

func TestSubscription(t *testing.T) {
	mockSubscription := new(mocks.Subscription)
	r := &Resolver{SubscriptionService: mockSubscription}
	sub := &subscriptionResolver{Resolver: r}

	ch := make(subscription.SubscriptionChan, 1)
	mockSubscription.On("Subscribe", mock.AnythingOfType("int")).
		Return(ch)

	result, err := sub.NewComment(ctx, 1)
	require.NoError(t, err)

	comment := &model.Comment{ID: 1, Author: "Мария", Content: "Привет"}
	ch <- comment

	received := <-result
	require.Equal(t, comment, received)
	mockSubscription.AssertExpectations(t)
}
