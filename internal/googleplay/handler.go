package googleplay

import (
	"net/http"
	"subscription-server/internal/storage"
)

func HandleClientNotification(w http.ResponseWriter, r *http.Request, store storage.Storage) {
	// Handle Client notifications
}

func HandleClientRequest(w http.ResponseWriter, r *http.Request, store storage.Storage) {
	// Handle Client requests
}
