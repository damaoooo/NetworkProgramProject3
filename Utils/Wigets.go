package Utils

import (
	"log"
	"net"
)

var ConnectionMap = make(map[string]net.Conn)
var MessageQueue = InitQueue()
var MessageBox = make(chan int)

func ErrHandle(e interface{}) {
	if e != nil {
		log.Fatal(e)
	}
}

func MessageListen() {
	for {
		if !MessageQueue.IsEmpty() {
			MessageBox <- 1
		}
	}
}
