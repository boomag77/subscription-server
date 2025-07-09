package http

import (
	"fmt"
	"net/http"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprintln(w, "pong!!!")
	})

	return mux
}
