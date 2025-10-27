package user_services

import (
	"biye/User/user_repository"
	"biye/model/user_model"
	"biye/share/error_code"
	"biye/share/hash"
	"biye/share/jwt"
	"biye/share/redis"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
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
		return nil, error_code.UserExists
	}
	if phoneExists {
		return nil, error_code.UserNumberExists
	}
	if emailExists {
		return nil, error_code.UserEmailExists
	}
	passwordHash, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, error_code.ServerError
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
		return nil, error_code.PasswordFail
	}
	token, err := jwt.GenerateToken(user.UserID)
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
		return error_code.CheckPhoneFail
	}
	if !hash.CheckPasswordHash(req.OldPassword, Info.PasswordHash) {
		return error_code.OldPasswordFail
	}

	if req.OldPassword == req.NewPassword {
		return error_code.PassWordSame
	}
	HashPassword, err := hash.HashPassword(req.NewPassword)
	if err != nil {
		return error_code.ServerError
	}
	err = s.userRepo.UpdatePassword(ctx, Info.UserID, HashPassword)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUserInfoByID(ctx context.Context, userID int64) (*user_model.UserInfo, error) {
	cachekey := fmt.Sprintf("user_info:%d", userID)
	val, err := redis.RedisClient.Get(ctx, cachekey).Result()
	if err == nil {
		var cacheUser user_model.UserInfo
		if json.Unmarshal([]byte(val), &cacheUser) == nil {
			return &cacheUser, nil
		}
	}
	userInfo, err := s.userRepo.GetUserInfoByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	userBytes, _ := json.Marshal(userInfo)
	redis.RedisClient.Set(ctx, cachekey, userBytes, 10*time.Minute)
	return userInfo, nil
}

func (s *UserService) UpdateUserAvatar(ctx context.Context, req *user_model.UpdateAvatarRequest) error {
	err := s.userRepo.UpdateUserAvatar(ctx, req.UserID, req.AvatarURL)
	if err != nil {
		return err
	}
	cacheKey := fmt.Sprintf("user_info:%d", req.UserID)
	err = redis.RedisClient.Del(ctx, cacheKey).Err()
	if err != nil {
		log.Printf("删除信息缓存失败:%v", err)
	} else {
		log.Printf("缓存已删除:%s", cacheKey)
	}
	return nil
}
