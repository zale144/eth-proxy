package main

import (
	"net"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

func PanicRecoveryMiddleware(next httprouter.Handle, log *zap.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Panic", zap.Any("error", err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next(w, r, p)
	}
}

// RateLimitMiddleware rate limits incoming requests by IP
// note: store could be an interface in case we want to use a different implementation (e.g. Redis)
func RateLimitMiddleware(store *RateLimiterStore, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		limiter := store.getLimiter(host)
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next(w, r, p)
	}
}
