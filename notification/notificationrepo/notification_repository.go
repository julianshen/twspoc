// store/notification_repository.go
package notificationrepo

import (
	"context"
	"notification/types"
)

// NotificationRepository defines CRUD and subscription for notifications.
type NotificationRepository interface {
	Create(ctx context.Context, n *types.Notification) error
	Get(ctx context.Context, id string) (*types.Notification, error)
	Update(ctx context.Context, n *types.Notification) error
	Delete(ctx context.Context, id string) error
	ListByUser(ctx context.Context, userID string) ([]*types.Notification, error)
	Subscribe(ctx context.Context, userID string) (<-chan *types.Notification, error)
}
