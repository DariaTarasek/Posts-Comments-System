package in_memory

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"OzonTestTask/internal/model"
)

type InMemoryStorage struct {
	mu               sync.RWMutex
	posts            map[int]model.Post
	postsByCreatedAt []int
	comments         map[int]model.Comment
	commentsByPost   map[int][]int
	replies          map[int][]int

	nextPostID    int
	nextCommentID int
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		posts:          make(map[int]model.Post),
		comments:       make(map[int]model.Comment),
		commentsByPost: make(map[int][]int),
		replies:        make(map[int][]int),
		nextPostID:     1,
		nextCommentID:  1,
	}
}

// контекст в in-memory не использую,
//т.к. все операции выполняются быстро,
//нет риска что, например, отвалится соединение, как может быть в БД

// CreatePost Создание нового поста
func (ms *InMemoryStorage) CreatePost(ctx context.Context, post *model.Post) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	post.ID = ms.nextPostID
	ms.nextPostID++
	post.CreatedAt = time.Now().UTC()
	ms.posts[post.ID] = *post
	ms.postsByCreatedAt = append(ms.postsByCreatedAt, post.ID)

	return nil
}

// GetAllPosts Получение всех постов
func (ms *InMemoryStorage) GetAllPosts(ctx context.Context) ([]model.Post, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	posts := make([]model.Post, 0, len(ms.postsByCreatedAt))
	// сначала выдаю все новые посты - мне кажется, это логично для новостной ленты
	for i := len(ms.postsByCreatedAt) - 1; i >= 0; i-- {
		id := ms.postsByCreatedAt[i]
		posts = append(posts, ms.posts[id])
	}
	return posts, nil
}

// GetPostByID Получение поста по ID
func (ms *InMemoryStorage) GetPostByID(ctx context.Context, id int) (*model.Post, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	p, ok := ms.posts[id]
	if !ok {
		return nil, fmt.Errorf("пост не найден")
	}
	return &p, nil
}

// CreateComment Создание комментария
func (ms *InMemoryStorage) CreateComment(ctx context.Context, comment *model.Comment) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	_, ok := ms.posts[comment.PostID]
	if !ok {
		return fmt.Errorf("пост для добавления комментария не найден")
	}

	comment.ID = ms.nextCommentID
	ms.nextCommentID++
	comment.CreatedAt = time.Now().UTC()

	if comment.ParentCommentID != nil {
		parent, ok := ms.comments[*comment.ParentCommentID]
		if !ok {
			return fmt.Errorf("комментарий для ответа не найден")
		}
		comment.Path = parent.Path + "." + strconv.Itoa(comment.ID)
	} else {
		comment.Path = strconv.Itoa(comment.ID)
	}

	ms.comments[comment.ID] = *comment

	// если коммент - ответ на другой коммент - кладу его в мапу ответов
	if comment.ParentCommentID != nil {
		ms.replies[*comment.ParentCommentID] = append(ms.replies[*comment.ParentCommentID], comment.ID)
	} else { // если коммент корневой - кладу в мапу корневых комментов
		ms.commentsByPost[comment.PostID] = append(ms.commentsByPost[comment.PostID], comment.ID)
	}

	return nil
}

// GetCommentsByPost Получение корневых комментариев к посту
func (ms *InMemoryStorage) GetCommentsByPost(ctx context.Context, postID, limit, offset int) ([]model.Comment, int, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	rootIDs := ms.commentsByPost[postID]
	amount := len(rootIDs)

	if offset >= len(rootIDs) {
		return []model.Comment{}, amount, nil
	}

	lastComment := offset + limit
	if lastComment > len(rootIDs) {
		lastComment = len(rootIDs)
	}

	result := make([]model.Comment, 0, lastComment-offset)
	for _, id := range rootIDs[offset:lastComment] {
		result = append(result, ms.comments[id])
	}

	return result, amount, nil
}

// GetReplies Получение ветки ответов на комментарий (все уровни вложенности)
func (ms *InMemoryStorage) GetReplies(ctx context.Context, parentID int) ([]model.Comment, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	if _, ok := ms.comments[parentID]; !ok {
		return nil, fmt.Errorf("комментарий не найден")
	}

	// результат отдаю плоским, пускай строит клиент,
	// чтобы не перегружать сервер при большом количестве комментов.
	var result []model.Comment
	stack := []int{}

	replies := ms.replies[parentID]
	for i := len(replies) - 1; i >= 0; i-- {
		stack = append(stack, replies[i])
	}

	for len(stack) > 0 {
		currentID := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		result = append(result, ms.comments[currentID])

		replies := ms.replies[currentID]
		for i := len(replies) - 1; i >= 0; i-- {
			stack = append(stack, replies[i])
		}
	}

	return result, nil
}
