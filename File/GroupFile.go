package File

import (
	"NPProj3/ORM"
	"myTest/Unit"
	"net"
	"os"
)

func GroupFile(req ORM.MessageBlock, connection net.Conn) {
	fileInfo := req.FileInfo
	if !isExist("./group") {
		err := os.Mkdir("./group", 0777)
		Unit.ErrHandle(err)
	}
	file, err := os.Create("./group/" + fileInfo.MD5)
	Unit.ErrHandle(err)
	fileRecv(file)

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

func fileRecv(file *os.File) {
	Unit.ErrHandle(file) //TODO:怎么收文件
}
