package log

import "fmt"

/* TODO: LOG
- LOG config
- LOG modes
- LOG file
*/

type LogOpts struct {
	Verbose bool
	Silent  bool
	// Colors bool
	// Emoji bool
	// Format ?
}

func LogInit(opts LogOpts) {
}

func Log(v ...any) {
	fmt.Println(v...)
}
