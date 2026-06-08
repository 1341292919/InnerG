package pack

import (
	"InnerG/dao/db/model"
	"InnerG/types"
	"strconv"
)

func BuildUser(user *model.User) *types.User {
	return &types.User{
		Id:            strconv.FormatInt(int64(user.ID), 10),
		Email:         user.Email,
		Avatar:        user.Avatar,
		Signature:     user.Signature,
		ProfilePublic: user.ProfilePublic,
		UserName:      user.Username,
		Account:       user.Account.String,
		RoleType:      int(user.RoleType),
		Gender:        int(user.Gender),
		CreatedAt:     user.CreatedAt.Unix(),
		UpdatedAT:     user.UpdatedAt.Unix(),
	}
}

func BuildUserProfile(user *model.User) *types.UserProfile {
	return &types.UserProfile{
		Id:        strconv.FormatInt(int64(user.ID), 10),
		Avatar:    user.Avatar,
		Signature: user.Signature,
		UserName:  user.Username,
		Account:   user.Account.String,
		RoleType:  int(user.RoleType),
		Gender:    int(user.Gender),
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAT: user.UpdatedAt.Unix(),
	}
}

func BuildProtectedUserProfile(user *model.User) *types.UserProfile {
	return &types.UserProfile{
		Id:        strconv.FormatInt(int64(user.ID), 10),
		Avatar:    user.Avatar,
		UserName:  user.Username,
		Protected: true,
	}
}
