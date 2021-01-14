package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("start home page"))
	fmt.Println("[info] home page")
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("[error]", err)
			return
		}
		log.Println("[info]", string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println("[error]", err)
			return
		}
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[error] %v\n", err)
	}
	log.Println("[info] client connected")

	if err = ws.WriteMessage(1, []byte("HI Client")); err != nil {
		log.Printf("[error] %v\n", err)
	}
	reader(ws)

}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	log.Print("start")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
