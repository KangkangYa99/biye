package device_handle

import (
	"biye/Devices/device_services"
	"biye/model/device_model"
	"biye/share/error_code"
	"biye/share/response"
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
		c.Error(error_code.NotLogin)
		return
	}
	userIDInt64, _ := userID.(int64)
	var bindReq struct {
		DeviceUID  string `json:"device_uid" binding:"required"`
		DeviceName string `json:"device_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&bindReq); err != nil {
		c.Error(error_code.ShouldBindError)
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
		c.Error(err)
		return
	}
	response.Success(c, gin.H{
		"device_uid": resp.DeviceUID,
	}, resp.Message)
}

func (d *DeviceHandle) Unbind(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.Error(error_code.NotLogin)
		return
	}

	userIDInt64, _ := userID.(int64)
	var unbindReq struct {
		DeviceUID string `json:"device_uid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&unbindReq); err != nil {
		c.Error(error_code.ShouldBindError)
		return
	}
	serviceReq := &device_model.UpdateDeviceUserIDRequest{
		DeviceUID:  unbindReq.DeviceUID,
		DeviceName: nil,
		UserID:     &userIDInt64,
	}
	resp, err := d.device_services.UnBindDevices(c.Request.Context(), serviceReq)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil, resp.Message)
}
func (d *DeviceHandle) GetDevicesByUserID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.Error(error_code.NotLogin)
		return
	}
	userIDInt64, _ := userID.(int64)

	resp, err := d.device_services.GetDevicesByUserID(c.Request.Context(), userIDInt64)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"success": true,
		"message": resp.Message,
		"data":    resp,
	})
}
