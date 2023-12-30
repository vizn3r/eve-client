package inp

import (
	"fmt"

	"github.com/eiannone/keyboard"
)

type Keyboard struct {
	Key chan rune
	Input
}

var KEYBOARD = new(Keyboard)

func CloseKeyboard() {
	k := KEYBOARD
	close(k.Exit)
	close(k.Key)
	k.IsRunning = false
}

func OpenKeyboard() {
	k := KEYBOARD
	k.WG.Add(1)
	defer k.WG.Done()

	if err := keyboard.Open(); err != nil {
		fmt.Println("Keyboard not found")
		CloseKeyboard()
		return
	}

	// aup 65517
	// ado 65516
	// ari 65514
	// ale 65515

	for {
		char, key, err := keyboard.GetKey()
		// fmt.Println(key)
		switch key {
		case 65517:
			char = 'w'
		case 65516:
			char = 's'
		case 65514, 13:
			char = 'd'
		case 65515, 8:
			char = 'a'
		}
		if err != nil {
			panic(err)
		}
		select {
		case <- k.Exit:
			CloseKeyboard()
			return
		case k.Key <- char:
		}
	}
}	
