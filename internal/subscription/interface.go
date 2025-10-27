package subscription

import "OzonTestTask/internal/model"

type SubscriptionChan chan *model.Comment

type Subscription interface {
	Subscribe(postID int) SubscriptionChan
	Publish(postID int, comment *model.Comment) error
	Close() error
}
