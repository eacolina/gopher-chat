package main

import (
	"fmt"
	"flag"
	"github.com/marcusolsson/tui-go"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"time"
)


var to_user string
var from_user string

type chat struct{
	View *chatView
	Conn *websocket.Conn
}

type chatView struct{
	History *tui.Box
	Input *tui.Entry
	Chat *tui.Box

}
type message struct {
	Sender string
	Receiver string
	Content string

}


func (chat *chat) initChat(socket_url string){
	chat.connectToSocket(socket_url)
	view := chatView{}
	view.SetupChatView()
	chat.View = &view
	chat.View.Input.OnSubmit(func(e *tui.Entry) {
		msg := message{from_user, to_user, e.Text()}
		chat.View.AppendToHistory(msg)
		err := chat.Conn.WriteJSON(msg)
		if err != nil {
			chat.View.Input.SetText(err.Error())
		} else {
			chat.View.Input.SetText("")
		}

	})

}

func (chat *chat) connectToSocket(url string) {
	header := make(http.Header)
	var Dialer websocket.Dialer

	header.Add("Origin", "http://localhost:3434/")
	header.Add("sender", from_user)


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

func (view *chatView) AppendToHistory(msg message){
	view.History.Append(
		tui.NewHBox(
			tui.NewLabel(time.Now().Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", msg.Sender))),
			tui.NewLabel(msg.Content),
			tui.NewSpacer(),
		),
	)
}

func (view *chatView) updateHistory(parentUI tui.UI, pipe chan message){
	for {
		message := <-pipe
		view.AppendToHistory(message)
		parentUI.Update(func(){})
	}
}

func checkSocket(conn *websocket.Conn, pipe chan message){
	for {
		message := message{}
		err:= conn.ReadJSON(&message)
		if err != nil {
			fmt.Println(err)
			return
		}
		pipe <- message
	}
}



func main() {
	fmt.Println("Starting Client...🚀")
	from := flag.String("from", "", "your user name")
	to := flag.String("to", "","their user name")
	flag.Parse()
	from_user = *from
	to_user = *to
	chat := chat{}
	chat.initChat("ws://localhost:3434/ws")

	root := tui.NewHBox(chat.View.Chat)
	ui, err := tui.New(root)

	pipe := make(chan message)
	go checkSocket(chat.Conn, pipe)
	go chat.View.updateHistory(ui, pipe)

	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() {
		chat.Conn.Close()
		ui.Quit()
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}