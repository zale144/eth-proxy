package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/zale144/eth-proxy/docs"
	"go.uber.org/zap"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/http-swagger"
)

// Router creates a new HTTP router for the service
func Router(store *RateLimiterStore, handler Handler) *httprouter.Router {
	router := httprouter.New()
	router.GET("/eth/balance/:address", PanicRecoveryMiddleware(RateLimitMiddleware(store, handler.GetBalance), handler.log))
	router.Handler("GET", "/metrics", promhttp.Handler())
	router.GET("/healthy", handler.Healthy)
	router.GET("/ready", handler.Ready)
	router.Handler("GET", "/swagger/*any", httpSwagger.WrapHandler)
	return router
}

// StartHTTPServer starts the HTTP server
func StartHTTPServer(router *httprouter.Router, port string, timeoutSec int, log *zap.Logger) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + port,
		WriteTimeout: time.Duration(timeoutSec) * time.Second,
		ReadTimeout:  time.Duration(timeoutSec) * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("Httpserver: ListenAndServe() error", zap.Error(err))
		} else {
			log.Info("Httpserver: ListenAndServe() started")
		}
	}()

	log.Info("The service is ready to listen and serve.", zap.String("port", port))

	switch killSignal := <-interrupt; killSignal {
	case os.Interrupt:
		log.Debug("Got SIGINT...")
	case syscall.SIGTERM:
		log.Debug("Got SIGTERM...")
	}

	log.Info("The service is shutting down...")

	if err := srv.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("HTTP server Shutdown: %w", err)
	}
	log.Info("HTTP server stopped")
	return nil
}
