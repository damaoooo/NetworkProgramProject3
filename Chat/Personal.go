package Chat

import (
	"NPProj3/ORM"
	"NPProj3/Utils"
	"NPProj3/Wigets"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"net"
)

// go round a conn queue to find and send message
func PersonalChat(connection net.Conn, request ORM.MessageBlock) {
	newUUID := uuid.Must(uuid.NewV4()).String()
	newEventReturn := ORM.Event{
		Type: "personal_chat",
		User: request.Username,
		Case: request.Plain,
	}
	eventResp := ORM.EventResponse{
		Uuid:        newUUID,
		MessageType: "personal",
		Event:       newEventReturn,
	}
	ret, err := json.Marshal(eventResp)
	Wigets.ErrHandle(err)
	commonReturn := ORM.CommonResponse{
		Result: "success",
		Uuid:   request.Uuid,
	}
	commonReturnByte, err := json.Marshal(commonReturn)
	Wigets.ErrHandle(err)
	Wigets.SendBuf(connection, commonReturnByte)

	for username, conn := range Utils.ConnectionMap {
		if username == request.SendTo {
			Wigets.SendBuf(conn, ret)
		}
	}

}
