package subscription

import (
	"OzonTestTask/internal/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSubscribeAndPublish(t *testing.T) {
	sub := NewInMemorySubscription()
	ch := sub.Subscribe(1)
	require.NotNil(t, ch)
	comment := &model.Comment{ID: 1, Author: "Дарья", Content: "Коммент"}
	go func() {
		err := sub.Publish(1, comment)
		require.NoError(t, err)
	}()
	result := <-ch
	require.Equal(t, comment, result)
}

func TestManySubscribers(t *testing.T) {
	sub := NewInMemorySubscription()
	ch1 := sub.Subscribe(1)
	ch2 := sub.Subscribe(1)
	comment := &model.Comment{ID: 1, Author: "Саша", Content: "Новый коммент"}

	go func() {
		err := sub.Publish(1, comment)
		require.NoError(t, err)
	}()

	received1 := <-ch1
	received2 := <-ch2
	require.Equal(t, comment, received1)
	require.Equal(t, comment, received2)
}

func TestSubscribeDifferentPosts(t *testing.T) {
	sub := NewInMemorySubscription()
	ch1 := sub.Subscribe(1)
	ch2 := sub.Subscribe(2)
	comment1 := &model.Comment{ID: 1, Content: "Пост 1"}
	comment2 := &model.Comment{ID: 2, Content: "Пост 2"}

	go func() {
		err := sub.Publish(1, comment1)
		require.NoError(t, err)
	}()
	go func() {
		err := sub.Publish(2, comment2)
		require.NoError(t, err)
	}()
	received1 := <-ch1
	received2 := <-ch2
	require.Equal(t, comment1, received1)
	require.Equal(t, comment2, received2)
}

func TestClose(t *testing.T) {
	sub := NewInMemorySubscription()
	ch1 := sub.Subscribe(1)
	err := sub.Close()
	require.NoError(t, err)
	select {
	case _, ok := <-ch1:
		require.False(t, ok)
	default:
		t.Fatal("ch1 должен быть закрыт")
	}
}
