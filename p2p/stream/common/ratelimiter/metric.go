package ratelimiter

import (
	prom "github.com/PositionExchange/posichain/api/service/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prom.PromRegistry().MustRegister(
		serverRequestCounter,
		serverRequestDelayDuration,
	)
}

var (
	serverRequestCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "psc",
			Subsystem: "stream",
			Name:      "num_server_request",
			Help:      "number of incoming requests as server",
		},
	)

	serverRequestDelayDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "psc",
			Subsystem: "stream",
			Name:      "server_request_delay",
			Help:      "delay in seconds of incoming requests of server",
			Buckets:   prometheus.ExponentialBuckets(0.01, 2, 5),
		},
	)
)
