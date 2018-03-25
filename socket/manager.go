package socket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Manager struct {
	upgrader websocket.Upgrader
}

func NewManager(in websocket.Upgrader) (obj *Manager) {
	obj = &Manager{}
	obj.upgrader = in
	return
}

func (m *Manager) Echo(w http.ResponseWriter, r *http.Request) {
	c, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
