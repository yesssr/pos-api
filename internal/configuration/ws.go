package configuration

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients map[string]*websocket.Conn
	mu      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*websocket.Conn),
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Hub) WSHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("userId")
    if userID == "" {
        http.Error(w, "userId required", http.StatusBadRequest)
        return
    }

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, "Failed to upgrade websocket", http.StatusInternalServerError)
        return
    }

    h.mu.Lock()
    h.clients[userID] = conn
    h.mu.Unlock()

    go func() {
        defer func() {
            h.mu.Lock()
            delete(h.clients, userID)
            h.mu.Unlock()
            conn.Close()
        }()
        for {
            if _, _, err := conn.NextReader(); err != nil {
                break
            }
        }
    }()
}


func (h *Hub) NotifyUser(userID string, message any) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	conn, ok := h.clients[userID]
	if !ok {
		return nil;
	}

	return conn.WriteJSON(message)
}
