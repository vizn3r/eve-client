package com

import (
	"eve-client/serv"
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/fasthttp/websocket"
)

type WSClient struct {
	Msg    chan string
	MsgRes chan string
	serv.Service
}

var (
	WSCLIENT = new(WSClient)
	addr     = flag.CommandLine.String("addr", "localhost:3000", "http service address")
)

func SendWS(msg string) string {
	if WSCLIENT.Status == serv.RUNNING {
		return "WSCLIENT is not running, please chceck WSCLIENT status"
	}
	WSCLIENT.Msg <- msg
	return <-WSCLIENT.MsgRes
}

func ConnectWS() {
	flag.Parse()
	log.SetFlags(0)

	u := url.URL{
		Scheme: "ws",
		Host:   *addr,
		Path:   "/ws/123",
	}
	log.Println("Connecting to:", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("Dial:", err)
		return
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			WSCLIENT.MsgRes <- string(msg)
		}
	}()

	WSCLIENT.Status = serv.RUNNING
	fmt.Println("Connected to server")
	for {
		var msg string
		select {
		case msg = <-WSCLIENT.Msg:
		default:
			msg = ""
		}
		fmt.Println(msg)
		if msg == "exit" {
			if err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
				fmt.Println(err)
				return
			}
			return
		}
		if err := c.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil || msg != "" {
			fmt.Println(err)
			return
		}
	}
}
