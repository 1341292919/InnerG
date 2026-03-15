package dao

import (
	"InnerG/dao/cache"
	"InnerG/dao/db"
)

func Init() {
	db.InitMySQL()
	cache.InitCache()
}
