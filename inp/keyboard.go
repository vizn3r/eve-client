package inp

import (
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
)

type Keyboard struct {
	Output chan KeyboardOutput
	Input
}

type KeyboardOutput struct {
	Key  keyboard.Key
	Char rune
}

var KEYBOARD = new(Keyboard)

func CloseKeyboard() {
	k := *KEYBOARD
	close(k.Exit)
	close(k.Output)
	k.IsRunning = false
	// _ = keyboard.Close()
}

func OpenKeyboard() {
	k := *KEYBOARD
	k.WG.Add(1)
	defer k.WG.Done()
	//
	// if err := keyboard.Open(); err != nil {
	// 	fmt.Println("Keyboard not found")
	// 	return
	// }

	// aup 65517
	// ado 65516
	// ari 65514
	// ale 65515

	for {
		char, key, err := keyboard.GetSingleKey()
		// Backdoor
		if key == 3 {
			keyboard.Close()
			os.Exit(0)
		}
		if err != nil {
			panic(err)
		}
		select {
		case <-k.Exit:
			keyboard.Close()
			return
		case k.Output <- KeyboardOutput{key, char}:
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
