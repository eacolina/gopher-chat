package main

import (
	"fmt"
	"log"
	"github.com/marcusolsson/tui-go"
	"net/http"
	"github.com/gorilla/websocket"
	"time"
)


type chat struct{
	View *chatView
	Conn *websocket.Conn
}

type chatView struct{
	History *tui.Box
	Input *tui.Entry
	Chat *tui.Box

}

func (chat *chat) initChat(socket_url string){
	chat.connectToSocket(socket_url)
	view := chatView{}
	view.SetupChatView()
	chat.View = &view
	chat.View.Input.OnSubmit(func(e *tui.Entry) {
		chat.View.AppendToHistory(e.Text())
		message := []byte(e.Text())
		chat.Conn.WriteMessage(websocket.TextMessage, message)
		chat.View.Input.SetText("")
	})

}

func (chat *chat) connectToSocket(url string) {
	header := make(http.Header)
	var Dialer websocket.Dialer

	header.Add("Origin", "http://localhost:3434/")


	conn, resp, err := Dialer.Dial(url, header)

	if err == websocket.ErrBadHandshake {
		fmt.Printf("handshake failed with status %d\n", resp.StatusCode)
		panic(err)
	}
	chat.Conn = conn
}

func (view *chatView) SetupChatView(){
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
	chat := chat{}
	chat.initChat("ws://localhost:3434/ws")

	root := tui.NewHBox(chat.View.Chat)
	ui, err := tui.New(root)

	pipe := make(chan string)
	go checkSocket(chat.Conn, pipe)
	go chat.View.updateHistory(ui, pipe)

	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}