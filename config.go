package main

// Config is the configuration for the service
type Config struct {
	ETHNetworkURL       string  `envconfig:"ETH_NETWORK_URL" required:"true"`
	HTTPPort            string  `envconfig:"HTTP_PORT" default:"8080"`
	HTTPTimeoutSec      int     `envconfig:"HTTP_TIMEOUT" default:"15"`
	ClientRetries       int     `envconfig:"CLIENT_RETRIES" default:"3"`
	ClientRetryDelaySec int     `envconfig:"CLIENT_RETRY_DELAY" default:"1"`
	RateLimit           float32 `envconfig:"RATE_LIMIT" default:"10.0"`
	RateBurst           int     `envconfig:"RATE_BURST" default:"10"`
}
