package File

import (
	"NPProj3/ORM"
	"NPProj3/Utils"
	"encoding/json"
	"myTest/Unit"
	"net"
)

func RecvFileMeta(connection net.Conn, req ORM.MessageBlock) {
	switch req.Plain {
	case "continue":
		if Utils.FileManager.IsUUIDExist(req.Uuid) {
			file := Utils.FileManager.FindFileItemByUUID(req.Uuid)
			err := file.WriteIn(req.Content)
			Unit.ErrHandle(err)
			respJson := ORM.CommonResponse{
				Result: "success",
				Uuid:   req.Uuid,
			}
			respRet, err := json.Marshal(respJson)
			Unit.ErrHandle(err)
			_, err = connection.Write(respRet)
			Unit.ErrHandle(err)
		}
	}
}
