package http

import (
	"fmt"
	"net/http"
	"subscription-server/internal/appstore"
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

	mux.HandleFunc("/appstoreconnectnotification/v2", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// Handle App Store Connect notifications
		appstore.HandleAppStoreNotification(w, r, storage)
	})

	return mux
}
