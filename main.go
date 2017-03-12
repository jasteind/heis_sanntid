package main

import (
	"./bcast"
	"./driver"
	"./localip"
	"./queue"
	"./peers"
	"./backup"
	"./fsm"
	"fmt"
	"os"
	"os/signal"
	. "./types"
)



func main() {

	//Event channels
	FloorReached_ch := make(chan int)
	ButtonPressed_ch := make(chan ButtonPress_s)
	NewOrder_ch := make(chan Order_s, 20)
	CompleteOrder_ch := make(chan Order_s)
	GlobalOrder_send_ch := make(chan Order_s)
	GlobalOrder_recv_ch := make(chan Order_s)
	ElevData_fsm_ch := make(chan ElevData_s)
	ElevData_send_ch := make(chan ElevData_s)
	ElevData_recv_ch := make(chan ElevData_s)

	//Peers channel
	PeerUpdate_ch := make(chan PeerUpdate)
	TransmitEnable_ch := make(chan bool)

	port := 12000
	myIP, error := localip.LocalIP()
	if error != nil {
		fmt.Println("Colud not fetch IP")
	}

	elevators := make(map[string]ElevData_s)

	go driver.Event(ButtonPressed_ch, FloorReached_ch)

	go bcast.Transmitter(port, GlobalOrder_send_ch, ElevData_send_ch)
	go bcast.Receiver(port, GlobalOrder_recv_ch, ElevData_recv_ch)
	go peers.Transmitter(port+1, myIP, TransmitEnable_ch)
	go peers.Receiver(port+1, PeerUpdate_ch)

	{
		e := fsm.Init(FloorReached_ch)
		e.IP = myIP
		backup.LoadOrder(&e)

		for f := range e.Orders {
			for b := range e.Orders[f] {
				if e.Orders[f][b] == 1 {
					NewOrder_ch <- Order_s{myIP, ElevButton_t(b), f, OS_NEW}
					driver.SetButtonLamp(ElevButton_t(b), f, 1)
				}
			}
		}
		elevators[myIP] = e
	}


	go fsm.Run(elevators[myIP], FloorReached_ch, NewOrder_ch, ElevData_fsm_ch, CompleteOrder_ch)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
	    <-c
        driver.SetMotorDir(DIRN_STOP)
        os.Exit(0)
	}()

	for {
		select {
		case a := <-ButtonPressed_ch:
			fmt.Printf("[MAIN] Button pressed: %+v\n", a)
			if a.ButtonType == BUTTON_COMMAND {
				NewOrder_ch <- Order_s{myIP, a.ButtonType, a.Floor, OS_NEW}
				driver.SetButtonLamp(BUTTON_COMMAND, a.Floor, 1)

			} else {
				bestIP := queue.BestSuitedElevator(elevators, a)
				fmt.Println("best: ", bestIP)
				order := Order_s{bestIP, a.ButtonType, a.Floor, OS_NEW}
				GlobalOrder_send_ch <- order
			}
		case a := <-CompleteOrder_ch:
			fmt.Printf("[MAIN] Complete order: %+v\n", a)
			if a.ButtonType == BUTTON_COMMAND {
				driver.SetButtonLamp(BUTTON_COMMAND, a.Floor, 0)
			} else {
				GlobalOrder_send_ch <- a
			}

		case a := <-GlobalOrder_recv_ch:
			fmt.Printf("[MAIN] Order from net: %+v\n", a)
			switch a.Status {
			case OS_NEW:
				if a.OwnerID == myIP {
					order := Order_s{myIP, a.ButtonType, a.Floor, OS_RECEIVED}
					GlobalOrder_send_ch <- order
					NewOrder_ch <- a
				}
			case OS_RECEIVED:
				driver.SetButtonLamp(a.ButtonType, a.Floor, 1)
			case OS_COMPLETE:
				driver.SetButtonLamp(a.ButtonType, a.Floor, 0)
			}

		case a := <-ElevData_recv_ch:
			fmt.Printf("[MAIN] Remote elevData: %+v\n", a)
			if a.IP != myIP {
				elevators[a.IP] = a
			}

		case a := <-ElevData_fsm_ch:
			fmt.Printf("[MAIN] Local elevData: %+v\n", a)
			elevators[myIP] = a
			backup.SaveOrder(elevators[myIP])
			ElevData_send_ch <- a


		case a := <-PeerUpdate_ch:
			fmt.Printf("[MAIN] Peers: %+v\n", a.Peers)
			fmt.Printf("[MAIN] New Peers: %+v\n", a.New)
			fmt.Printf("[MAIN] Lost Peers: %+v\n", a.Lost)

			for _, lostID := range a.Lost {
				lostElev := elevators[lostID]
				for f := range lostElev.Orders {
					for b := 0; b < N_BUTTONS-1; b++ {
						if lostElev.Orders[f][b] == 1 {
							NewOrder_ch <- Order_s{myIP, ElevButton_t(b), f, OS_NEW}
						}
					}
				}
			}
		}
	}
}