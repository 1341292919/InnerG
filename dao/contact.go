package dao

import (
	_interface "InnerG/dao/interface"
	"InnerG/dao/mongo"
	"context"
)

type ContactDao struct {
	Mongo _interface.ContactMongoDB
}

func NewContactDao(ctx context.Context) *ContactDao {
	return &ContactDao{
		Mongo: mongo.NewContactMongoDBClient(),
	}
}
