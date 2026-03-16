package dao

import (
	"InnerG/dao/cache"
	"InnerG/dao/db"
	"InnerG/dao/mongo"
)

func Init() {
	db.InitMySQL()
	cache.InitCache()
	mongo.InitMongoDb()
}
