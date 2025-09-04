package appstore

import (
	"fmt"
	"net/http"
	tools "subscription-server/internal/helpers"
	"subscription-server/internal/jws"
	"subscription-server/internal/logger"
	"subscription-server/internal/service"
	"subscription-server/internal/storage"
	"time"
)

type appleStoreService struct {
	storage   storage.Storage
	validator jws.JWSValidator
	logger    logger.Logger
}

func NewAppleStoreService(st storage.Storage, l logger.Logger) service.Service {
	v := NewAppleJWSValidator()
	return &appleStoreService{
		storage:   st,
		validator: v,
		logger:    l,
	}
}

func (s *appleStoreService) HandleProviderNotification(w http.ResponseWriter, r *http.Request) {

	if err := s.processProviderNotification(r); err != nil {
		http.Error(w, fmt.Sprintf("failed to process notification: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *appleStoreService) HandleClientNotification(w http.ResponseWriter, r *http.Request) {

	if err := s.processIOSClientNotification(r); err != nil {
		http.Error(w, fmt.Sprintf("failed to process iOS client notification: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *appleStoreService) HandleClientRequest(w http.ResponseWriter, r *http.Request) {
	// Handle Client requests
}

func (s *appleStoreService) processProviderNotification(r *http.Request) error {

	parsedNotification, err := parseAppStoreNotification(r.Body)
	if err != nil {
		return fmt.Errorf("failed to parse notification: %w", err)
	}

	parsedTx, err := parseTransaction(parsedNotification.Data.SignedTransactionInfo)
	if err != nil {
		return fmt.Errorf("failed to parse transaction: %w", err)
	}

	parsedRenewalInfo, err := parseRenewalInfo(parsedNotification.Data.SignedRenewalInfo)
	if err != nil {
		return fmt.Errorf("failed to parse renewal info: %w", err)
	}

	user := parsedNotification.Data.AppAccountToken
	if user == "" {
		user = "tx:" + parsedTx.OriginalTransactionID
	}

	expiresAt := tools.MsToTime(parsedTx.ExpiresDateMS)

	var grace time.Time

	if parsedRenewalInfo != nil {
		if t := tools.MsToTime(parsedRenewalInfo.GracePeriodExpiresDateMS); !t.IsZero() {
			grace = t
		}
	}

	activeUntil := expiresAt
	if grace.After(expiresAt) {
		activeUntil = grace
	}

	now := time.Now().UTC()

	isActive := !activeUntil.IsZero() && now.Before(activeUntil)

	if parsedTx.RevocationDateMS != nil && *parsedTx.RevocationDateMS > 0 {
		isActive = false
	}
	if parsedNotification.NotificationType == "EXPIRED" {
		isActive = false
	}

	status := &storage.SubscriptionStatus{
		ExpiresAt:             activeUntil,
		UserToken:             user,
		ProductID:             parsedTx.ProductID,
		OriginalTransactionID: parsedTx.OriginalTransactionID,
		IsActive:              isActive,
	}

	s.storage.SetSubscriptionStatus(status)

	return nil
}
