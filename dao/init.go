package dao

import (
	"InnerG/dao/cache"
	"InnerG/dao/db"
	"InnerG/dao/mongo"
	"InnerG/dao/rabbitmq"
)

func Init() {
	db.InitMySQL()
	cache.InitCache()
	mongo.InitMongoDb()
	rabbitmq.Init()
}
