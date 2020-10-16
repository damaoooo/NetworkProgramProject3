package Utils

import (
	"log"
	"net"
)

var ConnectionMap = make(map[string]net.Conn)
var MessageQueue = InitQueue()

func Errhandle(e interface{}) {
	if e != nil {
		log.Fatal(e)
	}
}
