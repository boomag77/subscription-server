package http

import (
	"fmt"
	"net/http"
	"subscription-server/internal/appstore"
	"subscription-server/internal/client"
	"subscription-server/internal/storage"
	"subscription-server/internal/android"
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
		// Handle Apple Store Connect notifications (Server-to-Server)
		appstore.HandleAppleStoreNotification(w, r, storage)
	})

	mux.HandleFunc("/api/v1/notifications/google", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// Handle Google Play notifications (Server-to-Server)
		android.HandleGooglePlayNotification(w, r, storage)
	})

	mux.HandleFunc("/api/v1/notifications/client/ios", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// Handle iOS Client notifications
		clientType := "ios"
		client.HandleClientNotification(w, r, storage, clientType)
	})

	mux.HandleFunc("/api/v1/notifications/client/android", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// Handle Android Client notifications
		clientType := "android"
		client.HandleClientNotification(w, r, storage, clientType)
	})

	mux.HandleFunc("/api/v1/requests/client/ios/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// Handle Client request
		client.HandleClientRequest(w, r, storage)
	})

	mux.HandleFunc("/api/v1/requests/client/android/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// Handle Client status
		client.HandleClientRequest(w, r, storage)
	})

	return mux
}
