package main

import (
	"eve-client/cli"
	"eve-client/com"
	"eve-client/inp"
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
	k.IsRunning = true
	k.WG = &wg

	go inp.OpenKeyboard()
	go inp.OpenController(0)

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
				Name: "Connect to WS",
				Func: com.ConnectWS,
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
