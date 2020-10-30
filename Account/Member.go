package Account

import (
	"NPProj3/ORM"
	"NPProj3/Utils"
	"NPProj3/Wigets"
	"encoding/json"
	"net"
)

type UserInfo struct {
	Username string `json:"username"`
}

type MemberResponse struct {
	Result string     `json:"result"`
	Uuid   string     `json:"uuid"`
	Member []UserInfo `json:"member"`
}

func GetMembers(connection net.Conn, requestJson ORM.MessageBlock) {
	uuid := requestJson.Uuid
	retJson := MemberResponse{
		Result: "success",
		Uuid:   uuid,
		Member: nil,
	}
	for key := range Utils.ConnectionMap {
		newMember := UserInfo{Username: key}
		retJson.Member = append(retJson.Member, newMember)
	}
	ret, err := json.Marshal(retJson)
	Wigets.ErrHandle(err)
	Wigets.SendBuf(connection, ret)
}
