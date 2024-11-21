package monitor

import (
	"context"
	"fmt"
	"log"
	"mongodb_exporter/metrics"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UpdateConnectionsMetrics 更新 MongoDB 连接数指标
func UpdateConnectionsMetrics(client *mongo.Client) error {
	var serverStatus bson.M
	err := client.Database("core").RunCommand(context.Background(), bson.D{{"serverStatus", 1}}).Decode(&serverStatus)
	if err != nil {
		return fmt.Errorf("failed to get serverStatus: %v", err)
	}

	connections, ok := serverStatus["connections"].(bson.M)
	if !ok {
		return fmt.Errorf("failed to parse connections field")
	}

	// 获取当前连接数
	if current, ok := connections["current"].(int32); ok {
		metrics.MongoDBConnections.Set(float64(current))
	} else {
		log.Println("Failed to get current connections")
	}

	return nil
}
