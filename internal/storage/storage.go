package storage

import (
	"context"
	"time"
)

type SubscriptionStatus struct {
	ExpiresAt             time.Time
	UserToken             string
	ProductID             string
	OriginalTransactionID string
	IsActive              bool
}

type Storage interface {
	GetSubscriptionStatus(ctx context.Context, userToken string) (*SubscriptionStatus, error)
	SetSubscriptionStatus(ctx context.Context, status *SubscriptionStatus) error
}
