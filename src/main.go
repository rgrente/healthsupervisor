package main

import (
	"context"
	"healthsupervisor/internal/config"
	"healthsupervisor/internal/httpServer"
	"healthsupervisor/internal/prober"
	"healthsupervisor/internal/rules"
	"healthsupervisor/internal/supervisor"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// Load configuration
	c, err := config.NewConfig("/config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Parse rules from config
	rules := rules.ParseRules(c.Rules)

	// Create prober and supervisor
	prober := prober.NewProber(c.Probes)
	supervisor, err := supervisor.NewSupervisor("supervisor", 1, c.RemoteSupervisors, prober, rules, c.Hooks)
	if err != nil {
		log.Fatal(err)
	}

	// Use a WaitGroup to manage goroutines
	var wg sync.WaitGroup
	wg.Add(3)

	// Start prober and supervisor concurrently
	go func() {
		defer wg.Done()
		log.Println("Starting Prober")
		prober.Run()
	}()

	go func() {
		defer wg.Done()
		log.Println("Starting Supervisor")
		supervisor.Run()
	}()

	// Set up HTTP server with a router
	router := http.NewServeMux()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		httpServer.HealthHandler(w, r, supervisor, prober)
	})

	srv := &http.Server{
		Addr:         ":80",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start HTTP server in a goroutine
	go func() {
		defer wg.Done()
		log.Println("Starting HTTP server on :80")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Set up graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	log.Println("Server stopped")
}
