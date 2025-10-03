package devicedata_handle

import (
	"biye/DeviceData/devicedata_services"
	"biye/model/devicedata_model"
	websocket "biye/share/webocket"

	"net/http"

	"github.com/gin-gonic/gin"
)

type DeviceDataHandle struct {
	deviceDataservice *devicedata_services.DeviceDataServices
	hub               *websocket.Hub
}

func NewDeviceDataHandle(
	deviceDataservice *devicedata_services.DeviceDataServices,
	hub *websocket.Hub,
) *DeviceDataHandle {
	return &DeviceDataHandle{
		deviceDataservice: deviceDataservice,
		hub:               hub,
	}
}
func (h *DeviceDataHandle) RecStm32Data(c *gin.Context) {
	deviceUID := c.Param("uid")
	if deviceUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"success": false,
			"message": "UID解析失败。",
		})
		return
	}
	var req devicedata_model.SendDataReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"success": false,
			"message": "请求数据格式错误",
			"detail":  err.Error(),
		})
		return
	}
	req.DeviceUID = c.Param("uid")
	resp, err := h.deviceDataservice.Insert(c.Request.Context(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"success": false,
			"message": "数据处理失败",
			"detail":  err.Error(),
		})
		return
	}
	broadcastData := map[string]interface{}{
		"type":         "sensor_data",
		"deviceUID":    req.DeviceUID,
		"temperature":  req.Temperature,
		"humidity":     req.Humidity,
		"light":        req.Light,
		"noise":        req.Noise,
		"fire":         req.Fire,
		"co":           req.CO,
		"light_status": req.LightStatus,
		"fan_status":   req.FanStatus,
		"timestamp":    req.DataTimeStamp,
		"message":      "设备数据已更新",
	}
	h.hub.Broadcast(deviceUID, broadcastData)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"success": true,
		"message": resp.Message,
		"data": gin.H{
			"device_uid": deviceUID,
			"timestamp":  req.DataTimeStamp,
		},
	})
}
