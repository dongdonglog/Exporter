package monitor

import (
	"context"
	"log"
	"time"

	"mongodb_exporter/client"
	"mongodb_exporter/metrics"
)

// CheckMongoDBHealth periodically checks the health of the MongoDB instance.
func CheckMongoDBHealth() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in CheckMongoDBHealth: %v", r)
			metrics.MongoDBHealthStatus.Set(0) // 设置为不可用
		}
	}()

	if client.Client == nil {
		log.Println("MongoDB client is not initialized")
		metrics.MongoDBHealthStatus.Set(0) // 设置为不可用
		return
	}

	// 发送 ping 命令测试连接
	err := client.Client.Ping(context.Background(), nil)
	if err != nil {
		log.Printf("MongoDB health check failed: %v", err)
		metrics.MongoDBHealthStatus.Set(0) // 设置为不可用
		return
	}

	// 如果成功，设置为 1
	//log.Println("MongoDB is alive")
	metrics.MongoDBHealthStatus.Set(1)
}

// StartHealthCheck starts a periodic health check for MongoDB.
func StartHealthCheck(interval time.Duration) {
	go func() {
		for {
			CheckMongoDBHealth()
			time.Sleep(interval)
		}
	}()
}
