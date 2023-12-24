package inp

import (
	"fmt"

	"github.com/0xcafed00d/joystick"
)

// Controller outputs:
// AXIS[J0X, J0Y, L2/R2, J1Y, J1X, DPX, DPY]

type Controller struct {
	Buttons chan uint32
	Axis chan []int
	Input
}

var CONTROLLER = new(Controller)

func OpenController(id int) {
	c := CONTROLLER
	c.WG.Add(1)
	defer c.WG.Done()

	js, err := joystick.Open(id)
	if err != nil {
		panic(err)
	}
	defer js.Close()

	for {
		state, err := js.Read()
		if err != nil {
			panic(err)
		}
		select {
		case <- c.Exit:
			close(c.Exit)
			close(c.Buttons)
			close(c.Axis)
			c.IsRunning = false
			return
		case c.Buttons <- state.Buttons:
		case c.Axis <- state.AxisData:
		}
	}
}

func TestController() {
	for {
		fmt.Println(<- CONTROLLER.Axis, <- CONTROLLER.Buttons)
		if axis := <- CONTROLLER.Axis; axis[0] > 30000 && <- CONTROLLER.Buttons == 2 {
			return
		}
	}
}