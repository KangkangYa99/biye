package websocket

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	UID  string          //设备uid
	Conn *websocket.Conn //Websocket
}
type Hub struct {
	clients    map[string]map[*Client]bool
	register   chan *Client          //注册通道
	unregister chan *Client          //注销
	broadcast  chan BroadcastMessage //广播
	mutex      sync.RWMutex          //锁

}
type BroadcastMessage struct {
	UID  string
	Data interface{}
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan BroadcastMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unRegisterClient(client)
		case msg := <-h.broadcast:
			h.broadcastMessage(msg)
		}
	}
}
func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.clients[client.UID] == nil {

		h.clients[client.UID] = make(map[*Client]bool)
	}
	h.clients[client.UID][client] = true
}
func (h *Hub) unRegisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if clients, ok := h.clients[client.UID]; ok {
		delete(clients, client)
	}
}
func (h *Hub) broadcastMessage(msg BroadcastMessage) {
	h.mutex.RLock()
	clients, ok := h.clients[msg.UID]
	h.mutex.RUnlock()
	if !ok {
		return
	}
	data, err := json.Marshal(msg.Data)
	if err != nil {
		return
	}
	for client := range clients {
		client.Conn.WriteMessage(websocket.TextMessage, data)
	}
}
func (h *Hub) Register(uid string, conn *websocket.Conn) {
	h.register <- &Client{UID: uid, Conn: conn}
}
func (h *Hub) Unregister(uid string, conn *websocket.Conn) {
	h.unregister <- &Client{UID: uid, Conn: conn}
}
func (h *Hub) Broadcast(uid string, data interface{}) {
	h.broadcast <- BroadcastMessage{UID: uid, Data: data}
}
