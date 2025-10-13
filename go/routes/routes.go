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

func RegisterRoutes(r *gin.Engine) {
	r.Use(middleware.ErrorHandler())
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
		dataGroup.Use(middleware.JWTAuthMiddleware())
		{
			dataGroup.GET("/GetDataHistory", deviceDataHandler.GetDeviceHistoryData)
		}
	}
	r.GET("/ws/device/:uid", websocket.HandleDeviceWS)

	r.GET("/ping", ping)
}
