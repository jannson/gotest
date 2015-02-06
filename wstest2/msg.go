package main

type Message struct {
	Seq     int    `json:"seq"`
	MsgType string `json:"msgType"`
	Data    string `json:"data"`
}

type SeqMessage interface {
	GetSeq() int
	SetSeq(seq int)
	PutResp(m string)
	GetData() *Message
	CloseChannel()
	WaitAck()
	StopAck()
}

func (seqData *Message) GetSeq() int {
	return seqData.Seq
}

func (seqData *Message) SetSeq(seq int) {
	seqData.Seq = seq
}
