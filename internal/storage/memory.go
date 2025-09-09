package storage

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
)

type memoryStorage struct {
	mu   sync.RWMutex
	data map[string]*SubscriptionStatus
}

func NewMemoryStorage() Storage {
	return &memoryStorage{
		data: make(map[string]*SubscriptionStatus),
	}
}

func (m *memoryStorage) GetSubscriptionStatus(ctx context.Context, userToken string) (*SubscriptionStatus, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		m.mu.RLock()
		defer m.mu.RUnlock()
		status, exists := m.data[userToken]
		if !exists {
			return nil, ErrSubscriptionNotFound
		}
		copy := *status

		return &copy, nil
	}

}

func (m *memoryStorage) SetSubscriptionStatus(ctx context.Context, status *SubscriptionStatus) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		m.mu.Lock()
		defer m.mu.Unlock()

		copy := *status
		m.data[status.UserToken] = &copy
		return nil
	}
}
