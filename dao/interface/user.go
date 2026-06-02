package _interface

import (
	"InnerG/dao/db/model"
	"context"
)

type UserDB interface {
	CreateNewUser(ctx context.Context, user *model.User) error
	IsUserExistById(ctx context.Context, id int64) (*model.User, bool, error)
	GetUserBasicByIds(ctx context.Context, ids []int64) (map[int64]*model.User, error)
	IsUserExistByEmail(ctx context.Context, email string) (*model.User, bool, error)
	IsUserExistByAccount(ctx context.Context, account string) (*model.User, bool, error)
	UpdateUserAccount(ctx context.Context, account string, id int64) error
	UpdateUserName(ctx context.Context, userName string, id int64) error
	UpdateUserGender(ctx context.Context, gender string, id int64) error
	UpdateUserAvatar(ctx context.Context, id, avatarUrl string) error
}
type UserCache interface {
	IsKeyExist(ctx context.Context, key string) bool
	SetEmailCode(ctx context.Context, key string, code string) error
	GetEmailCode(ctx context.Context, key string) (string, error)
	BlockToken(ctx context.Context, token string) error
}
