package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	MongoDBDataIOPS = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mongodb_data_iops",
		Help: "MongoDB data directory IOPS usage",
	})
	MongoDBLogIOPS = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mongodb_log_iops",
		Help: "MongoDB log directory IOPS usage",
	})
	MongoDBTotalIOPS = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mongodb_total_iops",
		Help: "MongoDB total IOPS usage",
	})
)

func init() {
	prometheus.MustRegister(MongoDBDataIOPS)
	prometheus.MustRegister(MongoDBLogIOPS)
	prometheus.MustRegister(MongoDBTotalIOPS)
}
