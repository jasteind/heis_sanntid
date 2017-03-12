package fsm

import (
	"../driver"
	"../queue"
	."../types"
	"fmt"
	"time"
)

func Run(
	e ElevData_s,
	FloorReached_ch chan int,
	NewOrder_ch chan Order_s,
	ElevDataSend_ch chan ElevData_s,
	CompleteOrder_ch chan Order_s) {

	var DoorClosed_ch <-chan time.Time

	fmt.Printf("[FSM] %+v\n", e)

	elevData := ElevData_s{e.IP, e.Direction, e.Floor, e.Orders, e.State}
	ElevDataSend_ch <- elevData

	for {
		select {
		case floor := <-FloorReached_ch:
			fmt.Printf("[FSM] New floor %+v\n", floor)

			driver.SetFloorIndicator(floor)
			e.Floor = floor

			switch e.State {
			case IDLE:
				break
			case MOVING:
				if queue.ShouldStop(e) {
					driver.SetMotorDir(DIRN_STOP)

					queue.DeleteLocalOrder(&e)
					CompleteOrder_ch <- Order_s{e.IP, BUTTON_CALL_UP, e.Floor, OS_COMPLETE}
					CompleteOrder_ch <- Order_s{e.IP, BUTTON_CALL_DOWN, e.Floor, OS_COMPLETE}
					CompleteOrder_ch <- Order_s{e.IP, BUTTON_COMMAND, e.Floor, OS_COMPLETE}

					driver.SetDoorOpenLamp(1)
					DoorClosed_ch = time.After(DoorOpenTime * time.Second)
					e.State = DOOROPEN
				}
				break
			}

		case order := <-NewOrder_ch:
			fmt.Printf("[FSM] New order: %+v\n", order)
			
			switch e.State {
			case IDLE:
				e.Orders[order.Floor][int(order.ButtonType)] = 1
				e.Direction = queue.SetDir(e)
				if e.Direction == DIRN_STOP {
					driver.SetDoorOpenLamp(1)

					queue.DeleteLocalOrder(&e)
					order := Order_s{OwnerID: e.IP, ButtonType: order.ButtonType, Floor: e.Floor, Status: OS_COMPLETE}
					CompleteOrder_ch <- order

					DoorClosed_ch = time.After(DoorOpenTime * time.Second)
				} else {
					driver.SetMotorDir(e.Direction)
					e.State = MOVING
				}
				break

			case DOOROPEN:

				e.Orders[order.Floor][int(order.ButtonType)] = 1
				if order.Floor == e.Floor {
					DoorClosed_ch = time.After(DoorOpenTime * time.Second)
				}
				break

			case MOVING:
				e.Orders[order.Floor][int(order.ButtonType)] = 1
				break
			default:
				break
			}

		case <-DoorClosed_ch:
			fmt.Println("Time out!")

			switch e.State {
			case IDLE:
				break

			case DOOROPEN:
				driver.SetDoorOpenLamp(0)
				direction := queue.SetDir(e)
				e.Direction = direction
				if direction == DIRN_STOP {
					e.State = IDLE
				} else {
					driver.SetMotorDir(direction)
					e.State = MOVING
					e.Direction = direction
				}
				break
			}
		}
		elevData := ElevData_s{e.IP, e.Direction, e.Floor, e.Orders, e.State}
		ElevDataSend_ch <- elevData
	}
}

func Init(FloorReached_ch chan int) ElevData_s {

	driver.Init()
	driver.SetMotorDir(DIRN_DOWN)
	for {
		if driver.GetFloorSensorSignal() != -1 {
			driver.SetMotorDir(DIRN_STOP)
			break
		}
	}
	f := driver.GetFloorSensorSignal()

	return ElevData_s{
		State:     IDLE,
		Floor:     f,
		Direction: DIRN_STOP,
	}
}