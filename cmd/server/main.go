package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ujjwal0563/ip-shakti-sahayak/internal/api"
	"github.com/ujjwal0563/ip-shakti-sahayak/internal/config"
	"github.com/ujjwal0563/ip-shakti-sahayak/internal/database"
)

func main() {
	log.Println("===============================================================")
	log.Println("🌿 IP-SAKTI Sahayak (आईपी-शक्ति सहायक) - Phase 1 Foundation")
	log.Println("   Intelligent Ayurveda IP, ABS & Regulatory Decision Support")
	log.Println("===============================================================")

	// 1. Load Configuration
	cfg := config.Load()
	log.Printf("[Init] Running with Port: %s | Zero-Data-Retention: %t | LLM: %s",
		cfg.Port, cfg.ZeroDataRetention, cfg.LLMProvider)

	// 2. Initialize Database Connection Pool
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := database.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Printf("[Warning] Database pool creation error: %v (running in standalone mode)", err)
	} else if db != nil {
		defer db.Close()
	}

	// 3. Initialize HTTP Server & Routes
	handler := api.NewServer(cfg, db)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.Port),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// 4. Background Server Runner
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("🚀 IP-SAKTI Sahayak server listening on http://localhost:%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	// 5. Graceful Shutdown Listener
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-serverErrors:
		log.Fatalf("[Fatal] Server startup error: %v", err)

	case sig := <-shutdown:
		log.Printf("\n[Shutdown] Received signal %v. Initiating graceful termination...", sig)

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("[Shutdown Error] Server forced shutdown: %v", err)
			_ = srv.Close()
		}

		log.Println("IP-SAKTI Sahayak gracefully stopped.")
	}
}
