package Chat

import (
	"NPProj3/ORM"
	"NPProj3/Utils"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"myTest/Unit"
	"net"
)

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
	Unit.ErrHandle(err)
	commonReturn := ORM.CommonResponse{
		Result: "success",
		Uuid:   request.Uuid,
	}
	commonReturnByte, err := json.Marshal(commonReturn)
	Unit.ErrHandle(err)
	_, err = connection.Write(commonReturnByte)
	Unit.ErrHandle(err)

	for username, conn := range Utils.ConnectionMap {
		if username == request.SendTo {
			_, err = conn.Write(ret)
			Unit.ErrHandle(err)
		}
	}

}
