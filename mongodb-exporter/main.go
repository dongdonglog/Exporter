package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"mongodb_exporter/client"
	"mongodb_exporter/monitor"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var cleanupMutex sync.Mutex // 清理逻辑的互斥锁

func main() {
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// 初始化 startupTime
	if err := monitor.InitializeStartupTime(); err != nil {
		log.Fatalf("Failed to initialize startup time: %v", err)
	}

	monitor.StartHealthCheck(10 * time.Second) // 每 10 秒检查一次

	// 定时清理的时间控制
	lastCleanupTime := time.Now()

	go func() {
		for {
			// 每 6 小时清理一次过期指标
			cleanupMutex.Lock()
			if time.Since(lastCleanupTime) >= 6*time.Hour {
				expiredTime := time.Now().Add(-7 * 24 * time.Hour) // 计算过期时间
				monitor.CleanupExpiredMetrics(expiredTime)        // 清理过期指标
				lastCleanupTime = time.Now()                      // 更新清理时间
			}
			cleanupMutex.Unlock()

			// 查询未来 7 天的慢操作
			if err := monitor.FetchFutureSlowOperations(); err != nil {
				log.Printf("Error monitoring MongoDB slow operations: %v", err)
			}

			// 其他监控任务
			if err := monitor.UpdateConnectionsMetrics(client.Client); err != nil {
				log.Printf("Error updating connections metrics: %v", err)
			}
			if err := monitor.UpdateMemoryMetrics(); err != nil {
				log.Printf("Error updating memory metrics: %v", err)
			}
			if err := monitor.UpdateCoreDatabaseMetrics(); err != nil {
				log.Println("Error updating database metrics:", err)
			}
			if err := monitor.UpdateIOPSMetrics(); err != nil {
				log.Println("Error updating IOPS metrics:", err)
			}

			time.Sleep(30 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Starting MongoDB Exporter on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
