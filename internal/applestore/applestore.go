package applestore

import (
	"fmt"
	"net/http"
	"subscription-server/internal/contracts"
	tools "subscription-server/internal/helpers"
	"subscription-server/internal/logger"
	"subscription-server/internal/storage"
	"time"
)

type appleStoreService struct {
	storage storage.Storage
	logger  logger.Logger
	parser  *appleParser
}

func NewAppleStoreService(st storage.Storage, l logger.Logger, p *appleParser) contracts.Service {
	return &appleStoreService{
		storage: st,
		parser:  p,
		logger:  l,
	}
}

func (s *appleStoreService) HandleProviderNotification(w http.ResponseWriter, r *http.Request) {

	if err := s.processProviderNotification(r); err != nil {
		http.Error(w, fmt.Sprintf("failed to process notification: %v", err), http.StatusInternalServerError)
		// s.logger.Log(logger.LogMessage{
		// 	Level:   "error",
		// 	Sender:  "appleStoreService",
		// 	Message: fmt.Sprintf("failed to process notification: %v", err),
		// 	Time:    time.Now().UTC(),
		// })
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *appleStoreService) processIOSClientNotification(r *http.Request) error {
	parsedClientNotification, err := s.parser.ParseClientNotification(r.Body)
	if err != nil {
		return fmt.Errorf("failed to parse client notification: %w", err)

	}

	signedTx := parsedClientNotification.SignedTransactionInfo
	parsedClientTx, err := s.parser.ParseTransaction(signedTx)
	if err != nil {
		return fmt.Errorf("failed to parse client transaction: %w", err)
	}
	user := parsedClientNotification.AppAccountToken
	if user == "" {
		user = "tx:" + parsedClientTx.OriginalTransactionID
	}
	expiresAt := tools.MsToTime(parsedClientTx.ExpiresDateMS)
	now := time.Now().UTC()
	isActive := !expiresAt.IsZero() && now.Before(expiresAt)
	if parsedClientTx.RevocationDateMS != nil && *parsedClientTx.RevocationDateMS > 0 {
		isActive = false
	}

	status := &storage.SubscriptionStatus{
		ExpiresAt:             expiresAt,
		UserToken:             user,
		ProductID:             parsedClientTx.ProductID,
		OriginalTransactionID: parsedClientTx.OriginalTransactionID,
		IsActive:              isActive,
	}
	s.storage.SetSubscriptionStatus(status)

	return nil
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

	parsedNotification, err := s.parser.ParseAppStoreNotification(r.Body)
	if err != nil {
		return fmt.Errorf("failed to parse notification: %w", err)
	}

	parsedTx, err := s.parser.ParseTransaction(parsedNotification.Data.SignedTransactionInfo)
	if err != nil {
		return fmt.Errorf("failed to parse transaction: %w", err)
	}

	parsedRenewalInfo, err := s.parser.ParseRenewalInfo(parsedNotification.Data.SignedRenewalInfo)
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
