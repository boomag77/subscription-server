package android

import (
	"fmt"
	"net/http"
	"subscription-server/internal/storage"
)

func HandleGooglePlayNotification(w http.ResponseWriter, r *http.Request, store storage.Storage) {

	if err := processGooglePlayNotification(r, store); err != nil {
		http.Error(w, fmt.Sprintf("failed to process notification: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func processGooglePlayNotification(r *http.Request, s storage.Storage) error {
	return nil
}

func ProcessAndroidClientNotification(r *http.Request, s storage.Storage) error {
	return nil
}
