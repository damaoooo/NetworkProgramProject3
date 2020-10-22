package Utils

import (
	"NPProj3/ORM"
	"encoding/json"
	"log"
	"net"
	"sync"
)

var ConnectionMap = make(map[string]net.Conn)
var MessageQueue = InitQueue()
var SessionMaps = make(map[string]Session)
var SessionM = SessionManager{
	Sessions: SessionMaps,
	Lock:     sync.Mutex{},
	Err:      SessionError{},
}

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

func SessionValidate(req ORM.MessageBlock, conn net.Conn) bool {
	if req.MessageType == "login" {
		return true
	} else {
		if SessionM.IsValid(req.Session) {
			return true
		} else {
			failedJson := ORM.WrongSession{Info: "Go Away!"}
			failedByte, err := json.Marshal(failedJson)
			ErrHandle(err)
			_, err = conn.Write(failedByte)
			ErrHandle(err)
			return false
		}
	}
}
