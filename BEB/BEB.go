package BEB

import (
	"fmt"
	PP2PLink "SD/PP2PLink"
)

type BEB_Module struct {
	ID      int
	PP2P    *PP2PLink.PP2PLink
	addresses []string

	Pp2plink *PP2PLink.PP2PLink
}

func NewBEB(_addresses []string, _id int, _dbg bool) *BEB_Module {
	if _dbg {
		fmt.Println("Initializing BEB")
	}

	pp2p := PP2PLink.NewPP2PLink(_addresses[_id], _dbg)

	beb := &BEB_Module{
		ID:        	_id,
		PP2P:      	pp2p,
		addresses: 	_addresses,
	}

	beb.Start()

	return beb
}

func (beb *BEB_Module) Start() {
	// Function to send predefined messages that identify the process
	go func() {
		for {
			msg := fmt.Sprintf("Hello from Process %d", beb.ID)

			for _, addr := range beb.addresses {
				beb.Send(msg, addr)
			}
		}
	}()

	// Function to handle incoming broadcast messages
	go func() {
		for {
			message := <-beb.PP2P.Ind
			beb.Deliver(message)
		}
	}()
}

func (beb *BEB_Module) Send(message string, destAddr string) {
	fmt.Printf("Process %d broadcasting to %s\n", beb.ID, destAddr)
	beb.PP2P.Req <- PP2PLink.PP2PLink_Req_Message{
		To:      destAddr,
		Message: message,
	}
}

func (p *BEB_Module) Deliver(message PP2PLink.PP2PLink_Ind_Message) {
	fmt.Printf("Process %d received from %s: %s\n", p.ID, message.From, message.Message)
}
