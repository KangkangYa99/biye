package routes

import (
	"biye/DeviceData/devicedata_handle"
	"biye/DeviceData/devicedata_services"
	"biye/Devices/device_handle"
	"biye/User/user_handle"
	"biye/share/middleware"
	websocket "biye/share/webocket"

	"github.com/gin-gonic/gin"
)

var (
	userHandler       *user_handle.UserHandle
	deviceHandler     *device_handle.DeviceHandle
	deviceDataHandler *devicedata_handle.DeviceDataHandle
	hub               *websocket.Hub
)

func SetHandle(
	s *devicedata_services.DeviceDataServices,
	u *user_handle.UserHandle,
	d *device_handle.DeviceHandle,
	h *websocket.Hub,
	data *devicedata_handle.DeviceDataHandle,
) {
	deviceDataHandler = data
	userHandler = u
	deviceHandler = d
	hub = h
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"alive": true,
		"from":  "rou.go",
	})
}

/*
	func SetDeviceData(c *gin.Context) {
		deviceUID := c.Param("uid")
		var req devicedata_model.SendDataReportRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "请求数据格式错误",
				"details": err.Error(),
			})
			return
		}
		req.DeviceUID = deviceUID
		resp, err := deviceDataService.SensorData(c.Request.Context(), &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		broadcastData := map[string]interface{}{
			"deviceUID":     req.DeviceUID,
			"temperature":   req.Temperature,
			"humidity":      req.Humidity,
			"light":         req.Light,
			"noise":         req.Noise,
			"fire":          req.Fire,
			"co":            req.CO,
			"light_Status":  req.LightStatus,
			"fan_Status":    req.FanStatus,
			"dataTimeStamp": req.DataTimeStamp,
			"message":       "设备数据已更新",
		}
		hub.Broadcast(deviceUID, broadcastData)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": resp.Message,
			"uid":     deviceUID,
		})
	}
*/
func RegisterRoutes(r *gin.Engine) {
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", userHandler.RegisterUser)
		userGroup.POST("/login", userHandler.LoginUser)
		userGroup.POST("/updatepassword", userHandler.UpdatePassword)
		userGroup.Use(middleware.JWTAuthMiddleware()) // JWT 保护
		{
			userGroup.GET("/info", userHandler.GetUserInfo)
			userGroup.POST("/avatar", userHandler.UploadAvatar)

		}

	}
	deviceGroup := r.Group("/device")
	{
		deviceGroup.Use(middleware.JWTAuthMiddleware())
		{
			deviceGroup.POST("/bind", deviceHandler.Bind)
			deviceGroup.POST("/unbind", deviceHandler.Unbind)
			deviceGroup.POST("/GetDeviceInfo", deviceHandler.GetDevicesByUserID)
		}
	}
	dataGroup := r.Group("/data")
	{
		dataGroup.POST("/:uid/SetDeviceData", deviceDataHandler.RecStm32Data)
	}

	r.GET("/ping", ping)
	r.GET("/ws/device/:uid", websocket.HandleDeviceWS)

}
