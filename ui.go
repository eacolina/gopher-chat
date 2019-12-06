package main

import (
	"fmt"
	"log"
	"time"

	"github.com/marcusolsson/tui-go"
	"net/http"
	"github.com/gorilla/websocket"
)

type post struct {
	username string
	message  string
	time     string
}

var posts = []post{
	{username: "john", message: "hi, what's up?", time: "14:41"},
	{username: "jane", message: "not much", time: "14:43"},
}

func main() {

	fmt.Println("Starting Client...🚀")
	header := make(http.Header)
	var Dialer websocket.Dialer

	header.Add("Origin", "http://localhost:3131/")


	conn, resp, err := Dialer.Dial("ws://localhost:3434/ws", header)

	if err == websocket.ErrBadHandshake {
		fmt.Printf("handshake failed with status %d\n", resp.StatusCode)
		panic(err)
	}

	history := tui.NewVBox()

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	input.OnSubmit(func(e *tui.Entry) {
		message := []byte(e.Text())
		conn.WriteMessage(websocket.TextMessage, message)
		history.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", "john"))),
			tui.NewLabel(e.Text()),
			tui.NewSpacer(),
		))
		input.SetText("")
	})

	root := tui.NewHBox(chat)
	ui, err := tui.New(root)
	for i := 0; i < 10; i++{
		history.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", "john"))),
			tui.NewLabel("message"),
			tui.NewSpacer(),
		))
	}
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}