package Chat

import (
	"NPProj3/ORM"
	"NPProj3/Utils"
	"NPProj3/Wigets"
	"encoding/json"
	"net"
)

func GroupChat(connect net.Conn, request ORM.MessageBlock) {
	Utils.MessageQueue.Add(request.Event)
	retJson := ORM.CommonResponse{
		Result: "success",
		Uuid:   request.Uuid,
	}
	ret, err := json.Marshal(retJson)
	Wigets.ErrHandle(err)
	Wigets.SendBuf(connect, ret)
}
