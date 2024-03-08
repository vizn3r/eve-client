package com

import (
	"bufio"
	"eve-client/inp"
	"eve-client/serv"
	"fmt"
	"log"
	"net/url"
	"strconv"
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
	WS_HOST  = "localhost:8080"
)

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
	if !WSCLIENT.IsRunning() {
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
		Host:   WS_HOST,
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

func MotorController() {
	numMotors := -1
	rawOut := SendWS("m0")
	s := bufio.NewScanner(strings.NewReader(rawOut))
	for s.Scan() {
		numMotors++
	}
	for {
		in := inp.Inp()
		if in == inp.Back {
			return
		}
		for i, con := range <-inp.CONTROLLER.Axis {
			dir := 0
			if con < 0 {
				dir = 1
			}
			// m3 <>
			SendWS("m3 " + strconv.Itoa(i) + " " + strconv.Itoa(dir) + " 5 10")
		}
	}
}
