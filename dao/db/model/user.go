package model

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Account      sql.NullString
	Username     string
	Email        string
	Avatar       string
	Status       int8
	PasswordHash string
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
}
