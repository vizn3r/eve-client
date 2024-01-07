package main

import (
	"eve-client/cli"
	"eve-client/com"
	"eve-client/inp"
	"eve-client/serv"
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	inp.CONTROLLER.Axis = make(chan []int)
	inp.CONTROLLER.Buttons = make(chan uint32)
	inp.CONTROLLER.Exit = make(chan bool)
	inp.CONTROLLER.Status = serv.STARTING
	inp.CONTROLLER.WG = &wg

	inp.KEYBOARD.Output = make(chan inp.KeyboardOutput)
	inp.KEYBOARD.Exit = make(chan bool)
	inp.KEYBOARD.GetKey = make(chan bool, 1)
	inp.KEYBOARD.Status = serv.STARTING
	inp.KEYBOARD.WG = &wg

	com.WSCLIENT.Msg = make(chan string)
	com.WSCLIENT.MsgRes = make(chan string)
	com.WSCLIENT.Status = serv.STARTING

	go inp.OpenKeyboard()
	inp.KEYBOARD.GetKey <- true
	go inp.OpenController(0)
	go com.ConnectWS()

	exitM := cli.Menu{
		Header: "Do you really want to exit?",
		Opts: []cli.Opt{
			{
				Name: "Yes",
				Func: func() {
					fmt.Println("Exited")
					os.Exit(0)
				},
			},
		},
	}

	wsmenu := cli.Menu{
		Header: "WebSocket communication",
		Opts: []cli.Opt{
			{
				Name: "Send Message",
				Func: func() {
					msg := com.SendWS(inp.StringInp())
					fmt.Println(msg)
					time.Sleep(time.Second)
				},
			},
			{
				Name: "Chat",
				Func: com.ChatWS,
			},
		},
	}

	menu := cli.Menu{
		Header: `
▄▄▄ . ▌ ▐·▄▄▄ .     ▄▄· ▄▄▌  ▪  ▄▄▄ . ▐ ▄ ▄▄▄▄▄
▀▄.▀·▪█·█▌▀▄.▀·    ▐█ ▌▪██•  ██ ▀▄.▀·•█▌▐█•██  
▐▀▀▪▄▐█▐█•▐▀▀▪▄    ██ ▄▄██▪  ▐█·▐▀▀▪▄▐█▐▐▌ ▐█.▪
▐█▄▄▌ ███ ▐█▄▄▌    ▐███▌▐█▌▐▌▐█▌▐█▄▄▌██▐█▌ ▐█▌·
 ▀▀▀ . ▀   ▀▀▀     ·▀▀▀ .▀▀▀ ▀▀▀ ▀▀▀ ▀▀ █▪ ▀▀▀ 

EVE Client v0.0.2
by vizn3r
 ` + "\nMain menu",
		Opts: []cli.Opt{
			{
				Name: "Input",
				Next: cli.Menu{
					Header: "Input menu",
					Opts: []cli.Opt{
						{
							Name: "Test controller",
							Func: inp.TestController,
						},
						{
							Name: "Test keyboard",
							Func: inp.TestKeyboard,
						},
					},
				},
			},
			{
				Name: "WebSocket",
				Next: wsmenu,
			},
			{
				Name: "Test",
				Func: func() {
					fmt.Println(inp.StringInp())
					time.Sleep(time.Second)
				},
			},
			{
				Name: "Exit",
				Next: exitM,
			},
		},
	}

	cli.Start(menu)
	wg.Wait()
}
