package main

import (
	"biye/DeviceData/devicedata_repository"
	"biye/DeviceData/devicedata_services"
	"biye/Devices/device_handle"
	"biye/Devices/device_repository"
	"biye/Devices/device_services"
	"biye/User/user_handle"
	"biye/User/user_repository"
	"biye/User/user_services"
	"biye/routes"
	"biye/share/logger"
	"biye/share/pgsql"
	websocket "biye/share/webocket"
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("程序开始运行...")
	if err := logger.Init("dev"); err != nil {
		fmt.Println("日志初始化失败:", err)
		panic(err)
	}
	defer logger.Sync()
	fmt.Println("开始加载环境变量...")
	if Enverr := loadEnv(); Enverr != nil {
		fmt.Println("环境变量加载失败:", Enverr)
		return
	}
	fmt.Println("开始初始化数据库...")
	if err := InitDataBase(); err != nil {
		fmt.Println("数据库初始化失败:", err)
		return
	}
	fmt.Println("数据库初始化完成")

	user := user_repository.NewUserRepository(pgsql.DB)
	userService := user_services.NewUserService(user)
	userHandle := user_handle.NewUserHandle(userService)

	deviceRepo := device_repository.NewDeviceRepository(pgsql.DB)
	deviceDataRepo := devicedata_repository.NewDeviceDataRepository(pgsql.DB)

	deviceService := device_services.NewDeviceService(deviceRepo)
	deviceHandle := device_handle.NewDeviceHandle(deviceService)

	deviceDataService := devicedata_services.NewDeviceDataServices(deviceDataRepo, deviceRepo)

	hub := websocket.NewHub()
	go hub.Run()
	websocket.SetHub(hub)
	routes.SetHandle(deviceDataService, userHandle, deviceHandle, hub)

	startServer()
}

func loadEnv() error {
	// 尝试多个可能的路径来加载.env文件
	envPaths := "/home/kang/biye/.env"
	_ = godotenv.Load(envPaths)
	return nil
}

func InitDataBase() error {
	startTime := time.Now()
	if err := pgsql.InitDB(); err != nil {
		logger.GlobalLogger.Error(err.Error(),
			zap.Error(err),
			zap.String("phase", "connection"),
			zap.Duration("duration", time.Since(startTime)),
			zap.String("error_type", "DATABASE_CONNECTION_ERROR"),
			zap.String("host", os.Getenv("DB_HOST")),
			zap.String("database", os.Getenv("DB_NAME")))
		return err
	}
	logger.GlobalLogger.Info("Database connected successfully",
		zap.String("phase", "connection"),
		zap.Duration("duration", time.Since(startTime)),
		zap.String("host", os.Getenv("DB_HOST")),
		zap.String("database", os.Getenv("DB_NAME")),
	)

	return nil
}

// startServer 启动 Gin 服务
func startServer() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
		"http://localhost:4000",
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Length",
		"Content-Type",
		"Authorization",
		"X-Requested-With",
	}
	config.AllowCredentials = true

	r.Use(cors.New(config))
	r.Static("/uploads", "/home/kang/biye/uploads")
	routes.RegisterRoutes(r)
	err := r.Run(":8888")
	if err != nil {
		return
	}
}

/*
func testUserRegistration() {
	ctx, cancle := context.WithTimeout(context.Background(), time.Second*10)
	defer cancle()
	serRepo := user_repository.NewUserRepository(pgsql.DB)
	userService := user_services.NewUserService(serRepo)

	registerReq := &user_model.RegisterRequest{
		Username:    "admin5",
		Password:    "admin",
		PhoneNumber: "555555",
		Email:       "5@qq.com",
		RoleID:      2, // 普通用户
	}
	result, err := userService.RegisterUser(ctx, registerReq)
	if err != nil {
		log.Printf("❌ Registration failed: %v", err)
		return
	}
	log.Printf(" %v", result.UserID)
	log.Printf(" %v", result.Username)
	log.Printf(" %v", result.PhoneNumber)
	log.Printf(" %v", result.Email)
	log.Printf(" %v", result.CreatedAt)
	log.Printf(" %v", result.RoleID)
}
func testbind() {
	ctx, cancle := context.WithTimeout(context.Background(), time.Second*10)
	defer cancle()
	deciveRepo := device_repository.NewDeviceRepository(pgsql.DB)
	deviceService := device_services.NewDeviceService(deciveRepo)

	bind := &device_model.UpdateDeviceUserIDRequest{
		DeviceUID: "device123",
		UserID:    Int64Ptr(52),
	}
	rep, err := deviceService.BindDevices(ctx, bind)
	if err != nil {
		log.Printf("❌ bind failed: %v", err)
		return

	}
	log.Printf("%v", rep.DeviceUID)
	log.Printf("%v", rep.Time)
	log.Printf("%v", rep.Message)
}
func testUnbind() {
	ctx, cancle := context.WithTimeout(context.Background(), time.Second*10)
	defer cancle()
	deciveRepo := device_repository.NewDeviceRepository(pgsql.DB)
	deviceService := device_services.NewDeviceService(deciveRepo)

	bind := &device_model.UpdateDeviceUserIDRequest{
		DeviceUID: "device123",
		UserID:    Int64Ptr(51),
	}
	rep, err := deviceService.UnBindDevices(ctx, bind)
	if err != nil {
		log.Printf("❌ bind failed: %v", err)
		return
	}
	log.Printf("%v", rep.DeviceUID)
	log.Printf("%v", rep.Message)
	log.Printf("%v", rep.Time)
}
func testDeviceInfo() {
	ctx, cancle := context.WithTimeout(context.Background(), time.Second*10)
	defer cancle()
	deciveRepo := device_repository.NewDeviceRepository(pgsql.DB)
	deviceService := device_services.NewDeviceService(deciveRepo)

	bind := &device_model.DeviceInfoRequest{
		DeviceUID: "device123",
		UserID:    Int64Ptr(51),
	}
	rep, err := deviceService.GetDeviceInfo(ctx, bind)
	if err != nil {
		log.Printf("❌ bind failed: %v", err)
		return
	}
	log.Printf("%v", rep.DeviceInfoResponse.DeviceID)
	log.Printf("%v", rep.DeviceInfoResponse.DeviceUID)
	log.Printf("%v", rep.DeviceInfoResponse.DeviceStatus)
	log.Printf("%v", rep.DeviceInfoResponse.LastOnline)
	log.Printf("%v", rep.DeviceInfoResponse.CreatedAt)
	log.Printf("%v", rep.DeviceInfoResponse.UpdatedAt)
	log.Printf("%v", rep.DeviceInfoResponse.IsOnline)
	log.Printf("%v", rep.DeviceInfoResponse.IsBound)

	log.Printf("%v", *rep.OwnerInfo.UserID)
	log.Printf("%v", rep.OwnerInfo.Username)
	log.Printf("%v", rep.OwnerInfo.PhoneNumber)
	log.Printf("%v", rep.OwnerInfo.Email)
	log.Printf("%v", rep.OwnerInfo.CreatedAt)
	log.Printf("%v", rep.OwnerInfo.UpdatedAt)
	log.Printf("%v", rep.OwnerInfo.RoleID)
}
*/
