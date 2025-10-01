package user_repository

import (
	"biye/model/user_model"

	"context"
)

type UserInterface interface {
	Create(ctx context.Context, user *user_model.RegisterInfo) error
	GetUserLoginForAuth(ctx context.Context, username string) (*user_model.LoginAuth, error)
	GetUserForPassword(ctx context.Context, username string) (*user_model.UpdatePassWordAuth, error)
	UpdatePassword(ctx context.Context, UserID int64, NewPassword string) error
	CheckUserExists(ctx context.Context, username, phone, email string) (bool, bool, bool, error)
	GetUserInfoByID(ctx context.Context, UserID int64) (*user_model.UserInfo, error)
	UpdateUserAvatar(ctx context.Context, UserID int64, AvatarURL string) error
}
