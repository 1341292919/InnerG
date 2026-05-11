package ctl

import (
	"InnerG/pkg/constants"
	"context"
)

var Key = constants.ContextIdKey

type UserInfo struct {
	Id    int64
	Token string
}

func GetUserInfo(ctx context.Context) *UserInfo {
	u, _ := ctx.Value(Key).(*UserInfo)
	return u
}
func NewContext(ctx context.Context, u *UserInfo) context.Context {
	return context.WithValue(ctx, Key, u)
}
