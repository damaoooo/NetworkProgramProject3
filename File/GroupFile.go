package File

import (
	"NPProj3/ORM"
	"NPProj3/Utils"
	"encoding/json"
	"myTest/Unit"
	"net"
	"os"
)

func GroupFile(connection net.Conn, req ORM.MessageBlock) {
	fileInfo := req.FileInfo
	if !isExist("./group") {
		err := os.Mkdir("./group", 0777)
		Unit.ErrHandle(err)
	}
	file, err := os.Create("./group/" + fileInfo.Name)
	Utils.ErrHandle(err)
	err = Utils.FileManager.AddFile(req.Uuid, file, fileInfo)
	Unit.ErrHandle(err)
	retJson := ORM.CommonResponse{
		Result: "success",
		Uuid:   req.Uuid,
	}
	ret, err := json.Marshal(retJson)
	Unit.ErrHandle(err)
	_, err = connection.Write(ret)
	Unit.ErrHandle(err)
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
