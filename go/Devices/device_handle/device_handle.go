package device_handle

import (
	"biye/Devices/device_services"
	"biye/model/device_model"
	"biye/share/error_code"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeviceHandle struct {
	device_services *device_services.DeviceServices
}

func NewDeviceHandle(device_services *device_services.DeviceServices) *DeviceHandle {
	return &DeviceHandle{
		device_services: device_services,
	}
}

func (d *DeviceHandle) Bind(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    error_code.NotLogin.Code,
			"success": false,
			"message": "请先登录",
		})
		return
	}
	userIDInt64, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    error_code.ServerErrorCode,
			"success": false,
			"message": "内部服务器错误",
		})
		return
	}
	var bindReq struct {
		DeviceUID  string `json:"device_uid" binding:"required"`
		DeviceName string `json:"device_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&bindReq); err != nil {
		log.Printf("Bind Device - Invalid JSON format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"success": false,
			"message": "请求数据格式错误",
			"detail":  err.Error(),
		})
		return
	}
	if bindReq.DeviceUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"success": false,
			"message": "设备UID不能为空",
		})
		return
	}

	if bindReq.DeviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"success": false,
			"message": "设备名称不能为空",
		})
		return
	}
	serviceReq := &device_model.UpdateDeviceUserIDRequest{
		DeviceUID:  bindReq.DeviceUID,
		DeviceName: &bindReq.DeviceName,
		UserID:     &userIDInt64,
	}
	resp, err := d.device_services.BindDevices(c.Request.Context(), serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    error_code.ServerErrorCode,
			"success": false,
			"message": error_code.ServerError.Message,
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"success": true,
		"message": resp.Message,
		"data": gin.H{
			"device_uid": resp.DeviceUID,
		},
	})
}
func (d *DeviceHandle) Unbind(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    error_code.NotLogin.Code,
			"success": false,
			"message": "用户未登录",
		})
		return
	}

	userIDInt64, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    error_code.ServerErrorCode,
			"success": false,
			"message": "内部服务器错误",
		})
		return
	}
	var unbindReq struct {
		DeviceUID string `json:"device_uid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&unbindReq); err != nil {
		log.Printf("Bind Device - Invalid JSON format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"success": false,
			"message": "请求数据错误",
		})
		return
	}
	serviceReq := &device_model.UpdateDeviceUserIDRequest{
		DeviceUID:  unbindReq.DeviceUID,
		DeviceName: nil,
		UserID:     &userIDInt64,
	}
	resp, err := d.device_services.UnBindDevices(c.Request.Context(), serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    100,
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"success": true,
		"message": resp.Message,
	})
}
func (d *DeviceHandle) GetDevicesByUserID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    error_code.NotLogin.Code,
			"success": false,
			"message": "请先登录",
		})
		return
	}
	userIDInt64, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    error_code.ServerErrorCode,
			"success": false,
			"message": "内部服务器错误",
		})
		return
	}
	resp, err := d.device_services.GetDevicesByUserID(c.Request.Context(), userIDInt64)
	if err != nil {
		apiErr, ok := err.(*error_code.APIError)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    error_code.ServerErrorCode,
				"success": false,
				"message": "服务器内部错误",
				"detail":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    apiErr.Code,
			"success": false,
			"message": apiErr.Message,
			"detail":  apiErr.Detail,
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
