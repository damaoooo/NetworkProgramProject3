package File

import (
	"NPProj3/ORM"
	"NPProj3/Utils"
	"NPProj3/Wigets"
	"encoding/json"
	"net"
)

func RecvFileMeta(connection net.Conn, req ORM.MessageBlock) {
	switch req.Plain {
	case "continue":
		if Utils.FileManager.IsUUIDExist(req.Uuid) {
			file := Utils.FileManager.FindFileItemByUUID(req.Uuid)
			err := file.WriteIn(req.Content)
			Wigets.ErrHandle(err)
			respJson := ORM.CommonResponse{
				Result: "success",
				Uuid:   req.Uuid,
			}
			respRet, err := json.Marshal(respJson)
			Wigets.ErrHandle(err)
			Wigets.SendBuf(connection, respRet)
		}
	case "finish":
		err := Utils.FileManager.Finish(req.Uuid)
		respJson := ORM.CommonResponse{
			Result: "",
			Uuid:   req.Uuid,
		}
		if err != nil {
			Wigets.ErrHandle(err)
			respJson.Result = err.Error()
		} else {
			respJson.Result = "success"
		}
		respRet, err := json.Marshal(respJson)
		Wigets.ErrHandle(err)
		Wigets.SendBuf(connection, respRet)
	}
}
