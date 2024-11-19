package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	EnableCompression: true,
}

type WebsocketClient struct {
	conn     *websocket.Conn
	deviceId string
	connType string
	mu       sync.RWMutex
	done     chan struct{}
}

func newClient(conn *websocket.Conn, deviceId string, connType string) *WebsocketClient {
	conn.EnableWriteCompression(true)
	client := &WebsocketClient{
		conn:     conn,
		deviceId: deviceId,
		connType: connType,
		done:     make(chan struct{}),
	}

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		log.Printf("Received pong from client: deviceId=%s", client.deviceId)
		return nil
	})

	return client
}

func (c *WebsocketClient) readPump(wsServer *WsServer) {
	defer func() {
		wsServer.unregister <- c
	}()

	for {
		select {
		case <-c.done:
			return
		default:
			_, _, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				return
			}
			// Handle other message types if needed
		}
	}
}

func (c *WebsocketClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Failed to send ping to client %s: %v", c.deviceId, err)
				return
			}
			log.Printf("Sent ping to client: deviceId=%s", c.deviceId)

		case <-c.done:
			return
		}
	}
}

type WsServer struct {
	clients    map[*WebsocketClient]bool
	register   chan *WebsocketClient
	unregister chan *WebsocketClient
	mu         sync.RWMutex
	done       chan struct{} // Channel to signal shutdown
}

func NewWsServer() *WsServer {
	return &WsServer{
		clients:    make(map[*WebsocketClient]bool),
		register:   make(chan *WebsocketClient),
		unregister: make(chan *WebsocketClient),
		done:       make(chan struct{}),
	}
}

func ValidateUser(token string) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, "http://server:5000/api/validate/admin", nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	log.Println("statusCode", res.StatusCode)
	return res.StatusCode == 200, nil
}

func ServeWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Sec-Websocket-Protocol")
	if token == "" {
		log.Println("Missing authentication token")
		http.Error(w, "Missing authentication token", http.StatusUnauthorized)
		return
	}

	valid, err := ValidateUser(token)
	if err != nil {
		log.Printf("Error validating user: %v", err)
		http.Error(w, "Authentication error", http.StatusInternalServerError)
		return
	}
	if !valid {
		log.Println("User not valid, refusing connection!")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	deviceId := r.URL.Query().Get("deviceId")
	connType := r.URL.Query().Get("connType")
	if deviceId == "" || connType == "" {
		log.Println("Invalid request! Missing device id or connection params!")
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	client := newClient(conn, deviceId, connType)
	wsServer.register <- client

	go client.writePump()
	go client.readPump(wsServer)

	log.Printf("New client connected: deviceId=%s, connType=%s", deviceId, connType)
}

func (server *WsServer) Run() {
	// Start the health check goroutine

	for {
		select {
		case client := <-server.register:
			server.registerClient(client)

		case client := <-server.unregister:
			server.unregisterClient(client)

		case <-server.done:
			return
		}
	}
}

func (server *WsServer) SendMessage(deviceId string, message []byte, connType string) {
	server.mu.RLock()
	defer server.mu.RUnlock()

	for client := range server.clients {
		if client.deviceId == deviceId && client.connType == connType {
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := client.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Error sending message to client %s: %v", client.deviceId, err)
				go server.unregisterClient(client)
			}
		}
	}
}

func (server *WsServer) Shutdown() {
	close(server.done)

	server.mu.Lock()
	defer server.mu.Unlock()

	for client := range server.clients {
		client.conn.Close()
	}
	server.clients = make(map[*WebsocketClient]bool)
}

func (server *WsServer) registerClient(client *WebsocketClient) {
	server.mu.Lock()
	defer server.mu.Unlock()

	server.clients[client] = true
	log.Printf("Client registered: deviceId=%s, connType=%s", client.deviceId, client.connType)
}

func (server *WsServer) unregisterClient(client *WebsocketClient) {
	server.mu.Lock()
	defer server.mu.Unlock()

	if _, ok := server.clients[client]; ok {
		close(client.done)
		delete(server.clients, client)
		client.conn.Close()
		log.Printf("Client unregistered: deviceId=%s, connType=%s", client.deviceId, client.connType)
	}
}

func HttpServer(port string, wsServer *WsServer) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(wsServer, w, r)
	})

	server := &http.Server{
		Addr:         port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting http server: %v", err)
	}
}
