package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	appstore "subscription-server/internal/applestore"
	"subscription-server/internal/deps"
	"subscription-server/internal/logger"
	"subscription-server/internal/storage"
	httpTransport "subscription-server/internal/transport/http"
)

func main() {

	localStorage := storage.NewMemoryStorage()

	port := ":443"

	//  TLS-certs
	cert, err := tls.LoadX509KeyPair(
		"/etc/letsencrypt/live/subscrsrv.boomag.org/fullchain.pem",
		"/etc/letsencrypt/live/subscrsrv.boomag.org/privkey.pem")
	if err != nil {
		log.Fatalf("failed to load TLS certs: %v", err)
	}

	// Setup TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	logger, err := logger.NewLogger()
	if err != nil {
		// Panic
		log.Panicf("failed to create logger: %v", err)
	}

	// Init dependencies
	deps := &deps.Deps{
		Storage:      localStorage,
		Logger:       logger,
		AppleService: appstore.NewAppleStoreService(localStorage, logger),
	}

	// HTTP server
	server := &http.Server{
		Addr:      port,
		Handler:   httpTransport.NewRouter(localStorage),
		TLSConfig: tlsConfig,
	}

	fmt.Println("Starting server on https://localhost" + port)

	// launch HTTPS
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
