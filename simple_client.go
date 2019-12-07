package main

import (
	"fmt"
	"github.com/gorilla/websocket"

	"net/http"
)

func main(){
	fmt.Println("Starting Client...🚀")
	header := make(http.Header)
	var Dialer websocket.Dialer

	header.Add("Origin", "http://localhost:3131/")


	conn, resp, err := Dialer.Dial("ws://localhost:3434/ws", header)

	if err == websocket.ErrBadHandshake {
		fmt.Printf("handshake failed with status %d\n", resp.StatusCode)
		panic(err)
	}
	for i:= 0; i < 10; i++{
		conn.WriteMessage(websocket.TextMessage, []byte("This is my reply to you!🤙"))

	}
}