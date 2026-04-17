package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API simulation
	mux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from go backend\n"))

		ip := getClientIP(r)
		fmt.Fprintf(w, "Hai from %s", ip)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// jalankan server di goroutine
	go func() {
		log.Println("Go backend running on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Println("Server exited cleanly")
}

func getClientIP(r *http.Request) string {
	fmt.Println("X-Forwarded-For:", r.Header.Get("X-Forwarded-For"))
	fmt.Println("X-Real-IP:", r.Header.Get("X-Real-IP"))
	fmt.Println("RemoteAddr:", r.RemoteAddr)

	// Cek X-Forwarded-For
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}

	// Cek X-Real-IP
	xrip := r.Header.Get("X-Real-IP")
	if xrip != "" {
		return xrip
	}

	// Fallback
	return r.RemoteAddr
}
