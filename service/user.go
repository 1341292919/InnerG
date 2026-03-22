package service

import (
	"InnerG/dao"
	"InnerG/dao/db/model"
	"InnerG/pkg/constants"
	"InnerG/pkg/ctl"
	"InnerG/pkg/oss"
	"InnerG/pkg/utils"
	"InnerG/types"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
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
	log.Println(req.Email, "发送验证码:", code)
	return nil
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
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("该邮箱未绑定账号，请先注册")
	}
	return u, nil
}

func (s *UserSrv) GetUserInfo(ctx context.Context) (*model.User, error) {
	u := ctl.GetUserInfo(ctx)
	userDao := dao.NewUserDao(ctx)
	user, exist, err := userDao.Db.IsUserExistById(ctx, u.Id)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("user not exist")
	}
	return user, nil
}

func (s *UserSrv) UpdateUserAccount(ctx context.Context, account string) error {
	u := ctl.GetUserInfo(ctx)
	userDao := dao.NewUserDao(ctx)
	user, exist, err := userDao.Db.IsUserExistByAccount(ctx, account)
	if err != nil {
		return err
	}
	if exist {
		if u.Id == strconv.FormatInt(int64(user.ID), 10) {
			return fmt.Errorf("新账号与原账号相同，无需修改")
		}
		return fmt.Errorf("该账号已存在")
	}
	return userDao.Db.UpdateUserAccount(ctx, account, u.Id)
}

func (s *UserSrv) UpdateUserName(ctx context.Context, userName string) error {
	u := ctl.GetUserInfo(ctx)
	userDao := dao.NewUserDao(ctx)
	user, exist, err := userDao.Db.IsUserExistById(ctx, u.Id)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("用户不存在")
	}

	userName = strings.TrimSpace(userName)
	if userName == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if user.Username == userName {
		return fmt.Errorf("新用户名与原用户名相同，无需修改")
	}
	return userDao.Db.UpdateUserName(ctx, userName, u.Id)
}

func (s *UserSrv) UpdateUserGender(ctx context.Context, gender string) error {
	u := ctl.GetUserInfo(ctx)
	userDao := dao.NewUserDao(ctx)
	_, exist, err := userDao.Db.IsUserExistById(ctx, u.Id)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("用户不存在")
	}

	if gender != "0" && gender != "1" {
		return fmt.Errorf("gender 参数错误，仅支持 0 或 1")
	}
	return userDao.Db.UpdateUserGender(ctx, gender, u.Id)
}

func (s *UserSrv) LogOut(ctx context.Context) error {
	u := ctl.GetUserInfo(ctx)
	userDao := dao.NewUserDao(ctx)
	key := fmt.Sprintf("token:%s", u.Token)
	return userDao.Cache.BlockToken(ctx, key)
}
func (s *UserSrv) UpdateUserAvatar(ctx context.Context, file *multipart.FileHeader) (string, error) {
	u := ctl.GetUserInfo(ctx)
	userDao := dao.NewUserDao(ctx)
	err := oss.IsImage(file)
	if err != nil {
		return "", fmt.Errorf("check image failed: %w", err)
	}
	// 识别图片信息
	fileName := fmt.Sprintf("%v_%v", u.Id, time.Now().Unix())
	err = oss.SaveFile(file, constants.StorePath, fileName)
	if err != nil {
		return "", fmt.Errorf("save file failed: %w", err)
	}
	filePath := filepath.Join(constants.StorePath, fileName)
	url, err := oss.Upload(filePath, fileName, u.Id, constants.OssOrigin)
	return url, userDao.Db.UpdateUserAvatar(ctx, u.Id, url)
}
func IsEmail(str string) bool {
	return strings.Contains(str, "@")
}
