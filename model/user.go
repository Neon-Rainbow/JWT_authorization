package model

import "gorm.io/gorm"

// User struct
type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
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
}

type UserRegisterResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}
