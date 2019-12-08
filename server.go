package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"sync"

)


type message struct {
	Sender string
	Receiver string
	Content string
}

type Hub struct {
	Connections map[string] *websocket.Conn
	ConnectionsMux sync.Mutex
	Upgrader websocket.Upgrader
}

func (hub *Hub) InitHub(){
	hub.Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	hub.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	hub.Connections = make(map[string] *websocket.Conn)
	hub.ConnectionsMux = sync.Mutex{}
}

var port = *flag.String("ip", "3434", "help message for flagname")


var connections = make(map [string]*websocket.Conn)

func (hub *Hub) listenForMessages(conn *websocket.Conn, sender string){
	for{
		message := message{}
		err := conn.ReadJSON(&message)
		if err != nil {
			fmt.Printf("Connection closed for %s\n", sender)
			return
		}
		fmt.Printf("Received messaged: `%s` from %s to %s\n", message.Content, sender, message.Receiver)
		message.Sender = sender
		hub.ConnectionsMux.Lock()
			out_conn := hub.Connections[message.Receiver]
		hub.ConnectionsMux.Unlock()
		if out_conn == nil{
			fmt.Println("Receiver not found!")
			continue
		}
		out_conn.WriteJSON(message)
	}
}

func main() {
	flag.Parse()
	fmt.Println("Starting server... 🚀")
	hub := Hub{}
	hub.InitHub()


	ws_handler := func(w http.ResponseWriter, r *http.Request) {
		sender := r.Header.Get("sender")

		conn, err := hub.Upgrader.Upgrade(w, r, nil)
		conn.SetCloseHandler(func(code int, text string) error {
			fmt.Println("Connection Closed")
			return nil
		})

		hub.ConnectionsMux.Lock()
			hub.Connections[sender] =  conn
		hub.ConnectionsMux.Unlock()

		go hub.listenForMessages(conn, sender)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s started connection\n", sender)
	}
	http.HandleFunc("/ws", ws_handler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil{
		panic(err)
	}
}
