package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)


type message struct {
	Sender string
	Receiver string
	Content string

}

var port = *flag.String("ip", "3434", "help message for flagname")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var connections = make(map [string]*websocket.Conn)

func listenForMessages(conn *websocket.Conn, sender string){
	for{
		message := message{}
		err := conn.ReadJSON(&message)
		if err != nil{
			panic(err)
		}
		fmt.Printf("Received messaged: `%s` from %s to %s\n", message.Content, sender, message.Receiver)
		message.Sender = sender
		connections[message.Receiver].WriteJSON(message)
	}
}

func main() {
	flag.Parse()

	fmt.Println("Starting server... 🚀")
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "A Go Server")
		w.Write([]byte("Hello\n"))
	}
	ws_handler := func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		sender := r.Header.Get("sender")

		conn, err := upgrader.Upgrade(w, r, nil)
		connections[sender] = conn
		go listenForMessages(conn, sender)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Done with connection for %s", sender)
	}
	http.HandleFunc("/", handler)
	http.HandleFunc("/ws", ws_handler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil{
		panic(err)
	}
}
