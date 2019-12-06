package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var port = *flag.String("ip", "3434", "help message for flagname")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	flag.Parse()
	fmt.Println("Starting server...")
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "A Go Server")
		w.Write([]byte("Hello\n"))
	}
	ws_handler := func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}
		for {
			_, p, err := conn.ReadMessage()
			if err != nil{
				break
			}
			fmt.Printf("Received messaged: %s\n", p)
			conn.WriteMessage(websocket.TextMessage, []byte("You sir have been aknowledged"))
		}
	}
	http.HandleFunc("/", handler)
	http.HandleFunc("/ws", ws_handler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil{
		panic(err)
	}
}
