package com

import (
	"eve-client/inp"
	"eve-client/serv"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/fasthttp/websocket"
)

type WSClient struct {
	Msg    chan string
	MsgRes chan string
	serv.Service
}

var WSCLIENT = new(WSClient)

func ChatWS() {
	for {
		msg := inp.StringInp()
		if strings.EqualFold(msg, "exit") {
			return
		}
		res := SendWS(msg)
		fmt.Println(res)
	}
}

func SendWS(msg string) string {
	if WSCLIENT.IsRunning() {
		return "WSCLIENT is not running, please chceck WSCLIENT status"
	}
	WSCLIENT.Msg <- msg
	return <-WSCLIENT.MsgRes
}

func CloseWS(c *websocket.Conn) {
	fmt.Println("Closing WSCLIENT")
	if err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		fmt.Println(err)
		return
	}
	WSCLIENT.Status = serv.STOPPED
	c.Close()
}

func ConnectWS() {
	u := url.URL{
		Scheme: "ws",
		Host:   "localhost:3000",
		Path:   "/ws/123",
	}
	fmt.Println("Connecting to:", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	defer CloseWS(c)
	if err != nil {
		log.Println("Dial:", err)
		return
	}

	WSCLIENT.Status = serv.RUNNING

	go func() {
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			WSCLIENT.MsgRes <- string(msg)
		}
	}()

	for {
		msg := <-WSCLIENT.Msg

		if msg != "" {
			if err := c.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
