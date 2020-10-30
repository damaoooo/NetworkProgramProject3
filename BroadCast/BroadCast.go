package BroadCast

import (
	"NPProj3/ORM"
	"NPProj3/Utils"
	"NPProj3/Wigets"
	"encoding/json"
	"github.com/satori/go.uuid"
)

func BroadCast(userName string, info []byte) {
	for username, conn := range Utils.ConnectionMap {
		if username != userName {
			Wigets.SendBuf(conn, info)
		}
	}
}

func EventBroadcast() {
	if !Utils.MessageQueue.IsEmpty() {
		for event := Utils.MessageQueue.Dequeue(); event != nil; event = Utils.MessageQueue.Dequeue() {
			eventUuid := uuid.Must(uuid.NewV4()).String()
			msgType := event.Type
			retJson := ORM.EventResponse{
				Uuid:        eventUuid,
				MessageType: "",
				Event:       *event,
			}
			if msgType == "group" {
				retJson.MessageType = "group_message"
			} else if msgType == "online" || msgType == "offline" {
				retJson.MessageType = "event"
			}
			ret, err := json.Marshal(retJson)
			Wigets.ErrHandle(err)
			username := event.User
			BroadCast(username, ret)
		}
	}
}

func MessageListen() {
	for {
		if !Utils.MessageQueue.IsEmpty() {
			EventBroadcast()
		}
	}
}
