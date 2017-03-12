package queue

import (
	. "../types"
	"time"
)

func BestSuitedElevator(e map[string]ElevData_s, buttonPress ButtonPress_s) string {
	bestTime := 1 * time.Hour
	bestIP := ""
	for IP, elevData := range e {
		elevDataCopy := elevData
		elevDataCopy.Orders[buttonPress.Floor][buttonPress.ButtonType] = 1
		time := timeToIdle(elevDataCopy)
		if time < bestTime {
			bestTime = time
			bestIP = IP
		}
	}
	return bestIP
}

func timeToIdle(e ElevData_s) time.Duration {
	dur := 0 * time.Millisecond

	switch e.State {
	case IDLE:
		e.Direction = SetDir(e)
		if e.Direction == DIRN_STOP {
			return dur
		}
	case MOVING:
		e.Floor = e.Floor + int(e.Direction)
		dur += TravelTime
	case DOOROPEN:
		dur -= DoorOpenTime
	}

	for {
		if ShouldStop(e) {
			DeleteLocalOrder(&e)
			dur += DoorOpenTime
			e.Direction = SetDir(e)
			if e.Direction == DIRN_STOP {
				return dur
			}
		}
		e.Floor = e.Floor + int(e.Direction)
		dur += TravelTime
	}
}

func IsLocalEmpty(e ElevData_s) bool {
	for nf := 0; nf < N_FLOORS; nf++ {
		for nb := 0; nb < N_BUTTONS; nb++ {
			if e.Orders[nf][nb] == 1 {
				return false
			}
		}
	}
	return true
}

func CheckOrdersAbove(e ElevData_s) bool {
	for nf := e.Floor + 1; nf < N_FLOORS; nf++ {
		for nb := 0; nb < N_BUTTONS; nb++ {
			if e.Orders[nf][nb] == 1 {
				return true
			}
		}
	}
	return false
}

func CheckOrdersBelow(e ElevData_s) bool {
	for nf := 0; nf < e.Floor; nf++ {
		for nb := 0; nb < N_BUTTONS; nb++ {
			if e.Orders[nf][nb] == 1 {
				return true
			}
		}
	}
	return false
}

func DeleteLocalOrder(e *ElevData_s) {
	e.Orders[e.Floor][BUTTON_CALL_UP] = 0
	e.Orders[e.Floor][BUTTON_CALL_DOWN] = 0
	e.Orders[e.Floor][BUTTON_COMMAND] = 0
}

func DeleteAll(e *ElevData_s) {
	for nf := 0; nf < N_FLOORS; nf++ {
		for nb := 0; nb < N_BUTTONS; nb++ {
			e.Floor = nf
			DeleteLocalOrder(e)
		}
	}
}

func ShouldStop(e ElevData_s) bool {
	switch e.Direction {
	case DIRN_DOWN:
		if e.Orders[e.Floor][BUTTON_CALL_DOWN] == 1 ||
			e.Orders[e.Floor][BUTTON_COMMAND] == 1 ||
			CheckOrdersBelow(e) == false {
			return true
		}

	case DIRN_UP:
		if e.Orders[e.Floor][BUTTON_CALL_UP] == 1 ||
			e.Orders[e.Floor][BUTTON_COMMAND] == 1 ||
			CheckOrdersAbove(e) == false {
			return true
		}
	default:
		return true
	}
	return false
}

func SetDir(e ElevData_s) ElevMotorDir_t {
	if IsLocalEmpty(e) {
		return DIRN_STOP
	}
	switch e.Direction {
	case DIRN_UP:
		if CheckOrdersAbove(e) {
			return DIRN_UP
		} else {
			return DIRN_DOWN
		}
	case DIRN_DOWN:
		if CheckOrdersBelow(e) {
			return DIRN_DOWN
		} else {
			return DIRN_UP
		}
	case DIRN_STOP:
		if CheckOrdersAbove(e) {
			return DIRN_UP
		} else if CheckOrdersBelow(e) {
			return DIRN_DOWN
		} else {
			return DIRN_STOP
		}
	default:
		return DIRN_STOP 
	}
}