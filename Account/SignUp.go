package Account

import (
	"NPProj3/ORM"
	"NPProj3/Utils"
	"NPProj3/Wigets"
	"encoding/json"
	"net"
)

func NewUserSignUp(connection net.Conn, req ORM.MessageBlock) {
	username := req.Username
	password := req.Plain
	respJson := ORM.MessageBlock{
		MessageType: "sign_up_back",
		Uuid:        req.Uuid,
		Plain:       "",
	}
	if Utils.UserPassManager.IsExist(username) {
		respJson.Plain = "multi-user"
	} else {
		Utils.UserPassManager.AddUser(username, password)
		respJson.Plain = "success"
	}
	respByte, err := json.Marshal(respJson)
	Wigets.ErrHandle(err)
	Wigets.SendBuf(connection, respByte)
}
