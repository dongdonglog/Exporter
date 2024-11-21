package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	MongoDBHealthStatus = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "mongodb_health_status",
			Help: "MongoDB health status (1 for alive, 0 for unreachable)",
		},
	)
)

func init() {
	prometheus.MustRegister(MongoDBHealthStatus)
}
