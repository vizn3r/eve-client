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

// func test() {
// 	fmt.Println("Test")
// 	time.Sleep(time.Second * 1)
// }

func main() {
	var wg sync.WaitGroup

	c := inp.CONTROLLER
	c.Axis = make(chan []int)
	c.Buttons = make(chan uint32)
	c.Exit = make(chan bool)
	c.IsRunning = true
	c.WG = &wg

	k := inp.KEYBOARD
	k.Output = make(chan inp.KeyboardOutput)
	k.Exit = make(chan bool)
	k.GetKey = make(chan bool, 1)
	k.IsRunning = true
	k.WG = &wg

	com.WSCLIENT.Msg = make(chan string)
	com.WSCLIENT.MsgRes = make(chan string)
	com.WSCLIENT.Status = serv.STARTING

	go inp.OpenKeyboard()
	k.GetKey <- true
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
			// {
			// 	Name: "Connect",
			// 	Func: ,
			// },
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
