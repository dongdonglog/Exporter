package metrics

import "github.com/prometheus/client_golang/prometheus"

var(
	MongoDBBytesRead = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mongodb_bytes_read_total",
			Help: "Total bytes read in slow operations",
		},
		[]string{"namespace", "operation"},
	)
)
func init() {
	prometheus.MustRegister(MongoDBBytesRead)
}