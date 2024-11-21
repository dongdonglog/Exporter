package utils

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 提取字符串值
func GetStringValue(data bson.M, key string, defaultValue string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return defaultValue
}

func GetIntValue(data bson.M, key string, defaultValue int) int {
	if value, ok := data[key]; ok {
		switch v := value.(type) {
		case int:
			return v
		case int32:
			return int(v)
		case int64:
			return int(v)
		case float64:
			return int(v)
		default:
			log.Printf("Unsupported numeric type for key '%s': %T", key, v)
		}
	}
	return defaultValue
}

// 提取嵌套字段的整数值
func GetNestedIntValue(data bson.M, path string, defaultValue int) int {
	keys := strings.Split(path, ".")
	var current interface{} = data
	for _, key := range keys {
		if nested, ok := current.(bson.M); ok {
			current = nested[key]
		} else {
			return defaultValue
		}
	}
	if value, ok := current.(int); ok {
		return value
	}
	return defaultValue
}

func ExtractFilter(slowOperation bson.M) string {
	if command, ok := slowOperation["command"].(bson.M); ok {
		// 处理查询操作的过滤器
		if filter, ok := command["filter"]; ok {
			filterJSON, err := json.Marshal(filter)
			if err != nil {
				log.Printf("Failed to serialize filter: %v", err)
				return "unknown"
			}
			return string(filterJSON)
		}

		// 处理更新和删除操作的过滤器
		if q, ok := command["q"]; ok {
			qJSON, err := json.Marshal(q)
			if err != nil {
				log.Printf("Failed to serialize q: %v", err)
				return "unknown"
			}
			return string(qJSON)
		}

		// 没有 filter 或 q
		log.Printf("No 'filter' or 'q' key found in 'command': %+v", command)
	} else {
		log.Printf("'command' key missing or invalid in slowOperation: %+v", slowOperation)
	}
	return "unknown"
}



// 提取执行阶段统计信息
func ExtractExecStats(slowOperation bson.M) string {
	if execStats, ok := slowOperation["execStats"].(bson.M); ok {
		execStatsJSON, _ := json.Marshal(execStats)
		return string(execStatsJSON)
	}
	return "unknown"
}
// ParseNumericValue 解析各种数值类型为 int64
func ParseNumericValue(value interface{}) int64 {
	switch v := value.(type) {
	case int32:
		return int64(v)
	case int64:
		return v
	case float64:
		return int64(v)
	default:
		//log.Printf("Unknown numeric type: %T", v)
		return 0
	}
}
func GetTimeValue(data bson.M, key string, defaultValue time.Time) time.Time {
	if value, ok := data[key]; ok {
		switch v := value.(type) {
		case time.Time:
			return v
		case primitive.DateTime: // 处理 MongoDB 的日期时间类型
			return v.Time()
		default:
			log.Printf("Unsupported timestamp type: %T", v)
		}
	}
	return defaultValue
}