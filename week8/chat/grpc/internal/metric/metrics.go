package metric

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "my_space"
	appName   = "my_app"
)

type Metrics struct {
	requestCounter  prometheus.Counter
	responseCounter *prometheus.CounterVec
}

var metrics *Metrics

func Init(_ context.Context) error {
	metrics = &Metrics{
		requestCounter: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "grpc",
				Name:      appName + "_request_total",
				Help:      "Количество запросов к серверу",
			},
		),
		responseCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "grpc",
				Name:      appName + "_response_total",
				Help:      "Количество ответов к серверу",
			},
			[]string{"status", "method"},
		),
	}
	prometheus.MustRegister(metrics.requestCounter)
	prometheus.MustRegister(metrics.responseCounter)

	return nil
}
func IncRequestCounter() {
	metrics.requestCounter.Inc()
}
func IncResponseCounter(status string, method string) {
	metrics.responseCounter.WithLabelValues(status, method).Inc()
}
