package inp

import (
	"eve-client/serv"
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
	close(CONTROLLER.Exit)
	close(CONTROLLER.Buttons)
	close(CONTROLLER.Axis)
	CONTROLLER.Status = serv.STOPPED
}

func ControllerIsReady() bool {
	if !CONTROLLER.IsRunning() {
		fmt.Println("[🎮]<⚠️ > Controller is not connected")
		time.Sleep(time.Second)
	}
	return CONTROLLER.IsRunning()
}

func OpenController(id int) {
	CONTROLLER.WG.Add(1)
	defer CONTROLLER.WG.Done()
	defer CloseController()

	js, err := joystick.Open(id)
	if err != nil {
		fmt.Println("[🎮]<⚠️ > Controller not found")
		time.Sleep(time.Second)
		return
	}
	defer js.Close()

	CONTROLLER.Status = serv.RUNNING

	for {
		state, err := js.Read()
		if err != nil {
			fmt.Println("[🎮]<⚠️ > Controller not found")
			time.Sleep(time.Second)
			return
		}
		select {
		case <-CONTROLLER.Exit:
			return
		case CONTROLLER.Buttons <- state.Buttons:
		case CONTROLLER.Axis <- state.AxisData:
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
