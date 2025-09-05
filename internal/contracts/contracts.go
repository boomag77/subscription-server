package contracts

import "net/http"

type JWSValidator interface {
	Validate(header string, payload string, signature string) error
}

type Service interface {
	HandleProviderNotification(w http.ResponseWriter, r *http.Request)
	HandleClientNotification(w http.ResponseWriter, r *http.Request)
	HandleClientRequest(w http.ResponseWriter, r *http.Request)
}
