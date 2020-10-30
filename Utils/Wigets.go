package Utils

import (
	"NPProj3/ORM"
	"NPProj3/Wigets"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

var FileFolder = "./"

var ConnectionMap = make(map[string]net.Conn)
var MessageQueue = InitQueue()
var sessionMaps = make(map[string]Session)
var SessionM = SessionManager{
	Sessions: sessionMaps,
	Lock:     sync.Mutex{},
	Err:      SessionError{},
}
var FileManager = FileList{
	Length: 0,
	List:   []FileItem(nil),
	Lock:   sync.Mutex{},
	Err:    FileErr{},
}

func FindConnByUsername(name string) net.Conn {
	for username, conn := range ConnectionMap {
		if username == name {
			return conn
		}
	}
	return nil
}

func UnSerialize(buf []byte) (string, *ORM.MessageBlock) {
	recvJson := new(ORM.MessageBlock)
	confirmJson := new(ORM.CommonResponse)
	err := json.Unmarshal(buf, confirmJson)
	err = json.Unmarshal(buf, recvJson)
	Wigets.ErrHandle(err)
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
			Wigets.ErrHandle(err)
			Wigets.SendBuf(conn, failedByte)
			return false
		}
	}
}

func FileMD5FileDescriptor(file *os.File) string {
	md5_ := md5.New()
	_, err := io.Copy(md5_, file)
	Wigets.ErrHandle(err)
	md5String := hex.EncodeToString(md5_.Sum(nil))
	return md5String
}

func ChapAuth(connection net.Conn) bool { //TODO: Key 可变
	key := 12345
	n := rand.Intn(10) + 2
	r := rand.New(rand.NewSource(time.Now().Unix()))
	nums := []uint32(nil)
	var res uint32 = 0
	for i := 1; i <= n; i++ {
		oneNum := r.Uint32()
		nums = append(nums, oneNum)
		res += oneNum & 0xffffffff
		res &= 0xffffffff
		r = rand.New(rand.NewSource(time.Now().Unix() + r.Int63()))
	}
	res ^= uint32(key)
	authAskJson := ORM.ChapAuthToClient{
		HowMany: uint32(n),
		Nums:    nums,
	}
	authAskByte, err := json.Marshal(authAskJson)
	Wigets.ErrHandle(err)
	Wigets.SendBuf(connection, authAskByte)
	recvBuf := make([]byte, 4096)
	cnt, err := connection.Read(recvBuf)
	recvBuf = recvBuf[:cnt]
	clientResultJson := new(ORM.AuthRecvResponse)
	err = json.Unmarshal(recvBuf, clientResultJson) //TODO: 给所有的read的反序列化操作加上切片
	Wigets.ErrHandle(err)
	responseJson := ORM.WrongSession{}
	if clientResultJson.Res == res {
		responseJson.Info = "ok"
		responseByte, err := json.Marshal(responseJson)
		Wigets.ErrHandle(err)
		Wigets.SendBuf(connection, responseByte)
		return true
	} else {
		responseJson.Info = "wrong"
		responseByte, err := json.Marshal(responseJson)
		Wigets.ErrHandle(err)
		Wigets.SendBuf(connection, responseByte)
		_ = connection.Close()
		return false
	}
}
