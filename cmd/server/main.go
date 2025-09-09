package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	appstore "subscription-server/internal/applestore"
	"subscription-server/internal/deps"
	"subscription-server/internal/logger"
	"subscription-server/internal/storage"
	httpTransport "subscription-server/internal/transport/http"
	"syscall"
	"time"
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

	validator := appstore.NewAppleJWSValidator()
	decoder := appstore.NewAppleDecoder(validator)
	parser := appstore.NewAppleParser(decoder)

	// Init dependencies
	deps := &deps.Deps{
		Storage:      localStorage,
		Logger:       logger,
		AppleService: appstore.NewAppleStoreService(localStorage, logger, parser),
	}

	// HTTP server
	server := &http.Server{
		Addr:      port,
		Handler:   httpTransport.NewRouter(deps),
		TLSConfig: tlsConfig,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fmt.Println("Starting server on https://localhost" + port)
	go func() {
		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
	fmt.Println("Server exited properly")
	logger.Close()

}
