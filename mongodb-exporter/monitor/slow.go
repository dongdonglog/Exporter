package monitor

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"mongodb_exporter/client"
	"mongodb_exporter/metrics"
	"mongodb_exporter/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	startupTime       time.Time // 程序启动时间
	lastProcessedTime time.Time // 上次处理的最早时间戳
	mu                sync.Mutex
)

// InitializeStartupTime 初始化 startupTime 和 lastProcessedTime
func InitializeStartupTime() error {
	var serverStatus bson.M
	err := client.Client.Database("admin").RunCommand(context.Background(), bson.D{{"serverStatus", 1}}).Decode(&serverStatus)
	if err != nil {
		return fmt.Errorf("failed to fetch serverStatus: %v", err)
	}

	// 获取 MongoDB 的 localTime
	localTime, ok := serverStatus["localTime"].(primitive.DateTime)
	if !ok {
		return fmt.Errorf("failed to parse localTime from serverStatus")
	}

	startupTime = localTime.Time()
	lastProcessedTime = startupTime // 初始化时查询未来 7 天的数据
	log.Printf("Exporter initialized with MongoDB localTime: %v", startupTime)
	return nil
}

func FetchFutureSlowOperations() error {
    mu.Lock()
    queryStartTime := lastProcessedTime
    queryEndTime := queryStartTime.Add(7 * 24 * time.Hour) // 查询未来 7 天
    mu.Unlock()

    if time.Now().After(queryEndTime) { // 如果当前时间超过了查询范围
        return nil // 跳过此次查询
    }

    // 构造查询
    query := bson.D{
        {"millis", bson.D{{"$gt", 100}}}, // 慢操作阈值
        {"ts", bson.D{
            {"$gt", primitive.NewDateTimeFromTime(queryStartTime)},
            {"$lte", primitive.NewDateTimeFromTime(queryEndTime)},
        }},
    }

    cursor, err := client.Client.Database("core").Collection("system.profile").Find(
        context.Background(),
        query,
    )
    if err != nil {
        return fmt.Errorf("failed to query system.profile: %v", err)
    }
    defer cursor.Close(context.Background())

    var maxTimestamp time.Time

    for cursor.Next(context.Background()) {
        var slowOperation bson.M
        if err := cursor.Decode(&slowOperation); err != nil {
            log.Printf("Failed to decode slow operation: %v", err)
            continue
        }

        // 提取时间戳
        timestamp := utils.GetTimeValue(slowOperation, "ts", time.Time{})
        if timestamp.After(maxTimestamp) {
            maxTimestamp = timestamp
        }

        // 提取基础信息
        operation := utils.GetStringValue(slowOperation, "op", "unknown")
        collection := utils.GetStringValue(slowOperation, "ns", "unknown")
        duration := utils.GetIntValue(slowOperation, "millis", 0)
        clientIP := utils.GetStringValue(slowOperation, "client", "unknown")
        user := utils.GetStringValue(slowOperation, "user", "unknown")
        appName := utils.GetStringValue(slowOperation, "appName", "unknown")
        planSummary := utils.GetStringValue(slowOperation, "planSummary", "unknown")

        // 打印简化后的日志
        log.Printf(
            "Future Slow Operation Detected: Operation: %s, Collection: %s, Duration: %dms, Filter: %s, PlanSummary: %s, Client: %s, User: %s, App: %s",
            operation, collection, duration, utils.ExtractFilter(slowOperation), planSummary, clientIP, user, appName,
        )

        // 更新 Prometheus 指标
        metrics.MongoDBSlowOperations.WithLabelValues(
            operation, collection, utils.ExtractFilter(slowOperation), fmt.Sprintf("%d", duration), planSummary, clientIP, user, appName,
        ).Set(float64(duration) / 1000)
    }

    // 更新全局处理时间
    if !maxTimestamp.IsZero() {
        mu.Lock()
        lastProcessedTime = maxTimestamp
        mu.Unlock()
    }

    return nil
}
// CleanupExpiredMetrics 清理过期的 Prometheus 指标
func CleanupExpiredMetrics(expiredTime time.Time) {
    log.Printf("Cleaning up metrics before: %v", expiredTime)

    // 遍历已有的 Prometheus 指标，并删除时间戳早于 expiredTime 的指标
    // 示例（需要结合你的缓存实现）：
    metrics.MongoDBSlowOperations.Reset() // 或者更细粒度的指标清理逻辑
}