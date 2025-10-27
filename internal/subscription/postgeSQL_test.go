package subscription

import (
	"OzonTestTask/internal/model"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDBSubscribeAndPublish(t *testing.T) {
	connStr := "postgres://postgres:password@localhost:5432/posts-comments-test-db?sslmode=disable"
	pool, err := NewPGXPool(connStr)
	require.NoError(t, err)
	defer pool.Close()
	sub := NewPostgresSubscription(pool)
	postID := 1
	sub1 := sub.Subscribe(postID)
	sub2 := sub.Subscribe(postID)
	time.Sleep(1 * time.Second)

	comment := &model.Comment{PostID: postID, Author: "Я", Content: "Тестик"}
	err = sub.Publish(postID, comment)
	require.NoError(t, err)
	res1 := <-sub1
	require.Equal(t, comment, res1)
	res2 := <-sub2
	require.Equal(t, comment, res2)
}
