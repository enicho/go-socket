package main

import (
	"log"
	"net/http"

	"github.com/enicho/go-socket/socket"
	"github.com/gorilla/websocket"
)

func main() {
	sockManage := socket.NewManager(websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		}})
	http.HandleFunc("/echo", sockManage.Echo)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
