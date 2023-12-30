package inp

import (
	"sync"
)

type Input struct {
	WG *sync.WaitGroup
	Exit chan bool
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
		axis := <- c.Axis
		butt := <- c.Buttons
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
		var key rune
		select {
		case key = <-k.Key:
		default:
			key = '0'
		}
		switch key {
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
