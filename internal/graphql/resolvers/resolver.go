package resolvers

import (
	"OzonTestTask/internal/service"
	"OzonTestTask/internal/subscription"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PostService         service.PostService
	CommentService      service.CommentService
	SubscriptionService subscription.Subscription
}
