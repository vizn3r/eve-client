package log

import (
	"fmt"
	"reflect"
	"strconv"
)

type (
	LogMode     int
	MessageType int
	ColorType   int
)

const (
	NORMAL    LogMode = iota // Normal logging
	SILENT                   // Minimal logging
	VERBOSE                  // Verbose logging
	COLORLESS                // Without color

	// Types logged with SILENT mode
	ERROR MessageType = iota - 4
	EXCEPTION

	// Types logged with NORMAL mode
	WARNING
	INFO

	// Types logged with VERBOSE mode
	MESSAGE
	DEBUG

	ANSI           string    = "\033["
	ANSI_SEPARATOR           = ";"
	ANSI_END                 = "m"
	RESET          ColorType = iota - 13
	BOLD
	FAINT
	ITALIC
	UNDERLINE
	SLOW_BLINK
	RAPID_BLINK
	INVERSE
	CONCEAL
	CROSSED_OUT
	DEFAULT
	ALT_FONT_1
	ALT_FONT_2
	ALT_FONT_3
	ALT_FONT_4
	ALT_FONT_5
	ALT_FONT_6
	ALT_FONT_7
	ALT_FONT_8
	ALT_FONT_9
	FRAKTUR
	BOLD_OFF
	NORMAL_COLOR
	ITALIC_OFF
	UNDERLINE_OFF
	BLINK_OFF
	_
	INVERSE_OFF
	REVEAl
	CROSSED_OFF
	FG_BLACK
	FG_RED
	FG_GREEN
	FG_YELLOW
	FG_BLUE
	FG_MAGENTA
	FG_CYAN
	FG_WHITE
	FG_8BIT
	FG_DEFAULT
	BG_BLACK
	BG_RED
	BG_YELLOW
	BG_BLUE
	BG_MAGENTA
	BG_CYAN
	BG_WHITE
	BG_8BIT
	BG_DEFAULT
	FRAMED
	ENCIRCLED
	OVERLINED
	FRAMED_OFF
	OVERLINED_OFF
)

// MAKE COLORS CONFIG

type LogConfig struct {
	Mode LogMode
}

type Logger struct {
	Emoji string
}

func SetColor(colors ...ColorType) (out string) {
	out += ANSI
	for i, color := range colors {
		out += strconv.Itoa(int(color))
		if i != len(colors)-1 {
			out += ANSI_SEPARATOR
		}
	}
	out += ANSI_END
	return
}

func Color(args ...any) (out string) {
	for i, arg := range args {
		if reflect.TypeOf(arg) == reflect.TypeFor[ColorType]() {
			out += SetColor(arg.(ColorType))
		} else {
			out += fmt.Sprintf("%v", arg)
			if i != len(args)-1 {
				out += " "
			}
		}
	}
	out += SetColor(RESET)
	return
}

func PrintColor(args ...any) {
	fmt.Println(Color(args...))
}

func (l *Logger) Log(msgType MessageType, args ...any) {
	buff := []any{"[" + l.Emoji + "]"}
	switch msgType {
	case ERROR:
		buff = append(buff, "[❗]")
		color := []any{FG_RED}
		color = append(color, buff...)
		color = append(color, args...)
		PrintColor(color...)
	case EXCEPTION:
		buff = append(buff, "[❗❗❗]")
		color := []any{FG_RED, BG_BLACK}
		color = append(color, buff...)
		color = append(color, args...)
		PrintColor(color...)
	case WARNING:
		buff = append(buff, "[⚠️ ]")
		color := []any{BG_YELLOW, FG_BLACK}
		color = append(color, buff...)
		color = append(color, args...)
		PrintColor(color...)
	case INFO:
		buff = append(buff, "[❔]")
		color := []any{FG_GREEN}
		color = append(color, buff...)
		color = append(color, args...)
		PrintColor(color...)
	case MESSAGE:
		buff = append(buff, "[💬]")
		color := []any{FG_WHITE}
		color = append(color, buff...)
		color = append(color, args...)
		PrintColor(color...)
	case DEBUG:
		buff = append(buff, "[⚙️ ]")
		color := []any{FG_BLUE}
		color = append(color, buff...)
		color = append(color, args...)
		PrintColor(color...)
	}
}

func (l *Logger) Error(args ...any) {
	l.Log(ERROR, args...)
}

func (l *Logger) Warning(args ...any) {
	l.Log(WARNING, args...)
}

func (l *Logger) Message(args ...any) {
	l.Log(MESSAGE, args...)
}

func (l *Logger) Info(args ...any) {
	l.Log(INFO, args...)
}
