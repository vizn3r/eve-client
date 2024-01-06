package cli

import (
	"eve-client/com"
	"eve-client/inp"
	"eve-client/util"
	"fmt"
)

type Opt struct {
	Name string
	Func func()
	Next Menu
}

type Menu struct {
	Header string
	Opts   []Opt
}

func Start(m Menu) {
	m.Start()
}

var prevMenu = new(Menu)

// it works, don't touch it
func (m *Menu) Start() {
	index := 0
	prevIndex := 0
	printed := false
	trig := false
	prevIn := inp.Err

	for {
		in := inp.Inp()
		if in == prevIn {
			trig = true
		} else {
			trig = false
			prevIn = in
		}

		if index != prevIndex || !printed {
			util.Clear()
			fmt.Println("STATUS:\nKEYBOARD:", inp.KEYBOARD.IsRunning, "\nCONTROLLER:", inp.CONTROLLER.IsRunning, "\nWEBSOCKET:", com.WSCLIENT.Status.String())
			fmt.Print(m.Header, "\n\n")
			for i, o := range m.Opts {
				if index == i {
					fmt.Print("> ")
				} else {
					fmt.Print(" ")
				}
				fmt.Println(" " + o.Name)
			}
			prevIndex = index
			printed = true
		} else if in == inp.Down && !trig {
			prevIn = in

			if index == len(m.Opts)-1 {
				index = 0
			} else {
				index++
			}
		} else if in == inp.Up && !trig {
			prevIn = in

			if index == 0 {
				index = len(m.Opts) - 1
			} else {
				index--
			}
		} else if (in == inp.Back || in == inp.Left) && !trig {
			prevIn = in

			if prevMenu.Header == "" {
				return
			}
			prevMenu.Start()
		} else if (in == inp.Select || in == inp.Right) && !trig {
			prevIn = in

			if f := m.Opts[index].Func; f != nil {
				util.Clear()
				f()
				printed = false
			} else if next := m.Opts[index].Next; next.Header != "" {
				prevMenu = m
				next.Start()
			}
		}
	}
}
