package File

import (
	"NPProj3/ORM"
	"NPProj3/Utils"
	"NPProj3/Wigets"
	"encoding/json"
	"net"
)

type ListResponse struct {
	Uuid          string         `json:"uuid"`
	Result        string         `json:"result"`
	GroupFileList []ORM.FileInfo `json:"group_file_list"`
}

func GroupFileList(connection net.Conn, req ORM.MessageBlock) {
	fileList := []ORM.FileInfo(nil)
	for _, file := range Utils.FileManager.List {
		if file.State == Utils.FileState.Finish {
			fileList = append(fileList, file.FileInfo)
		}
	}
	retJson := ListResponse{
		Uuid:          req.Uuid,
		Result:        "success",
		GroupFileList: fileList,
	}
	retByte, err := json.Marshal(retJson)
	Wigets.ErrHandle(err)
	Wigets.SendBuf(connection, retByte)
}
