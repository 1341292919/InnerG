package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Account          string
	Username         string
	Email            string
	PhoneNumber      string
	Avatar           string
	Status           int8
	PasswordHash     string
	AccountChangedAt time.Time
}

const (
	PassWordCost      = 12 // 密码加密难度
	Active       int8 = 1  // 激活用户
)

// SetPassword 设置密码
func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(bytes)
	return nil
}

// CheckPassword 校验密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// SetDefaultAvatar 设置默认头像
func (u *User) SetDefaultAvatar() {
	u.Avatar = constants.DefaultAvatarUrl
}
