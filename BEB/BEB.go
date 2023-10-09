package BEB

import (
	"fmt"
	"strconv"
	"strings"
	PP2PLink "SD/PP2PLink"
)

type BroadcastMessage struct {
	SenderID   int
	Message    string
}

type BEB_Module struct {
	ID               int
	PP2P             *PP2PLink.PP2PLink
	BroadcastMessageChannel_IN chan BroadcastMessage
	BroadcastMessageChannel_OUT chan BroadcastMessage
	addresses        []string
}

func NewBEB(_addresses []string, _id int, _dbg bool) *BEB_Module {
	if _dbg {
		fmt.Println("Initializing BEB")
	}

	pp2p := PP2PLink.NewPP2PLink(_addresses[_id], _dbg)

	beb := &BEB_Module{
		ID:                    _id,
		PP2P:                  pp2p,
		BroadcastMessageChannel_IN: make(chan BroadcastMessage),
		BroadcastMessageChannel_OUT: make(chan BroadcastMessage),
		addresses:             _addresses,
	}

	beb.Start()

	return beb
}

func (beb *BEB_Module) Start() {
	// Function to send predefined messages that identify the process
	go func() {
		for {
			msg := <-beb.BroadcastMessageChannel_IN
			for _, addr := range beb.addresses {
                beb.Send(msg, addr)
            }		
		}
	}()

	// Function to handle incoming broadcast messages
	go func() {
		for {
			message := <- beb.PP2P.Ind
			beb.Deliver(message)
		}
	}()
}

func (beb *BEB_Module) Send(bm BroadcastMessage, destAddr string) {
    fmt.Printf("Process %d broadcasting to %s: %s\n", bm.SenderID, destAddr, bm.Message)
	msg := fmt.Sprintf("%d/%s", bm.SenderID, bm.Message)
	beb.PP2P.Req <- PP2PLink.PP2PLink_Req_Message{
        To:        destAddr,
        Message:   msg,
    }
}

func (beb *BEB_Module) Deliver(pp2p PP2PLink.PP2PLink_Ind_Message) {
	split := strings.Split(pp2p.Message, "/")
	
	senderID, _ := strconv.Atoi(split[0])
	message := split[1]
	
	broadcast := BroadcastMessage{
		SenderID: senderID,
		Message:  message,
	}
	
	fmt.Printf("Process %d received from %d: %s\n", beb.ID, broadcast.SenderID, broadcast.Message)

	// beb.BroadcastMessageChannel_OUT <- broadcast
}