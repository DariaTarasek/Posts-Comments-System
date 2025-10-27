package subscription

import (
	"OzonTestTask/internal/model"
	"sync"
)

type InMemorySubscription struct {
	mu          sync.RWMutex
	subscribers map[int][]SubscriptionChan
}

func NewInMemorySubscription() *InMemorySubscription {
	return &InMemorySubscription{subscribers: make(map[int][]SubscriptionChan)}
}

// Subscribe Добавление подписчиков на пост
func (sub *InMemorySubscription) Subscribe(postID int) SubscriptionChan {
	ch := make(SubscriptionChan)
	sub.mu.Lock()
	defer sub.mu.Unlock()
	sub.subscribers[postID] = append(sub.subscribers[postID], ch)
	return ch
}

// Publish Публикация нового комментария в канал
func (sub *InMemorySubscription) Publish(postID int, comment *model.Comment) error {
	sub.mu.RLock()
	postSubs := sub.subscribers[postID]
	sub.mu.RUnlock()

	for _, postSub := range postSubs {
		go func(s SubscriptionChan) {
			s <- comment
		}(postSub)
	}
	return nil
}

// Close Закрытие всех каналов и очистка структуры
func (sub *InMemorySubscription) Close() error {
	sub.mu.Lock()
	defer sub.mu.Unlock()
	for _, v := range sub.subscribers {
		for _, ch := range v {
			close(ch)
		}
	}
	sub.subscribers = make(map[int][]SubscriptionChan)
	return nil
}
