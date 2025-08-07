package storage

import (
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

func (m *memoryStorage) GetSubscriptionStatus(userToken string) (*SubscriptionStatus, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status, exists := m.data[userToken]
	if !exists {
		return nil, ErrSubscriptionNotFound
	}
	copy := *status

	return &copy, nil
}

func (m *memoryStorage) SetSubscriptionStatus(status *SubscriptionStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[status.UserToken] = status

	copy := *status
	m.data[status.UserToken] = &copy

	return nil
}
