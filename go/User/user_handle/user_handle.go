package user_handle

import (
	"biye/User/user_services"
	"biye/model/user_model"
	"biye/share/error_code"
	"biye/share/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"

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
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"error":   "请求数据格式错误",
			"message": err.Error(),
		})
		return
	}
	resp, err := h.userServices.RegisterUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"success": true,
		"message": resp.Message,
		"data":    resp,
	})

}
func (h *UserHandle) LoginUser(c *gin.Context) {
	var req user_model.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"success": false,
			"message": err.Error(),
		})
		log.Printf(err.Error())
		return
	}
	resp, err := h.userServices.LoginUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"success": false,
			"message": err.Error(),
		})
		log.Printf(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"success": true,
		"message": resp.Message,
		"data": gin.H{
			"token": resp.Token,
		},
	})
}

func (h *UserHandle) UpdatePassword(c *gin.Context) {
	var req user_model.UpdatePasswordRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"success": false,
			"message": err.Error(),
		})
		return
	}
	err := h.userServices.UpdatePassword(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"success": true,
		"message": "修改密码成功。",
	})
}

func (h *UserHandle) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    error_code.NotLogin.Code,
			"success": false,
			"message": "请先登录",
		})
		return
	}
	userIDInt64 := userID.(int64)
	userInfo, err := h.userServices.GetUserInfoByID(c, userIDInt64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    error_code.DatabaseError.Code,
			"message": "获取用户信息失败",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data":    userInfo,
	})
}

func (h *UserHandle) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"success": false,
			"message": "请上传头像",
			"detail":  err.Error(),
		})
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
			"success": false,
			"code":    400,
			"messgae": "文件大小不能超过5MB",
		})
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    error_code.NotLogin.Code,
			"message": "请先登录",
		})
		return
	}
	filename := utils.GenerateUniqueFilename(file.Filename)
	uploadDir := "/home/kang/biye/uploads/avatars"
	filepath := filepath.Join(uploadDir, filename)
	err = os.MkdirAll(uploadDir, 0755)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建上传目录失败",
			"detail":  err.Error(),
		})
		return
	}
	err = c.SaveUploadedFile(file, filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "文件保存失败",
			"detail":  err.Error(),
		})
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新头像失败",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "头像上传成功",
		"data": gin.H{
			"avatar_url": avatarURL,
			"filename":   filename,
		},
	})
}
