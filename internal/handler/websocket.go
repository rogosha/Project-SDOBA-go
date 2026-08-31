package handler

import (
	"strconv"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type WSHandler struct {
	clients map[uint]map[*websocket.Conn]bool
	mu      sync.RWMutex
}

func NewWSHandler() *WSHandler {
	return &WSHandler{
		clients: make(map[uint]map[*websocket.Conn]bool),
	}
}

func (h *WSHandler) Handle(c *websocket.Conn) {
	conversationID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return
	}

	id := uint(conversationID)

	h.mu.Lock()

	if h.clients[id] == nil {
		h.clients[id] = make(map[*websocket.Conn]bool)
	}

	h.clients[id][c] = true

	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients[id], c)

		if len(h.clients[id]) == 0 {
			delete(h.clients, id)
		}

		h.mu.Unlock()

		c.Close()
	}()

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			break
		}

		h.broadcast(id, messageType, message)
	}
}

func (h *WSHandler) broadcast(conversationID uint, messageType int, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients[conversationID] {
		if err := client.WriteMessage(messageType, message); err != nil {
			continue
		}
	}
}
