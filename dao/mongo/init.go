package mongo

import (
	"InnerG/config"
	_interface "InnerG/dao/interface"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"time"
)

var MongoDB *mongo.Database

// InitMongoDb 初始化 MongoDB 连接
func InitMongoDb() {
	clientOptions := options.Client().ApplyURI(MongoDbDSN())
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		panic(err)
	}
	MongoDB = client.Database(config.MongoDb.Database)
	log.Println("✅ MongoDB 连接成功！")
}

func MongoDbDSN() string {
	log.Println(config.MongoDb.Username)
	log.Println(config.MongoDb.Password)
	log.Println(config.MongoDb.Addr)
	log.Println(config.MongoDb.Database)
	return fmt.Sprintf(
		"mongodb://%s:%s@%s/%s?authSource=%s",
		config.MongoDb.Username, config.MongoDb.Password,
		config.MongoDb.Addr,
		config.MongoDb.Database,
		config.MongoDb.Database,
	)
}

// CloseMongoDb 关闭连接（程序退出时调用）
func CloseMongoDb(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Printf("关闭 MongoDB 连接失败: %v", err)
	}
}

func NewContactMongoDBClient() _interface.ContactMongoDB {
	return NewContactMongoDB(MongoDB)
}
