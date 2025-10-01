package user_services

import (
	"biye/User/user_repository"
	"biye/model/user_model"
	"biye/share/error_code"
	"biye/share/hash"
	"biye/share/utils"
	"context"
)

type UserService struct {
	userRepo user_repository.UserInterface
}

func NewUserService(userRepo user_repository.UserInterface) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, req *user_model.RegisterRequest) (*user_model.RegisterResponse, error) {
	usernameExists, phoneExists, emailExists, err := s.userRepo.CheckUserExists(
		ctx,
		req.Username,
		req.PhoneNumber,
		req.Email,
	)
	if err != nil {
		return nil, err
	}
	if usernameExists {
		return nil, error_code.UserExists.WithDetail("用户名已存在")
	}
	if phoneExists {
		return nil, error_code.UserNumberExists.WithDetail("手机号已存在")
	}
	if emailExists {
		return nil, error_code.UserEmailExists.WithDetail("邮箱已存在")
	}
	passwordHash, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, error_code.ServerError.WithDetail("密码加密失败:" + err.Error())
	}

	user := &user_model.RegisterInfo{
		Username:     req.Username,
		PasswordHash: passwordHash,
		PhoneNumber:  req.PhoneNumber,
		Email:        req.Email,
		RoleID:       req.RoleID,
	}
	err = s.userRepo.Create(ctx, user)

	if err != nil {
		return nil, err
	}
	return &user_model.RegisterResponse{
		Message: "注册成功。",
	}, nil
}
func (s *UserService) LoginUser(ctx context.Context, req *user_model.LoginRequest) (*user_model.LoginResponse, error) {
	user, err := s.userRepo.GetUserLoginForAuth(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if !hash.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, error_code.PasswordFail.WithDetail("密码错误。")
	}
	token, err := utils.GenerateToken(user.UserID)
	if err != nil {
		return nil, err
	}

	response := &user_model.LoginResponse{
		Token:   token,
		Message: "登录成功。",
	}
	return response, nil
}
func (s *UserService) UpdatePassword(ctx context.Context, req *user_model.UpdatePasswordRequest) error {
	Info, err := s.userRepo.GetUserForPassword(ctx, req.Username)
	if err != nil {
		return err
	}
	if Info.PhoneNumber != req.PhoneNumber {
		return error_code.CheckPhoneFail.WithDetail("手机号输入错误。")
	}
	if !hash.CheckPasswordHash(req.OldPassword, Info.PasswordHash) {
		return error_code.OldPasswordFail.WithDetail(error_code.OldPasswordFail.Message)
	}

	if req.OldPassword == req.NewPassword {
		return error_code.PassWordSame.WithDetail(error_code.PassWordSame.Message)
	}
	HashPassword, err := hash.HashPassword(req.NewPassword)
	if err != nil {
		return error_code.ServerError.WithDetail(err.Error())
	}
	err = s.userRepo.UpdatePassword(ctx, Info.UserID, HashPassword)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUserInfoByID(ctx context.Context, userID int64) (*user_model.UserInfo, error) {
	userInfo, err := s.userRepo.GetUserInfoByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (s *UserService) UpdateUserAvatar(ctx context.Context, req *user_model.UpdateAvatarRequest) error {
	return s.userRepo.UpdateUserAvatar(ctx, req.UserID, req.AvatarURL)
}
