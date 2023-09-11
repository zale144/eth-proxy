package main

import (
	"sync"

	"golang.org/x/time/rate"
)

// RateLimiterStore is a synchronized store for rate limiters
type RateLimiterStore struct {
	sync.RWMutex
	rateLimiters map[string]*rate.Limiter
	limit        rate.Limit
	burst        int
}

// NewRateLimiterStore creates a new RateLimiterStore
func NewRateLimiterStore(limit float32, burst int) *RateLimiterStore {
	return &RateLimiterStore{
		rateLimiters: make(map[string]*rate.Limiter),
		limit:        rate.Limit(limit),
		burst:        burst,
	}
}

// Retrieve or create a new rate limiter for the given IP
func (s *RateLimiterStore) getLimiter(ip string) *rate.Limiter {
	s.Lock()
	defer s.Unlock()

	limiter, exists := s.rateLimiters[ip]
	if !exists {
		limiter = rate.NewLimiter(s.limit, s.burst)
		s.rateLimiters[ip] = limiter
	}
	return limiter
}
