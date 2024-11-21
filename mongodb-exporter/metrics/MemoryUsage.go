package metrics

import "github.com/prometheus/client_golang/prometheus"

// 定义 Prometheus 指标
var (
	MongoDBMemoryResident = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "mongodb_memory_resident_bytes",
			Help: "MongoDB resident memory usage in bytes",
		},
	)
	MongoDBMemoryVirtual = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "mongodb_memory_virtual_bytes",
			Help: "MongoDB virtual memory usage in bytes",
		},
	)
)

func init() {
	// 注册 Prometheus 指标
	prometheus.MustRegister(MongoDBMemoryResident)
	prometheus.MustRegister(MongoDBMemoryVirtual)
}