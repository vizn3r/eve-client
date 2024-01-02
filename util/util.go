package util

import (
	"os"
	"os/exec"
	"runtime"
)

func Clear() {
	switch runtime.GOOS {
	case "windows":
		c := exec.Command("cmd", "/c", "cls")
		c.Stdout = os.Stdout
		if e := c.Run(); e != nil {
			return
		}
	case "linux":
		c := exec.Command("printf", `\033c`)
		c.Stdout = os.Stdout
		if e := c.Run(); e != nil {
			return
		}
	}
}
