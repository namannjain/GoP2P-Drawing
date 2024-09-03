package websocket

//his file implements the WebSocket hub for real-time collaboration.
//It manages WebSocket connections, rooms (based on artboard IDs), and message broadcasting.
// It includes functionality to:
// 1) Register and unregister clients
// 2) Broadcast messages to all clients in a specific artboard room
// 3) Handle WebSocket upgrades
// 4) Implement read and write pumps for each client

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Client struct {
	hub        *Hub
	conn       *websocket.Conn
	send       chan []byte
	artboardID string
	userID     string
}

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	rooms      map[string]map[*Client]bool
	mutex      sync.Mutex
}

type Message struct {
	Type       string          `json:"type"`
	ArtboardID string          `json:"artboard_id"`
	UserID     string          `json:"user_id"`
	Data       json.RawMessage `json:"data"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Adjust this for production!
	},
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			if _, ok := h.rooms[client.artboardID]; !ok {
				h.rooms[client.artboardID] = make(map[*Client]bool)
			}
			h.rooms[client.artboardID][client] = true
			h.mutex.Unlock()
		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.rooms[client.artboardID]; ok {
				delete(h.rooms[client.artboardID], client)
				close(client.send)
				if len(h.rooms[client.artboardID]) == 0 {
					delete(h.rooms, client.artboardID)
				}
			}
			h.mutex.Unlock()
		case message := <-h.broadcast:
			var msg Message
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}
			h.mutex.Lock()
			if clients, ok := h.rooms[msg.ArtboardID]; ok {
				for client := range clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(clients, client)
						if len(clients) == 0 {
							delete(h.rooms, msg.ArtboardID)
						}
					}
				}
			}
			h.mutex.Unlock()
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	artboardID := vars["artboardID"]
	userID := r.URL.Query().Get("user_id") // Assume user_id is passed as a query parameter

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), artboardID: artboardID, userID: userID}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		msg.UserID = c.userID
		msg.ArtboardID = c.artboardID

		updatedMessage, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshaling message: %v", err)
			continue
		}

		c.hub.broadcast <- updatedMessage
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}
