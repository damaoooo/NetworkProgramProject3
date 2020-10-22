package Account

import (
	"NPProj3/ORM"
	"NPProj3/Utils"
	"encoding/json"
	"log"
	"net"
)

func Login(connection net.Conn, request ORM.MessageBlock) {
	username := request.Username
	uuid := request.Uuid
	retJson := ORM.LoginResponse{Uuid: uuid}
	if _, ok := Utils.ConnectionMap[username]; ok {
		retJson.Result = "multi-username"
	} else {
		retJson.Result = "Success"
		Utils.ConnectionMap[username] = connection
		event := ORM.Event{
			Type: "online",
			User: username,
			Case: "online",
		}
		Utils.MessageQueue.Add(event)
		retJson.Session = Utils.SessionM.GenerateNew(username)
	}
	ret, err := json.Marshal(retJson)
	Utils.ErrHandle(err)
	_, err = connection.Write(ret)
	Utils.ErrHandle(err)
}

func Logout(connection net.Conn, request ORM.MessageBlock) {
	username := request.Username
	uuid := request.Uuid
	retJson := ORM.CommonResponse{Uuid: uuid}

	if _, ok := Utils.ConnectionMap[username]; ok {
		retJson.Result = "success"
		delete(Utils.ConnectionMap, username)
		event := ORM.Event{
			Type: "offline",
			User: username,
			Case: "offline",
		}
		Utils.MessageQueue.Add(event)
		err := Utils.SessionM.Destroy(request.Session)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		retJson.Result = "no-user"
	}

	ret, err := json.Marshal(retJson)
	Utils.ErrHandle(err)
	_, err = connection.Write(ret)
	Utils.ErrHandle(err)
}
