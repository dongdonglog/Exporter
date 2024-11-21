package monitor

import (
	"context"
	"fmt"
	"log"
	"mongodb_exporter/client"
	"mongodb_exporter/metrics"

	"go.mongodb.org/mongo-driver/bson"
)

// UpdateMemoryMetrics 更新内存相关指标
func UpdateMemoryMetrics() error {
	var serverStatus bson.M
	err := client.Client.Database("admin").RunCommand(context.Background(), bson.D{{"serverStatus", 1}}).Decode(&serverStatus)
	if err != nil {
		return fmt.Errorf("failed to fetch serverStatus: %v", err)
	}
	if memInfo, ok := serverStatus["mem"].(bson.M); ok {
		var resident, virtual int64
	
		// 解析 resident
		if resVal, ok := memInfo["resident"]; ok {
			switch v := resVal.(type) {
			case int32:
				resident = int64(v) * 1024 * 1024 // MB 转换为字节
			case int64:
				resident = v * 1024 * 1024
			case float64:
				resident = int64(v * 1024 * 1024)
			default:
				log.Printf("Unknown resident type: %T", v)
			}
		}
	
		// 解析 virtual
		if virtVal, ok := memInfo["virtual"]; ok {
			switch v := virtVal.(type) {
			case int32:
				virtual = int64(v) * 1024 * 1024 // MB 转换为字节
			case int64:
				virtual = v * 1024 * 1024
			case float64:
				virtual = int64(v * 1024 * 1024)
			default:
				log.Printf("Unknown virtual type: %T", v)
			}
		}
	
		// 更新 Prometheus 指标
		metrics.MongoDBMemoryResident.Set(float64(resident))
		metrics.MongoDBMemoryVirtual.Set(float64(virtual))
	
		// 日志输出
		//log.Printf(" 物理内存: %dMB, 虚拟内存: %dMB", resident/1024/1024, virtual/1024/1024)
	} else {
		log.Println("Memory information not available in serverStatus")
	}
	

	return nil
}