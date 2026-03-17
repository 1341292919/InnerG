package pack

import (
	"InnerG/dao/db/model"
	"InnerG/types"
	"strconv"
)

func BuildUser(user *model.User) *types.User {
	return &types.User{
		Id:        strconv.FormatInt(int64(user.ID), 10),
		Email:     user.Email,
		Avatar:    user.Avatar,
		UserName:  user.Username,
		Account:   user.Account,
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAT: user.UpdatedAt.Unix(),
	}
}
