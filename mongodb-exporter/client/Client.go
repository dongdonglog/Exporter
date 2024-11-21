package client

import (
	"context"
	"fmt"
	"log"
	"os" // 导入 os 包来读取环境变量

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

// Connect 连接到 MongoDB（确保连接只初始化一次）
func Connect() error {
	if Client != nil {
		return nil // 已经连接，无需重复连接
	}

	// 从环境变量中获取 MongoDB 连接字符串
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return fmt.Errorf("environment variable MONGODB_URI is not set")
	}

	// 设置连接选项
	clientOptions := options.Client().ApplyURI(mongoURI)
	var err error
	Client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	log.Println("Successfully connected to MongoDB")
	return nil
}
