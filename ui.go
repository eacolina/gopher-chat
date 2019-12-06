package main

import (
	"fmt"
	"log"
	"github.com/marcusolsson/tui-go"
	"net/http"
	"github.com/gorilla/websocket"
	"time"
)


type chatView struct{
	History *tui.Box
	Input *tui.Entry
	Chat *tui.Box

}

func (view *chatView) SetupChat(){
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
	view.Chat = chat
	view.History = history
	view.Input = input
}

func (view *chatView) AppendToHistory(message string){
	view.History.Append(
		tui.NewHBox(
			tui.NewLabel(time.Now().Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", "john"))),
			tui.NewLabel(message),
			tui.NewSpacer(),
		),
	)
}

func connectToSocket(url string) *websocket.Conn{
	header := make(http.Header)
	var Dialer websocket.Dialer

	header.Add("Origin", "http://localhost:3434/")


	conn, resp, err := Dialer.Dial(url, header)

	if err == websocket.ErrBadHandshake {
		fmt.Printf("handshake failed with status %d\n", resp.StatusCode)
		panic(err)
	}
	return conn
}

func (view *chatView) updateHistory(parentUI tui.UI, pipe chan string){
	for {
		message := <-pipe
		view.AppendToHistory(message)
		parentUI.Update(func(){})
	}
}

func checkSocket(conn *websocket.Conn, pipe chan string){
	for {
		_, m, err:= conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		pipe <- string(m)
	}
}


func main() {

	fmt.Println("Starting Client...🚀")
	conn := connectToSocket("ws://localhost:3434/ws")
	chatView := chatView{}
	chatView.SetupChat()

	chatView.Input.OnSubmit(func(e *tui.Entry) {
		chatView.AppendToHistory(e.Text())
		message := []byte(e.Text())
		conn.WriteMessage(websocket.TextMessage, message)
		chatView.Input.SetText("")
	})

	root := tui.NewHBox(chatView.Chat)
	ui, err := tui.New(root)

	pipe := make(chan string)
	go checkSocket(conn, pipe)
	go chatView.updateHistory(ui, pipe)

	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}