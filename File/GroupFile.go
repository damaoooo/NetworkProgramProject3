package File

import (
	"NPProj3/ORM"
	"NPProj3/Utils"
	"NPProj3/Wigets"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func GroupFileUpload(connection net.Conn, req ORM.MessageBlock) {
	groupFilePath := filepath.Join(Utils.FileFolder, "./group")
	fileInfo := req.FileInfo
	if !isExist(groupFilePath) {
		err := os.Mkdir(groupFilePath, 0777)
		Wigets.ErrHandle(err)
	}
	thisFilePath := filepath.Join(groupFilePath, fileInfo.Name)
	file, err := os.Create(thisFilePath)
	Wigets.ErrHandle(err)
	retJson := ORM.CommonResponse{
		Result: "success",
		Uuid:   req.Uuid,
	}
	err = Utils.FileManager.AddFile(req.Uuid, file, fileInfo, thisFilePath)
	if err != nil && err.Error() == "duplicate file" {
		retJson.Result = "duplicate_file"
	} else {
		Wigets.ErrHandle(err)
	}
	ret, err := json.Marshal(retJson)
	Wigets.ErrHandle(err)
	Wigets.SendBuf(connection, ret)
}

func GroupFileDownload(connection net.Conn, req ORM.MessageBlock) {
	fileMD5 := strings.ToLower(req.FileInfo.MD5)
	respJson := ORM.CommonResponse{
		Result: "",
		Uuid:   req.Uuid,
	}
	if fileHandle := Utils.FileManager.FindFileItemByMD5(fileMD5); fileHandle != nil {
		respJson.Result = "success"
		respByte, err := json.Marshal(respJson)
		Wigets.ErrHandle(err)
		Wigets.SendBuf(connection, respByte)
		ch := make(chan string)
		Utils.FileChanManager.Add(ch, req.Uuid)
		go SendFileMeta(connection, fileHandle.Path, req)
	} else {
		respJson.Result = "no_such_file"
		respByte, err := json.Marshal(respJson)
		Wigets.ErrHandle(err)
		Wigets.SendBuf(connection, respByte)
	}

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
