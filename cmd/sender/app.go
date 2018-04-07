package main

import (
	"bufio"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type (
	SendData struct {
		clientID string
		message  string
	}

	ReceiveData struct {
		clientID string
		message  string
	}
)

func main() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	tmp := "ws://localhost:8080/ws?name=elliot"

	// u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	u, _ := url.Parse(tmp)
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	sendMsg := make(chan string)
	go chat(sendMsg)

	for {
		select {
		case <-done:
			return
		case t := <-sendMsg:
			err := c.WriteMessage(websocket.TextMessage, []byte(t))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func chat(c chan string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		c <- text
	}
}
