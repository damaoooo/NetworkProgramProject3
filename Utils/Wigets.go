package Utils

import (
	"NPProj3/ORM"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

var ConnectionMap = make(map[string]net.Conn)
var MessageQueue = InitQueue()
var SessionM = SessionManager{
	Sessions: map[string]Session(nil),
	Lock:     sync.Mutex{},
	Err:      SessionError{},
}
var FileManager = FileList{
	Length: 0,
	List:   []FileItem(nil),
	Lock:   sync.Mutex{},
	Err:    FileErr{},
}

func ErrHandle(e interface{}) {
	if e != nil {
		log.Fatal(e)
	}
}

func UnSerialize(buf []byte) (string, *ORM.MessageBlock) {
	recvJson := new(ORM.MessageBlock)
	confirmJson := new(ORM.CommonResponse)
	err := json.Unmarshal(buf, confirmJson)
	err = json.Unmarshal(buf, recvJson)
	ErrHandle(err)
	if confirmJson.Result == "success" {
		return "ack", nil
	} else {
		return recvJson.MessageType, recvJson
	}
}

func SessionValidate(req ORM.MessageBlock, conn net.Conn) bool {
	if req.MessageType == "login" {
		return true
	} else {
		if SessionM.isValid(req.Session) {
			return true
		} else {
			failedJson := ORM.WrongSession{Info: "Go Away!"}
			failedByte, err := json.Marshal(failedJson)
			ErrHandle(err)
			_, err = conn.Write(failedByte)
			ErrHandle(err)
			return false
		}
	}
}

func FileMD5Path(path string) string {
	file, err := os.Open(path)
	ErrHandle(err)
	md5_ := md5.New()
	_, err = io.Copy(md5_, file)
	ErrHandle(err)
	md5String := hex.EncodeToString(md5_.Sum(nil))
	return md5String
}

func FileMD5FileDescriptor(file *os.File) string {
	md5_ := md5.New()
	_, err := io.Copy(md5_, file)
	ErrHandle(err)
	md5String := hex.EncodeToString(md5_.Sum(nil))
	return md5String
}
