package main

import "github.com/prometheus/client_golang/prometheus"

var (
	ethBalanceRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "eth_balance_requests_total",
			Help: "Total number of requests to /eth/balance",
		},
		[]string{"status"},
	)
)

func init() {
	prometheus.MustRegister(ethBalanceRequests)
}
