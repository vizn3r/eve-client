package inp

import (
	"bufio"
	"eve-client/log"
	"eve-client/serv"
	"fmt"
	"os"
	"sync"
)

type Input struct {
	WG   *sync.WaitGroup
	Exit chan bool
	serv.Service
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

var (
	RAW string
	LOG = log.Logger{
		Emoji: "🔤",
	}
)

func Inp() InpType {
	cok := CONTROLLER.IsRunning()
	kok := KEYBOARD.IsRunning()
	for !cok && !kok {
		cok = CONTROLLER.IsRunning()
		kok = KEYBOARD.IsRunning()
		fmt.Println(cok, kok)
	}
	if !cok && !kok {
		fmt.Println("No input device open")
		os.Exit(0)
	}

	// analog treshold
	tresh := 30000
	if cok {
		axis := <-CONTROLLER.Axis
		butt := <-CONTROLLER.Buttons
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
		case out = <-KEYBOARD.Output:
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

func WaitForBack() {
	LOG.Message("Press ")
	for {
		in := Inp()
		if in == Back || in == Left {
			return
		}
	}
}

func WaitForAny() {
	LOG.Message("Press 'any' movement")
	for {
		in := Inp()
		if in != 0 {
			return
		}
	}
}
