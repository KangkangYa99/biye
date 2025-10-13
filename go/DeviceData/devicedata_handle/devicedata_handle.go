package devicedata_handle

import (
	"biye/DeviceData/devicedata_services"
	"biye/model/devicedata_model"
	"biye/share/error_code"
	"biye/share/response"
	websocket "biye/share/webocket"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeviceDataHandle struct {
	deviceDataService *devicedata_services.DeviceDataServices
	hub               *websocket.Hub
}

func NewDeviceDataHandle(
	deviceDataservice *devicedata_services.DeviceDataServices,
	hub *websocket.Hub,
) *DeviceDataHandle {
	return &DeviceDataHandle{
		deviceDataService: deviceDataservice,
		hub:               hub,
	}
}
func (h *DeviceDataHandle) RecStm32Data(c *gin.Context) {
	deviceUID := c.Param("uid")
	if deviceUID == "" {
		c.Error(error_code.NotLogin)
		return
	}
	var req devicedata_model.SendDataReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(error_code.ShouldBindError)
		return
	}
	req.DeviceUID = deviceUID
	resp, err := h.deviceDataService.Insert(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
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
	response.Success(c, gin.H{
		"device_uid": deviceUID,
		"timestamp":  req.DataTimeStamp,
	}, resp.Message)
}
func (d *DeviceDataHandle) GetDeviceHistoryData(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		c.Error(error_code.NotLogin)
		return
	}
	deviceUID := c.Query("device_uid")
	limitStr := c.DefaultQuery("limit", "50")
	if deviceUID == "" || limitStr == "" {
		c.Error(error_code.ShouldBindError)
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"success": false,
			"message": "limit 参数必须是数字",
		})
		return
	}
	req := devicedata_model.DataHistoryRequest{
		DeviceUID: deviceUID,
		Limit:     limit,
	}
	historyData, err := d.deviceDataService.GetDeviceHistoryData(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, historyData, "获取历史数据成功。")

}
