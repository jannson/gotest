package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"
	"time"
)

var addr = flag.String("addr", ":8080", "http service address")
var homeTempl = template.Must(template.ParseFiles("home.html"))

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTempl.Execute(w, r.Host)
}

type WsMessage struct {
	Message
	resp    chan string
	ackTick *time.Ticker
}

func (wsMsg *WsMessage) PutResp(m string) {
	wsMsg.resp <- m
}

func (wsMsg *WsMessage) GetData() *Message {
	return &wsMsg.Message
}

func (wsMsg *WsMessage) CloseChannel() {
	close(wsMsg.resp)
}

func (wsMsg *WsMessage) WaitAck() {
	if wsMsg.Seq < 0 {
		panic("the seq should not < 0")
	}
	wsMsg.ackTick = time.NewTicker(time.Second * 5)
	go func(seq int) {
		<-wsMsg.ackTick.C
		h.timeouts <- seq
	}(wsMsg.GetSeq())
}

func (wsMsg *WsMessage) StopAck() {
	wsMsg.ackTick.Stop()
}

func serveTest(w http.ResponseWriter, r *http.Request) {
	msg := &WsMessage{Message{-1, "test", "hello"}, make(chan string, 0), nil}
	h.messages <- msg

	respData, ok := <-msg.resp
	if ok {
		msg.CloseChannel()
		io.WriteString(w, fmt.Sprintf("got message %v", respData))
	} else {
		io.WriteString(w, "got message timeout")
	}
}

func main() {
	flag.Parse()
	go h.Run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/test", serveTest)
	http.HandleFunc("/ws", serveWs)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
