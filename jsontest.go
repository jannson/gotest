package main

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Seq     int
	MsgType string
	Data    string
}

func main() {
	s := `{"seq":66,"msgType":"just test","data":"lkjlkj"}`
	msg := &Message{}
	fmt.Printf("the s is :%s\n", s)
	err := json.Unmarshal([]byte(s), &msg)
	if err == nil {
		fmt.Printf("msg is seq:%v %v\n", msg.Seq, msg)
	} else {
		fmt.Printf("convert error\n")
	}
}
