package http

import (
	"fmt"
	"net/http"
	"subscription-server/internal/appstore"
	"subscription-server/internal/googleplay"
	"subscription-server/internal/storage"
)

func NewRouter(storage storage.Storage) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprintln(w, "pong!!! new SSH_KEY")
	})

	mux.HandleFunc("/api/v1/notifications/apple/v2", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// Handle App Store Connect notifications (Server-to-Server)
		appstore.HandleAppStoreNotification(w, r, storage)
	})

	mux.HandleFunc("/api/v1/notifications/client/ios", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// Handle Client notifications
		appstore.HandleClientNotification(w, r, storage)
	})

	mux.HandleFunc("/api/v1/notifications/client/android", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// Handle Client notifications
		googleplay.HandleClientNotification(w, r, storage)
	})

	mux.HandleFunc("/api/v1/requests/client/ios/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// Handle Client request
		appstore.HandleClientRequest(w, r, storage)
	})

	mux.HandleFunc("/api/v1/requests/client/android/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// Handle Client status
		googleplay.HandleClientRequest(w, r, storage)
	})

	return mux
}
