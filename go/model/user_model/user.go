package user_model

import (
	"database/sql"
	"time"
)

type User struct {
	UserID       int64          `json:"user_id"`
	Username     string         `json:"username"`
	PasswordHash string         `json:"password_hash"`
	PhoneNumber  string         `json:"phone_number"`
	AvatarURL    sql.NullString `json:"avatar_url"`
	Email        string         `json:"email"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	RoleID       int            `json:"role_id"`
}
type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Password    string `json:"password" binding:"required,min=6,max=50"`
	PhoneNumber string `json:"phone" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	RoleID      int    `json:"role_id" binding:"required,min=1,max=3"`
}
type RegisterInfo struct {
	Username     string `binding:"required,min=3,max=50"`
	PasswordHash string
	PhoneNumber  string `binding:"required"`
	Email        string `binding:"required,email"`
	RoleID       int    `binding:"required,min=1,max=3"`
}
type RegisterResponse struct {
	Message string `json:"message"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required"`
}
type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"-"`
}
type LoginAuth struct {
	UserID       int64
	PasswordHash string
}
type UpdatePasswordRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	PhoneNumber string `json:"phone" binding:"required"`
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
type UpdatePassWordAuth struct {
	UserID       int64
	PasswordHash string
	PhoneNumber  string
}

type UserInfo struct {
	UserID      int64
	Username    string
	PhoneNumber string
	AvatarURL   sql.NullString
	Email       string
	CreatedAt   time.Time
	RoleID      int
}
type UpdateAvatarRequest struct {
	UserID    int64
	AvatarURL string
}
