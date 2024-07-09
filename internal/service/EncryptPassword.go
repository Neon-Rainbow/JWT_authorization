package service

import (
	"JWT_authorization/config"
	"crypto/md5"
	"encoding/hex"
)

// EncryptPassword 用于对密码进行加密,使用了md5加密算法
func EncryptPassword(password string) string {
	// secret 用于存储密码加密的密钥
	var secret = config.GetConfig().PasswordSecret
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(password)))
}
