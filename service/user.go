package service

import (
	"InnerG/dao"
	"InnerG/dao/db/model"
	"InnerG/pkg/constants"
	"InnerG/pkg/utils"
	"InnerG/types"
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
)

var UserSrvIns *UserSrv
var UserSrvOnce sync.Once

type UserSrv struct{} // 空结构体，只包含方法

func GetUserSrv() *UserSrv {
	UserSrvOnce.Do(func() {
		UserSrvIns = &UserSrv{}
	})
	return UserSrvIns
}

func (s *UserSrv) GetEmailCode(ctx context.Context, req *types.UserGetEmailCodeReq) error {
	userDao := dao.NewUserDao(ctx)
	emailCodeKey := fmt.Sprintf("emailCode:%s", req.Email)
	code := utils.GenerateRandomCode(constants.CommonEmailCodeLength)
	err := userDao.Cache.SetEmailCode(ctx, emailCodeKey, code)
	if err != nil {
		return err
	}
	// 发送验证码
	return utils.MailSendCode(req.Email, code)
}

func (s *UserSrv) VerifyEmailAndRegister(ctx context.Context, req *types.UserVerifyEmailAndRegisterReq) error {
	userDao := dao.NewUserDao(ctx)
	emailCodeKey := fmt.Sprintf("emailCode:%s", req.Email)
	if !userDao.Cache.IsKeyExist(ctx, emailCodeKey) {
		return fmt.Errorf("验证码错误")
	}
	code, err := userDao.Cache.GetEmailCode(ctx, emailCodeKey)
	if err != nil {
		return err
	}
	if code != req.VerifyCode {
		return fmt.Errorf("验证码错误")
	}

	_, exist, err := userDao.Db.IsUserExistByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("该邮箱已注册其他账号")
	}

	// 插入数据库
	newUser := &model.User{
		Email:  req.Email,
		Status: model.Active,
	}
	newUser.SetDefaultAvatar()
	if err = newUser.SetPassword(req.Password); err != nil {
		return err
	}

	// 默认头像等
	return userDao.Db.CreateNewUser(ctx, newUser)

}

// Login 同时支持邮箱或者账号登录
func (s *UserSrv) Login(ctx context.Context, req *types.UserLoginReq) (*model.User, error) {
	userDao := dao.NewUserDao(ctx)
	var u *model.User
	var exist bool
	var err error
	switch {
	case IsEmail(req.Account):
		u, exist, err = userDao.Db.IsUserExistByEmail(ctx, req.Account)
	default:
		u, exist, err = userDao.Db.IsUserExistByAccount(ctx, req.Account)
	}
	if err != nil {
		return nil, err
	}
	if !exist || !u.CheckPassword(req.Password) {
		return nil, fmt.Errorf("账号或密码错误")
	}
	return u, nil
}

func (s *UserSrv) VerifyEmailAndLogin(ctx context.Context, req *types.UserVerifyEmailAndLoginReq) (*model.User, error) {
	userDao := dao.NewUserDao(ctx)
	emailCodeKey := fmt.Sprintf("emailCode:%s", req.Email)
	if !userDao.Cache.IsKeyExist(ctx, emailCodeKey) {
		return nil, fmt.Errorf("验证码错误")
	}
	code, err := userDao.Cache.GetEmailCode(ctx, emailCodeKey)
	if err != nil {
		return nil, err
	}
	if code != req.VerifyCode {
		return nil, fmt.Errorf("验证码错误")
	}

	u, exist, err := userDao.Db.IsUserExistByEmail(ctx, req.Email)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("该邮箱未绑定账号，请先注册")
	}
	return u, nil
}

func IsEmail(str string) bool {
	return strings.Contains(str, "@")
}
