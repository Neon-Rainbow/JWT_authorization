package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string `gorm:"type:varchar(255);unique;not null"`
	Password  string `gorm:"type:varchar(255)"`
	Telephone string `gorm:"type:varchar(255);unique;not null"`
	IsFrozen  bool
	IsAdmin   bool
}

// UserLoginRequest struct
type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserRegisterRequest struct
type UserRegisterRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Telephone string `json:"telephone"`
}

type UserLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       uint   `json:"user_id"`
	Username     string `json:"username"`
	IsAdmin      bool   `json:"is_admin"`
	IsFrozen     bool   `json:"is_frozen"`
}

type UserRegisterResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}
