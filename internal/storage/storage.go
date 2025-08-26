package storage

import "time"

type SubscriptionStatus struct {
	UserToken             string
	ProductID             string
	OriginalTransactionID string
	ExpiresAt             time.Time
	IsActive              bool
}

type Storage interface {
	GetSubscriptionStatus(userToken string) (*SubscriptionStatus, error)
	SetSubscriptionStatus(status *SubscriptionStatus) error
}
