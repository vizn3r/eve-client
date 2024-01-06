package com

import (
	"eve-client/inp"
	"eve-client/serv"
	"flag"
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

var (
	WSCLIENT = new(WSClient)
	addr     = flag.CommandLine.String("addr", "localhost:3000", "http service address")
)

func ChatWS() {
	for {
		msg := inp.StringInp()
		fmt.Println("Test:", msg)
		if strings.EqualFold(msg, "exit") {
			return
		}
		res := SendWS(msg)
		fmt.Println(res)
	}
}

func SendWS(msg string) string {
	if WSCLIENT.Status != serv.RUNNING {
		return "WSCLIENT is not running, please chceck WSCLIENT status"
	}
	WSCLIENT.Msg <- msg
	return <-WSCLIENT.MsgRes
}

func CloseWS(c *websocket.Conn) {
	if err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		fmt.Println(err)
		return
	}
	c.Close()
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
	defer CloseWS(c)
	if err != nil {
		log.Println("Dial:", err)
		return
	}

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
	for {
		msg := <-WSCLIENT.Msg
		// if msg == "exit" {
		// 	CloseWS(c)
		// 	return
		// }
		if msg != "" {
			if err := c.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
