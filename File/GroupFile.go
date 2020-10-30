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
	file, err := os.Create(filepath.Join(groupFilePath, fileInfo.Name))
	Wigets.ErrHandle(err)
	err = Utils.FileManager.AddFile(req.Uuid, file, fileInfo)
	Wigets.ErrHandle(err)
	retJson := ORM.CommonResponse{
		Result: "success",
		Uuid:   req.Uuid,
	}
	ret, err := json.Marshal(retJson)
	Wigets.ErrHandle(err)
	Wigets.SendBuf(connection, ret)
}

func GroupFileDownload(connection net.Conn, req ORM.MessageBlock) {
	fileMD5 := strings.ToLower(req.FileInfo.MD5)
	if fileHandle := Utils.FileManager.FindFileItemByMD5(fileMD5); fileHandle != nil {
		go SendFileMeta(connection, fileHandle.FileDescriptor, req)
	} else {
		respJson := ORM.CommonResponse{
			Result: "No Such File",
			Uuid:   req.Uuid,
		}
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
