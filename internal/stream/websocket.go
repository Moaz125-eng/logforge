package stream

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

type TailHandler struct {
	hub    *Hub
	filter Filter
}

func NewTailHandler(hub *Hub) *TailHandler {
	return &TailHandler{hub: hub}
}

func (h *TailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filter := Filter{
		Service: r.URL.Query().Get("service"),
		Level:   r.URL.Query().Get("level"),
		Keyword: r.URL.Query().Get("q"),
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	ch := h.hub.Subscribe(128)
	defer h.hub.Unsubscribe(ch)
	for entry := range ch {
		if !filter.Allow(entry) {
			continue
		}
		data, err := json.Marshal(entry)
		if err != nil {
			continue
		}
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return
		}
	}
}
