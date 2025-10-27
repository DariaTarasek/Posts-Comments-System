package subscription

import (
	"OzonTestTask/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSubscription struct {
	pool   *pgxpool.Pool // pgx для работы с механизмом Listen/Notify в PostgreSQL
	cancel context.CancelFunc
}

func NewPostgresSubscription(pool *pgxpool.Pool) *PostgresSubscription {
	return &PostgresSubscription{pool: pool}
}

// Subscribe Отправка Listen в БД
func (sub *PostgresSubscription) Subscribe(postID int) SubscriptionChan {
	ch := make(SubscriptionChan)
	ctx, cancel := context.WithCancel(context.Background())
	sub.cancel = cancel

	// горутина, подписывающая на канал и ожидающая нового комментария из базы
	go func() {
		conn, err := pgx.Connect(ctx, sub.pool.Config().ConnConfig.ConnString())
		if err != nil {
			close(ch)
			return
		}
		defer conn.Close(ctx)

		channel := fmt.Sprintf("post_%d", postID)
		_, err = conn.Exec(ctx, fmt.Sprintf("LISTEN %s;", channel))
		if err != nil {
			close(ch)
			return
		}

		// ожидание Notify из базы
		for {
			notification, err := conn.WaitForNotification(ctx)
			if err != nil {
				close(ch)
				return
			}
			var comment model.Comment
			err = json.Unmarshal([]byte(notification.Payload), &comment)
			if err != nil {
				return
			}
			ch <- &comment
		}
	}()
	return ch
}

// Publish Отправка Notify в БД
func (sub *PostgresSubscription) Publish(postID int, comment *model.Comment) error {
	commentJSON, err := json.Marshal(comment)
	if err != nil {
		return fmt.Errorf("не удалось сериализовать комментарий в JSON: %v", err)
	}

	_, err = sub.pool.Exec(context.Background(),
		fmt.Sprintf("NOTIFY post_%d, '%s';", postID, string(commentJSON)))
	return err
}

func (sub *PostgresSubscription) Close() error {
	if sub.cancel != nil {
		sub.cancel()
	}
	return nil
}
