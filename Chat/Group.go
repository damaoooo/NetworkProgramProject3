package Chat

import (
	"NPProj3/ORM"
	"NPProj3/Utils"
	"encoding/json"
	"myTest/Unit"
	"net"
)

func GroupChat(connect net.Conn, request ORM.MessageBlock) {
	Utils.MessageQueue.Add(request.Event)
	retJson := ORM.CommonResponse{
		Result: "success",
		Uuid:   request.Uuid,
	}
	ret, err := json.Marshal(retJson)
	Unit.ErrHandle(err)
	_, err = connect.Write(ret)
	Unit.ErrHandle(err)
}
