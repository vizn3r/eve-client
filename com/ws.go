package com

import (
	"eve-client/inp"
	"eve-client/serv"
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
	WS_HOST  = "10.0.0.111:8080"
)

func ChatWS() {
	for {
		msg := inp.StringInp()
		if strings.EqualFold(msg, "exit") {
			return
		}
		res := SendWS(msg)
		LOG.Message(res)
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
	LOG.Info("Closing WSCLIENT")
	if err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		LOG.Error(err)
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
	LOG.Info("Connecting to:", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	defer CloseWS(c)
	if err != nil {
		LOG.Warning("Dial:", err)
		return
	}

	WSCLIENT.Status = serv.RUNNING

	go func() {
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				LOG.Error(err)
				return
			}
			WSCLIENT.MsgRes <- string(msg)
		}
	}()

	for {
		msg := <-WSCLIENT.Msg

		if msg != "" {
			if err := c.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				LOG.Error(err)
				return
			}
		}
	}
}

func SendController() {
	if !inp.ControllerIsReady() {
		LOG.Warning("Controller is not connected")
		inp.WaitForAny()
	}
	LOG.Info("Use controller to control arm")
	LOG.Info("Exit by pullig LEFT 🕹️ to ➡️  and pressing 🅱️ ")
	for {
		if axis := <-inp.CONTROLLER.Axis; axis[0] > 30000 && <-inp.CONTROLLER.Buttons == 2 {
			return
		}
		msg := "CON" + strconv.Itoa(int(<-inp.CONTROLLER.Buttons)) + "/"
		for i, d := range <-inp.CONTROLLER.Axis {
			msg += strconv.Itoa(d)
			if i != len(<-inp.CONTROLLER.Axis)-1 {
				msg += "/"
			}
		}
		WSCLIENT.Msg <- msg
	}
}
