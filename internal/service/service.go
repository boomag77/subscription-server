package service

import (
	"net/http"
)

type Service interface {
	HandleProviderNotification(w http.ResponseWriter, r *http.Request)
	HandleClientNotification(w http.ResponseWriter, r *http.Request)
	HandleClientRequest(w http.ResponseWriter, r *http.Request)
}
