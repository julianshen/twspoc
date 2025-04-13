// notificationrepo/mock_notification_repository.go
package notificationrepo

import (
	"context"
	"notification/types"

	"github.com/stretchr/testify/mock"
)

// MockNotificationRepository is a mock implementation of NotificationRepository for unit testing.
type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) Create(ctx context.Context, n *types.Notification) error {
	args := m.Called(ctx, n)
	return args.Error(0)
}

func (m *MockNotificationRepository) Get(ctx context.Context, id string) (*types.Notification, error) {
	args := m.Called(ctx, id)
	var n *types.Notification
	if args.Get(0) != nil {
		n = args.Get(0).(*types.Notification)
	}
	return n, args.Error(1)
}

func (m *MockNotificationRepository) Update(ctx context.Context, n *types.Notification) error {
	args := m.Called(ctx, n)
	return args.Error(0)
}

func (m *MockNotificationRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockNotificationRepository) ListByUser(ctx context.Context, userID string) ([]*types.Notification, error) {
	args := m.Called(ctx, userID)
	var notifications []*types.Notification
	if args.Get(0) != nil {
		notifications = args.Get(0).([]*types.Notification)
	}
	return notifications, args.Error(1)
}

func (m *MockNotificationRepository) Subscribe(ctx context.Context, userID string) (<-chan *types.Notification, error) {
	args := m.Called(ctx, userID)
	var ch <-chan *types.Notification
	if args.Get(0) != nil {
		ch = args.Get(0).(<-chan *types.Notification)
	}
	return ch, args.Error(1)
}
