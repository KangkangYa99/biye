package user_handle

import (
	"biye/User/user_services"
	"biye/model/user_model"
	"biye/share/error_code"
	"biye/share/jwt"
	"biye/share/redis"
	"biye/share/response"
	"biye/share/utils"
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandle struct {
	userServices *user_services.UserService
}

func NewUserHandle(userServices *user_services.UserService) *UserHandle {
	return &UserHandle{
		userServices: userServices,
	}
}
func (h *UserHandle) RegisterUser(c *gin.Context) {
	var req user_model.RegisterRequest

	if err := c.ShouldBind(&req); err != nil {
		c.Error(error_code.ShouldBindError)
		return
	}
	resp, err := h.userServices.RegisterUser(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil, resp.Message)
}
func (h *UserHandle) LoginUser(c *gin.Context) {
	var req user_model.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.Error(error_code.ShouldBindError)
		return
	}
	resp, err := h.userServices.LoginUser(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, gin.H{
		"token": resp.Token,
	}, resp.Message)
}
func (h *UserHandle) LoginOut(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.Error(error_code.NotLogin)
		return
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.Error(error_code.InvalidToken)
		return
	}
	tokenString := parts[1]
	claims, err := jwt.ParseToken(parts[1])
	if err != nil {
		c.Error(error_code.InvalidToken)
		return
	}
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl > 0 {
		key := "jwt_blacklist:" + tokenString
		err = redis.RedisClient.Set(context.Background(), key, "blacklist", ttl).Err()
		if err != nil {
			c.Error(err)
		}
	}
	response.Success(c, nil, "登出成功。")
}
func (h *UserHandle) UpdatePassword(c *gin.Context) {
	var req user_model.UpdatePasswordRequest
	if err := c.ShouldBind(&req); err != nil {
		c.Error(err)
		return
	}
	err := h.userServices.UpdatePassword(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil, "修改密码成功。")
}

func (h *UserHandle) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.Error(error_code.NotLogin)
		return
	}
	userIDInt64 := userID.(int64)
	userInfo, err := h.userServices.GetUserInfoByID(c, userIDInt64)

	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, userInfo, "获取成功。")
}

func (h *UserHandle) UploadAvatar(c *gin.Context) {

	userID, exists := c.Get("userID")
	if !exists {
		c.Error(error_code.NotLogin)
		return
	}
	file, err := c.FormFile("avatar")
	if err != nil {
		c.Error(err)
		return
	}
	if !utils.IsVaildImageFile(file.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"success": false,
			"message": "只支持 JPG,PNG,GIF格式的图片",
		})
		return
	}
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"success": false,
			"message": "文件大小不能超过5MB",
		})
		return
	}
	filename := utils.GenerateUniqueFilename(file.Filename)
	uploadDir := "/home/kang/biye/uploads/avatars"
	filepath := filepath.Join(uploadDir, filename)
	err = os.MkdirAll(uploadDir, 0755)
	if err != nil {
		c.Error(err)
		return
	}
	err = c.SaveUploadedFile(file, filepath)
	if err != nil {
		c.Error(err)
		return
	}
	avatarURL := "/uploads/avatars/" + filename
	req := &user_model.UpdateAvatarRequest{
		UserID:    userID.(int64),
		AvatarURL: avatarURL,
	}
	err = h.userServices.UpdateUserAvatar(c, req)
	if err != nil {
		os.Remove(filepath)
		c.Error(err)
		return
	}
	response.Success(c, gin.H{
		"avatar_url": avatarURL,
		"filename":   filename,
	}, "头像上传成功。")
}
