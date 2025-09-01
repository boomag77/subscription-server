package appstore

import (
	"fmt"
	"net/http"
	tools "subscription-server/internal/helpers"
	"subscription-server/internal/storage"
	"time"
)

func HandleClientNotification(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	// Handle Client notifications
	
}

func HandleClientRequest(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	// Handle Client requests
}

func HandleAppStoreNotification(w http.ResponseWriter, r *http.Request, s storage.Storage) {

	if err := processAppStoreNotification(r, s); err != nil {
		http.Error(w, fmt.Sprintf("failed to process notification: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func processAppStoreNotification(r *http.Request, s storage.Storage) error {

	parsedNotification, err := parseNotification(r.Body)
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

	s.SetSubscriptionStatus(status)

	return nil
}
