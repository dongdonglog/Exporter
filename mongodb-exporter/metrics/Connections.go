package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// MongoDB 连接数
	MongoDBConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "mongodb_connections_current",
			Help: "Current number of active connections",
		},
	)
)

func init() {
	prometheus.MustRegister(MongoDBConnections)
}
