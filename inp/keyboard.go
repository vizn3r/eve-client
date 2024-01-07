package inp

import (
	"eve-client/serv"
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
)

type Keyboard struct {
	Output chan KeyboardOutput
	GetKey chan bool
	Input
}

type KeyboardOutput struct {
	Key  keyboard.Key
	Char rune
}

var KEYBOARD = new(Keyboard)

func CloseKeyboard() {
	close(KEYBOARD.Exit)
	close(KEYBOARD.Output)
	KEYBOARD.Status = serv.STARTING
	keyboard.Close()
}

func OpenKeyboard() {
	KEYBOARD.WG.Add(1)
	defer KEYBOARD.WG.Done()

	// Arrow key codes
	// aup 65517
	// ado 65516
	// ari 65514
	// ale 65515

	KEYBOARD.Status = serv.RUNNING
	for {
		var char rune
		var key keyboard.Key
		var err error
		select {
		case gk := <-KEYBOARD.GetKey:
			if gk {
				char, key, err = keyboard.GetSingleKey()
			}
		default:
			KEYBOARD.GetKey <- true
		}
		// Backdoor
		if key == 3 {
			keyboard.Close()
			os.Exit(0)
		}
		if err != nil {
			panic(err)
		}
		select {
		case <-KEYBOARD.Exit:
			keyboard.Close()
			return
		case KEYBOARD.Output <- KeyboardOutput{key, char}:
			keyboard.Close()
		}
	}
}

func TestKeyboard() {
	for {
		o := <-KEYBOARD.Output
		if o.Key == 13 {
			return
		}
		fmt.Println(o)
	}
}
