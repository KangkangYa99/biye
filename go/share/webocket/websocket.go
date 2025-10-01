package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var globalHub *Hub

func SetHub(hub *Hub) {
	globalHub = hub
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleDeviceWS(c *gin.Context) {
	uid := c.Param("uid")
	if uid == "" {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "设备UID不能为空",
		})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "升级失败",
		})
		return
	}

	globalHub.Register(uid, conn)
	defer func() {
		globalHub.Unregister(uid, conn) // 👈 注销
		conn.Close()
	}()
	println("✅ 前端 WebSocket 连接已建立，设备 UID:", uid)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			println("❌ 前端断开连接，设备 UID:", uid)
			break
		}
	}
}
