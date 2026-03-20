package _interface

import (
	"InnerG/dao/db/model"
	"context"
)

type UserDB interface {
	CreateNewUser(ctx context.Context, user *model.User) error
	IsUserExistById(ctx context.Context, id string) (*model.User, bool, error)
	IsUserExistByEmail(ctx context.Context, email string) (*model.User, bool, error)
	IsUserExistByAccount(ctx context.Context, account string) (*model.User, bool, error)
	UpdateUserAccount(ctx context.Context, account string, id string) error
	UpdateUserAvatar(ctx context.Context, id, avatarUrl string) error
}
type UserCache interface {
	IsKeyExist(ctx context.Context, key string) bool
	SetEmailCode(ctx context.Context, key string, code string) error
	GetEmailCode(ctx context.Context, key string) (string, error)
	BlockToken(ctx context.Context, token string) error
}
