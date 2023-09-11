package main

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()

	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		logger.Fatal("Failed to process env var", zap.Error(err))
	}

	client, err := NewClient(cfg.ETHNetworkURL, cfg.ClientRetries, cfg.ClientRetryDelaySec, logger)
	if err != nil {
		logger.Fatal("Failed to initialize ETH client", zap.Error(err))
	}

	handler, err := NewHandler(client, logger)
	if err != nil {
		logger.Fatal("Failed to initialize HTTP handler", zap.Error(err))
	}

	store := NewRateLimiterStore(cfg.RateLimit, cfg.RateBurst)
	router := Router(store, handler)

	if err := StartHTTPServer(router, cfg.HTTPPort, cfg.HTTPTimeoutSec, logger); err != nil {
		logger.Fatal("Failed to start HTTP server", zap.Error(err))
	}
}
