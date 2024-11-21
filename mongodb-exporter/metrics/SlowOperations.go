package metrics

import "github.com/prometheus/client_golang/prometheus"

// var (
// 	MongoDBSlowOperations = prometheus.NewGaugeVec(
// 		prometheus.GaugeOpts{
// 			Name: "mongodb_real_time_slow_operations_seconds",
// 			Help: "MongoDB real-time slow operations execution time in seconds",
// 		},
// 		[]string{"operation", "collection", "filter", "duration", "planSummary", "client", "user", "appName"},
// 	)
// 	MongoDBDocsExamined = prometheus.NewGaugeVec(
// 		prometheus.GaugeOpts{
// 			Name: "mongodb_docs_examined_total",
// 			Help: "Total documents examined in slow operations",
// 		},
// 		[]string{"namespace", "operation"},
// 	)

// 	MongoDBKeysExamined = prometheus.NewGaugeVec(
// 		prometheus.GaugeOpts{
// 			Name: "mongodb_keys_examined_total",
// 			Help: "Total keys examined in slow operations",
// 		},
// 		[]string{"namespace", "operation"},
// 	)


// )

// func init() {
// 	prometheus.MustRegister(MongoDBSlowOperations)
// 	prometheus.MustRegister(MongoDBDocsExamined)
// 	prometheus.MustRegister(MongoDBKeysExamined)
// }


var (
	MongoDBSlowOperations = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mongodb_slow_operations",
			Help: "Details about slow operations in MongoDB.",
		},
		[]string{"operation", "collection", "filter", "duration", "planSummary", "clientIP", "user", "appName"},
	)
	
)

func init() {
	prometheus.MustRegister(MongoDBSlowOperations)
	//prometheus.MustRegister(MongoDBDocsExamined)

}
