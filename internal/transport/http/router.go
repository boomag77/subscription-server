package http

import (
	"encoding/json"
	"net/http"
	"subscription-server/internal/deps"
)

func NewRouter(d *deps.Deps) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		response := map[string]string{"status": "ok"}
		json.NewEncoder(w).Encode(response)
	})

	mux.HandleFunc("/api/v1/notifications/apple/v2", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// Handle Apple Store Connect notifications (Server-to-Server)
		d.AppleService.HandleProviderNotification(w, r)
	})

	mux.HandleFunc("/api/v1/notifications/google", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// Handle Google Play notifications (Server-to-Server)
		d.GoogleService.HandleProviderNotification(w, r)
	})

	mux.HandleFunc("/api/v1/notifications/client/ios", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// Handle iOS Client notifications
		d.AppleService.HandleClientNotification(w, r)
	})

	mux.HandleFunc("/api/v1/notifications/client/android", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// Handle Android Client notifications
		d.GoogleService.HandleClientNotification(w, r)
	})

	mux.HandleFunc("/api/v1/requests/client/ios/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// Handle Client request
		d.AppleService.HandleClientRequest(w, r)
	})

	mux.HandleFunc("/api/v1/requests/client/android/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// Handle Client status
		d.GoogleService.HandleClientRequest(w, r)
	})

	return mux
}
