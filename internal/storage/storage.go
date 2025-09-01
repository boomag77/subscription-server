package storage

import "time"

type SubscriptionStatus struct {
	ExpiresAt             time.Time
	UserToken             string
	ProductID             string
	OriginalTransactionID string
	IsActive              bool
}

type Storage interface {
	GetSubscriptionStatus(userToken string) (*SubscriptionStatus, error)
	SetSubscriptionStatus(status *SubscriptionStatus) error
}

