package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	httpTransport "subscription-server/internal/transport/http"
)

func main() {

	port := ":443"
	// 13.
	// Загружаем TLS-сертификаты
	cert, err := tls.LoadX509KeyPair(
		"/etc/letsencrypt/live/subscrsrv.boomag.org/fullchain.pem",
		"/etc/letsencrypt/live/subscrsrv.boomag.org/privkey.pem")
	if err != nil {
		log.Fatalf("failed to load TLS certs: %v", err)
	}

	// Настройки TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	// HTTP сервер
	server := &http.Server{
		Addr:      port,
		Handler:   httpTransport.NewRouter(),
		TLSConfig: tlsConfig,
	}

	fmt.Println("🔐 Starting server on https://localhost" + port)

	// Запуск HTTPS
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
