package devicedata_handle

import (
	"biye/DeviceData/devicedata_services"
	"biye/model/devicedata_model"
	websocket "biye/share/webocket"
	"net/http"
	"strconv"

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
	req.DeviceUID = deviceUID
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
func (d *DeviceDataHandle) GetDeviceHistoryData(c *gin.Context) {
	// 从 context 获取用户 ID（token 已在中间件验证）
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"success": false,
			"message": "未授权访问，缺少用户信息",
		})
		return
	}

	// 从 Query 读取参数
	deviceUID := c.Query("device_uid")
	limitStr := c.DefaultQuery("limit", "50")

	// 参数校验
	if deviceUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"success": false,
			"message": "缺少 device_uid 参数",
		})
		return
	}

	// 转换 limit
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"success": false,
			"message": "limit 参数必须是数字",
		})
		return
	}

	// 构造请求结构体
	req := devicedata_model.DataHistoryRequest{
		DeviceUID: deviceUID,
		Limit:     limit,
	}

	// 调用 service 层逻辑
	historyData, err := d.deviceDataservice.GetDeviceHistoryData(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"success": false,
			"detail":  err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"success": true,
		"message": "获取历史数据成功",
		"user_id": userID, // 可选，调试时可返回
		"data":    historyData,
	})
}
