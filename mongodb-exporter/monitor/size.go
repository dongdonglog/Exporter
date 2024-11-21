package monitor

import (
	"context"
	"fmt"
	"log"
	"mongodb_exporter/client"
	"mongodb_exporter/metrics"
	"mongodb_exporter/utils"

	"go.mongodb.org/mongo-driver/bson"
)

// UpdateCoreDatabaseMetrics 更新核心数据库(core)的详细使用量
func UpdateCoreDatabaseMetrics() error {
	databaseName := "core" // 核心数据库名称

	// 获取数据库统计信息 (DataSize)
	var dbStats bson.M
	err := client.Client.Database(databaseName).RunCommand(context.Background(), bson.D{{"dbStats", 1}}).Decode(&dbStats)
	if err != nil {
		return fmt.Errorf("failed to fetch dbStats for %s: %v", databaseName, err)
	}

	// 提取数据大小 (DataSize)
	dataSize := utils.ParseNumericValue(dbStats["dataSize"])

	// 获取日志大小 (LogSize) 从 serverStatus 的 wiredTiger.log 中提取
	logSize := fetchLogSize()

	// 计算 InsSize
	insSize := dataSize + logSize

	// 获取系统磁盘总大小
	fsTotalSize := utils.ParseNumericValue(dbStats["fsTotalSize"])

	// 计算磁盘占用比例
	var diskPercentage float64
	if fsTotalSize > 0 {
		diskPercentage = (float64(insSize) / float64(fsTotalSize)) * 100
	}

	// 更新 Prometheus 指标
	metrics.MongoDBDatabaseDataSize.WithLabelValues(databaseName).Set(float64(dataSize))
	metrics.MongoDBDatabaseLogSize.WithLabelValues(databaseName).Set(float64(logSize))
	metrics.MongoDBDatabaseInsSize.WithLabelValues(databaseName).Set(float64(insSize))
	metrics.MongoDBDatabaseDiskPercentage.WithLabelValues(databaseName).Set(diskPercentage)

	// 打印日志
	// log.Printf(
	// 	"Database: %s - DataSize: %.2fMB, LogSize: %.2fMB, InsSize: %.2fMB, Disk Usage: %.2f%%",
	// 	databaseName, float64(dataSize)/1024/1024, float64(logSize)/1024/1024, float64(insSize)/1024/1024, diskPercentage,
	// )

	return nil
}

// fetchLogSize 从 serverStatus 中获取日志大小
func fetchLogSize() int64 {
	var serverStatus bson.M
	err := client.Client.Database("admin").RunCommand(context.Background(), bson.D{{"serverStatus", 1}}).Decode(&serverStatus)
	if err != nil {
		log.Printf("Failed to fetch serverStatus: %v", err)
		return 0
	}

	// 提取 wiredTiger.log 信息
	if wiredTiger, ok := serverStatus["wiredTiger"].(bson.M); ok {
		if logDetails, ok := wiredTiger["log"].(bson.M); ok {
			return utils.ParseNumericValue(logDetails["log bytes written"])
		}
	}

	return 0
}
