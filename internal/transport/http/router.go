package http

import (
	"fmt"
	"net/http"
	"subscription-server/internal/deps"
)

func NewRouter(d *deps.Deps) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprintln(w, "pong!!! almast full router")
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
