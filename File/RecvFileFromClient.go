package File

import (
	"NPProj3/ORM"
	"NPProj3/Utils"
	"encoding/json"
	"net"
)

func RecvFileMeta(connection net.Conn, req ORM.MessageBlock) {
	switch req.Plain {
	case "continue":
		if Utils.FileManager.IsUUIDExist(req.Uuid) {
			file := Utils.FileManager.FindFileItemByUUID(req.Uuid)
			err := file.WriteIn(req.Content)
			Utils.ErrHandle(err)
			respJson := ORM.CommonResponse{
				Result: "success",
				Uuid:   req.Uuid,
			}
			respRet, err := json.Marshal(respJson)
			Utils.ErrHandle(err)
			_, err = connection.Write(respRet)
			Utils.ErrHandle(err)
		}
	case "finish":
		err := Utils.FileManager.Finish(req.Uuid)
		respJson := ORM.CommonResponse{
			Result: "",
			Uuid:   req.Uuid,
		}
		if err != nil {
			Utils.ErrHandle(err)
			respJson.Result = err.Error()
		} else {
			respJson.Result = "success"
		}
		respRet, err := json.Marshal(respJson)
		Utils.ErrHandle(err)
		_, err = connection.Write(respRet)
		Utils.ErrHandle(err)
	}
}
