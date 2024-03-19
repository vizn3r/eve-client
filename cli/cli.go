package cli

import (
	"eve-client/inp"
	"eve-client/log"
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

var (
	CURRENT_MENU Menu
	menu         = 0
)

// it works, don't touch it
// lmao, touched it :P
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
			CURRENT_MENU = *m
			util.Clear()
			fmt.Print(m.Header, "\n\n")
			for i, o := range m.Opts {
				if index == i {
					log.PrintColor("> ", log.UNDERLINE, log.BOLD, o.Name)
					// fmt.Print("> ", o.Name)
				} else {
					fmt.Println(" ", o.Name)
				}
			}
			prevIndex = index
			printed = true
		} else if in == inp.Down && !trig {
			if index == len(m.Opts)-1 {
				index = 0
			} else {
				index++
			}
		} else if in == inp.Up && !trig {
			if index == 0 {
				index = len(m.Opts) - 1
			} else {
				index--
			}
		} else if (in == inp.Back || in == inp.Left) && !trig {
			if menu > 0 {
				break
			}
		} else if (in == inp.Select || in == inp.Right) && !trig {
			if f := m.Opts[index].Func; f != nil {
				util.Clear()
				f()
				printed = false
			} else if next := m.Opts[index].Next; next.Header != "" {
				printed = false
				menu++
				next.Start()
			}
		}
	}
}
