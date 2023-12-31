package inp

import (
	"fmt"
	"time"

	"github.com/0xcafed00d/joystick"
)

// Controller outputs:
// AXIS[J0X, J0Y, L2/R2, J1Y, J1X, DPX, DPY]

type Controller struct {
	Buttons chan uint32
	Axis    chan []int

	Input
}

var CONTROLLER = new(Controller)

// Close CONTROLLER channels and set IsRunning to false
func CloseController() {
	c := CONTROLLER
	close(c.Exit)
	close(c.Buttons)
	close(c.Axis)
	c.IsRunning = false
}

func ControllerIsReady() bool {
	if !CONTROLLER.IsRunning {
		fmt.Println("Controller is not connected")
		time.Sleep(time.Second)
	}
	return CONTROLLER.IsRunning
}

func OpenController(id int) {
	c := CONTROLLER
	c.WG.Add(1)
	defer c.WG.Done()
	defer CloseController()

	js, err := joystick.Open(id)
	if err != nil {
		fmt.Println("Controller not found")
		return
	}
	defer js.Close()

	for {
		state, err := js.Read()
		if err != nil {
			fmt.Println("Controller not found")
			return
		}
		select {
		case <-c.Exit:
			return
		case c.Buttons <- state.Buttons:
		case c.Axis <- state.AxisData:
		}
	}
}

func TestController() {
	if !ControllerIsReady() {
		return
	}
	for {
		fmt.Println(<-CONTROLLER.Axis, <-CONTROLLER.Buttons)
		if axis := <-CONTROLLER.Axis; axis[0] > 30000 && <-CONTROLLER.Buttons == 2 {
			return
		}
	}
}
