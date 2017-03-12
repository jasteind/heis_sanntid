package driver

import (
	. "../types"
	"fmt"
	"time"
)

var lampChannelMatrix = [N_FLOORS][N_BUTTONS]int{

	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var buttonChannelMatrix = [N_FLOORS][N_BUTTONS]int{

	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

func Init() {

	initSucsess := io_init()
	if initSucsess == false {
		fmt.Println("Feil i initiering av heis")
	}

	for f := 0; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTONS; b++ {
			SetButtonLamp(ElevButton_t(b), f, 0)
		}
	}
	SetDoorOpenLamp(0)
	SetFloorIndicator(0)
}

func Event(ButtonPressed_ch chan ButtonPress_s, FloorReached_ch chan int) {
	var prevFloor int
	var prevButtons [N_FLOORS][N_BUTTONS]bool

	for {
		floor := GetFloorSensorSignal()
		if floor != prevFloor && floor != -1 {
			FloorReached_ch <- floor
		}
		prevFloor = floor

		for floor := 0; floor < N_FLOORS; floor++ {
			for button := BUTTON_CALL_UP; int(button) < N_BUTTONS; button++ {
				b := GetButtonSignal(ElevButton_t(button), floor)
				if b != prevButtons[floor][int(button)] && b {
					buttonPressed := ButtonPress_s{ButtonType: ElevButton_t(button), Floor: floor}
					ButtonPressed_ch <- buttonPressed
				}
				prevButtons[floor][button] = b
			}
		}
		time.Sleep(25 * time.Millisecond)
	}
}

func SetFloorIndicator(floor int) {

	if !(floor >= 0) {
		panic(fmt.Sprintf("Error: Floor out of reach"))
	}
	if !(floor < N_FLOORS) {
		panic(fmt.Sprintf("Error: Too many floors"))
	}
	if floor&0x02 != 0 {
		io_set_bit(LIGHT_FLOOR_IND1)
	} else {
		io_clear_bit(LIGHT_FLOOR_IND1)
	}
	if floor&0x01 != 0 {
		io_set_bit(LIGHT_FLOOR_IND2)
	} else {
		io_clear_bit(LIGHT_FLOOR_IND2)
	}
}

func SetMotorDir(dirn ElevMotorDir_t) {

	if dirn == 0 {
		io_write_analog(MOTOR, 0)
	} else if dirn > 0 {
		io_clear_bit(MOTORDIR)
		io_write_analog(MOTOR, MOTOR_SPEED)
	} else if dirn < 0 {
		io_set_bit(MOTORDIR)
		io_write_analog(MOTOR, MOTOR_SPEED)
	}
}

func SetButtonLamp(button ElevButton_t, floor int, value int) {

	if !(floor >= 0) {
		panic(fmt.Sprintf("Error: Floor out of reach"))
	}
	if !(floor < N_FLOORS) {
		panic(fmt.Sprintf("Error: Too  many floors"))
	}
	if !(button >= 0) {
		panic(fmt.Sprintf("Error: Button is less than 0"))
	}
	if !(button < N_BUTTONS) {
		panic(fmt.Sprintf("Error: Button is larger than 3"))
	}

	if !(value == 0) {
		io_set_bit(lampChannelMatrix[floor][button])
	} else {
		io_clear_bit(lampChannelMatrix[floor][button])
	}
}

func SetDoorOpenLamp(value int) {

	if !(value == 0) {
		io_set_bit(LIGHT_DOOR_OPEN)

	} else {
		io_clear_bit(LIGHT_DOOR_OPEN)
	}
}

func SetStopLamp(value int) {

	if !(value == 0) {
		io_set_bit(LIGHT_STOP)

	} else {
		io_clear_bit(LIGHT_STOP)
	}
}

func GetButtonSignal(button ElevButton_t, floor int) bool {

	if !(floor >= 0) {
		panic(fmt.Sprintf("Error: Floor  out of reach"))
	}
	if !(floor < N_FLOORS) {
		panic(fmt.Sprintf("Error: Too  many floors"))
	}
	if !(button >= 0) {
		panic(fmt.Sprintf("Error: Button is less than 0"))
	}
	if !(button < N_BUTTONS) {
		panic(fmt.Sprintf("Error: Button is larger than 3"))
	}
	return io_read_bit(buttonChannelMatrix[floor][button])
}

func GetFloorSensorSignal() int {

	if io_read_bit(SENSOR_FLOOR1) {
		return 0
	} else if io_read_bit(SENSOR_FLOOR2) {
		return 1
	} else if io_read_bit(SENSOR_FLOOR3) {
		return 2
	} else if io_read_bit(SENSOR_FLOOR4) {
		return 3
	} else {
		return -1
	}
}
