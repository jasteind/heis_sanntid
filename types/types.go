package types

const MOTOR_SPEED = 2800

const N_FLOORS = 4
const N_BUTTONS = 3

const TravelTime = 2
const DoorOpenTime = 3

const (
	IDLE = iota
	MOVING
	DOOROPEN
)

type ElevMotorDir_t int

const (
	DIRN_DOWN ElevMotorDir_t = -1
	DIRN_STOP ElevMotorDir_t = 0
	DIRN_UP   ElevMotorDir_t = 1
)

type ElevButton_t int

const (
	BUTTON_CALL_UP   ElevButton_t = 0
	BUTTON_CALL_DOWN ElevButton_t = 1
	BUTTON_COMMAND   ElevButton_t = 2
)



type ButtonPress_s struct {
	ButtonType ElevButton_t
	Floor      int
}

type ElevData_s struct {
	IP        string
	Direction ElevMotorDir_t
	Floor     int
	Orders    [N_FLOORS][N_BUTTONS]int
	State     int
}

type Order_s struct {
	OwnerID    string
	ButtonType ElevButton_t
	Floor      int
	Status     OrderStatus
}

type OrderStatus int

const (
	OS_NEW OrderStatus = iota
	OS_RECEIVED
	OS_COMPLETE
)

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}