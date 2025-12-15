package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chess-training/internal/config"
	db "chess-training/internal/db/sqlc"
	"chess-training/internal/http/router"
	"chess-training/internal/service"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	defer pool.Close()

	queries := db.New(pool)
	authSvc := service.NewAuthService(queries, cfg.JWTSecret, cfg.JWTTTLMin)

	r := router.New(authSvc, cfg.JWTSecret)

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: r,
	}

	go func() {
		log.Printf("api listening on %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen failed: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("shutting down...")
	_ = srv.Shutdown(ctxShutdown)
}
