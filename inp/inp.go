package inp

import (
	"bufio"
	"os"
	"sync"
)

type Input struct {
	WG        *sync.WaitGroup
	Exit      chan bool
	IsRunning bool
}

type InpType int

const (
	Err InpType = iota

	Down
	Up
	Right
	Left

	Select
	Back
)

var RAW string

func Inp() InpType {
	c := CONTROLLER
	k := KEYBOARD

	cok := c.IsRunning
	kok := k.IsRunning
	if !cok && !kok {
		panic("No input device is open")
	}

	// analog treshold
	tresh := 30000
	if cok {
		axis := <-c.Axis
		butt := <-c.Buttons
		for i, a := range axis {
			switch i {
			case 0, 5:
				if a > tresh {
					return Right
				} else if a < -tresh {
					return Left
				}
			case 1, 6:
				if a > tresh {
					return Down
				} else if a < -tresh {
					return Up
				}
			}
		}
		switch butt {
		case 1:
			return Select
		case 2:
			return Back
		}
	}

	if kok {
		var out KeyboardOutput
		var char rune
		select {
		case out = <-k.Output:
		default:
			out = KeyboardOutput{0, '0'}
		}
		switch out.Key {
		case 65517:
			char = 'w'
		case 65516:
			char = 's'
		case 65514, 13:
			char = 'd'
		case 65515, 8:
			char = 'a'
		default:
			char = out.Char
		}
		switch char {
		case 'w', 'W', 'k', 'K':
			return Up
		case 's', 'S', 'j', 'J':
			return Down
		case 'd', 'D', 'l', 'L':
			return Right
		case 'a', 'A', 'h', 'H':
			return Left
		case 'q', 'Q':
			return Back
		case 'e', 'E':
			return Select
		}
	}

	return 0
}

func StringInp() string {
	s := bufio.NewScanner(os.Stdin)
	for s.Err() == nil {
		if s.Scan() {
			return s.Text()
		}
	}
	return ""
}

// func StringInp() string {
// 	buff := []byte{}
// 	for {
// 		o := <-KEYBOARD.Output
// 		util.Clear()
// 		switch o.Key {
// 		case 13:
// 			return string(buff)
// 		case 127, 8:
// 			if len(buff) > 0 {
// 				buff = buff[:len(buff)-1]
// 			}
// 		case 32:
// 			buff = append(buff, ' ')
// 		default:
// 			buff = append(buff, byte(o.Char))
// 		}
// 		fmt.Println(string(buff))
// 	}
// }
