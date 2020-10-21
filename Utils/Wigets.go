package Utils

import (
	"NPProj3/ORM"
	"encoding/json"
	"log"
	"net"
)

var ConnectionMap = make(map[string]net.Conn)
var MessageQueue = InitQueue()

func ErrHandle(e interface{}) {
	if e != nil {
		log.Fatal(e)
	}
}

func UnSerialize(buf []byte) (string, *ORM.MessageBlock) {
	recvJson := new(ORM.MessageBlock)
	confirmJson := new(ORM.CommonResponse)
	err := json.Unmarshal(buf, confirmJson)
	err = json.Unmarshal(buf, recvJson)
	ErrHandle(err)
	if confirmJson.Result == "success" {
		return "ack", nil
	} else {
		return recvJson.MessageType, recvJson
	}
}
