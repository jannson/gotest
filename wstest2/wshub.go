package main

import "fmt"

//http://stackoverflow.com/questions/13511203/why-cant-i-assign-a-struct-to-an-interface
type WsHub struct {
	connections map[*Connection]bool
	register    chan *Connection
	unregister  chan *Connection
	messages    chan SeqMessage
	timeouts    chan int
	responses   chan Message
	seqMap      *SeqMap
}

var h = WsHub{
	connections: make(map[*Connection]bool),
	register:    make(chan *Connection),
	unregister:  make(chan *Connection),
	messages:    make(chan SeqMessage, 100),
	timeouts:    make(chan int, 100),
	responses:   make(chan Message, 100),
	seqMap:      NewSeqMap(100),
}

func (h *WsHub) Run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.sendMsg)
			}
		case m := <-h.messages:
			oldData := h.seqMap.NewSeq(m)
			if oldData != nil {
				oldData.(SeqMessage).CloseChannel()
			}

			m.WaitAck()
			for c := range h.connections {
				select {
				case c.sendMsg <- m.GetData():
				default:
					close(c.sendMsg)
					delete(h.connections, c)
				}
			}
		case resp := <-h.responses:
			fmt.Printf("bus got resp seq %d", resp.Seq)
			data := h.seqMap.GetData(resp.Seq)
			if data != nil {
				//Not close channel at ok
				seqMsg := data.(SeqMessage)
				seqMsg.StopAck()
				seqMsg.PutResp(resp.Data)
			} else {
				fmt.Printf("data not find in seqmap\n")
			}
		case timeoutSeq := <-h.timeouts:
			fmt.Printf("bus got timeout seq %d", timeoutSeq)
			data := h.seqMap.GetData(timeoutSeq)
			if data != nil {
				seqMsg := data.(SeqMessage)
				seqMsg.CloseChannel()
			} else {
				fmt.Printf("data not find in seqmap\n")
			}
		}
	}
}
