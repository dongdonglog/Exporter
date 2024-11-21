package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	MongoDBDatabaseDataSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mongodb_database_data_size_bytes",
			Help: "MongoDB database data size in bytes",
		},
		[]string{"database"},
	)

	MongoDBDatabaseLogSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mongodb_database_log_size_bytes",
			Help: "MongoDB database log size in bytes",
		},
		[]string{"database"},
	)

	MongoDBDatabaseInsSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mongodb_database_ins_size_bytes",
			Help: "MongoDB database ins size (dataSize + logSize) in bytes",
		},
		[]string{"database"},
	)

	MongoDBDatabaseDiskPercentage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mongodb_database_disk_percentage",
			Help: "MongoDB database disk usage percentage relative to total disk size",
		},
		[]string{"database"},
	)
)

func init() {
	prometheus.MustRegister(
		MongoDBDatabaseDataSize,
		MongoDBDatabaseLogSize,
		MongoDBDatabaseInsSize,
		MongoDBDatabaseDiskPercentage,
	)
}
