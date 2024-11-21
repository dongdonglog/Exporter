package monitor

import (
	"context"
	"fmt"
	// "log"
	"mongodb_exporter/client"
	"mongodb_exporter/metrics"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	previousLogBytesWritten   int64
	previousCacheBytesWritten int64
	previousCacheBytesRead    int64
	lastCheckTime             time.Time
)
const blockSize = 4096 // 块大小 4KB

// UpdateIOPSMetrics 监控 IOPS 使用量
func UpdateIOPSMetrics() error {
	var serverStatus bson.M
	err := client.Client.Database("admin").RunCommand(context.Background(), bson.D{{"serverStatus", 1}}).Decode(&serverStatus)
	if err != nil {
		return fmt.Errorf("failed to fetch serverStatus: %v", err)
	}

	// 提取 wiredTiger.metrics 数据
	if wiredTiger, ok := serverStatus["wiredTiger"].(bson.M); ok {
		// 日志写入字节数
		logBytesWritten := parseNumericValue(wiredTiger, "log.log bytes written")
		// 缓存写入字节数
		cacheBytesWritten := parseNumericValue(wiredTiger, "cache.bytes written from cache")
		// 缓存读取字节数
		cacheBytesRead := parseNumericValue(wiredTiger, "cache.bytes read into cache")

		// 当前时间
		currentTime := time.Now()

		// 如果是首次运行，初始化计数器
		if lastCheckTime.IsZero() {
			previousLogBytesWritten = logBytesWritten
			previousCacheBytesWritten = cacheBytesWritten
			previousCacheBytesRead = cacheBytesRead
			lastCheckTime = currentTime
			return nil
		}

		// 时间间隔
		elapsed := currentTime.Sub(lastCheckTime).Seconds()

		// 计算 IOPS 使用量
		dataIOPS := float64(cacheBytesRead+cacheBytesWritten-previousCacheBytesRead-previousCacheBytesWritten) / elapsed / blockSize
		logIOPS := float64(logBytesWritten-previousLogBytesWritten) / elapsed / blockSize
		totalIOPS := dataIOPS + logIOPS

		// 更新 Prometheus 指标
		metrics.MongoDBDataIOPS.Set(dataIOPS)
		metrics.MongoDBLogIOPS.Set(logIOPS)
		metrics.MongoDBTotalIOPS.Set(totalIOPS)

		// 打印日志
		//log.Printf("Data IOPS: %.2f, Log IOPS: %.2f, Total IOPS: %.2f", dataIOPS, logIOPS, totalIOPS)

		// 更新前一次统计
		previousLogBytesWritten = logBytesWritten
		previousCacheBytesWritten = cacheBytesWritten
		previousCacheBytesRead = cacheBytesRead
		lastCheckTime = currentTime
	}

	return nil
}

// parseNumericValue 辅助函数：解析数字
func parseNumericValue(data bson.M, key string) int64 {
	keys := strings.Split(key, ".")
	for _, k := range keys[:len(keys)-1] {
		if nested, ok := data[k].(bson.M); ok {
			data = nested
		} else {
			return 0
		}
	}
	if value, ok := data[keys[len(keys)-1]]; ok {
		switch v := value.(type) {
		case int32:
			return int64(v)
		case int64:
			return v
		case float64:
			return int64(v)
		}
	}
	return 0
}
