package File

import (
	"NPProj3/ORM"
	"NPProj3/Wigets"
	"encoding/base64"
	"encoding/json"
	"io"
	"net"
	"os"
)

func SendFileMeta(connection net.Conn, file *os.File, req ORM.MessageBlock) {
	// TODO: 加上base64 检查获取文件目录
	buf := make([]byte, 1024)
	resp := ORM.SendFileResponse{
		Uuid:        req.Uuid,
		MessageType: "file_send",
		Plain:       "",
		Content:     nil,
	}
	for {
		cnt, err := file.Read(buf)
		buf = buf[:cnt]
		if err != nil && err != io.EOF {
			resp.Plain = "server error"
			respByte, err := json.Marshal(resp)
			Wigets.ErrHandle(err)
			Wigets.SendBuf(connection, respByte)
			break

		} else if cnt == 0 {
			resp.Plain = "finish"
			respByte, err := json.Marshal(resp)
			Wigets.ErrHandle(err)
			Wigets.SendBuf(connection, respByte)
			break

		} else {
			resp.Plain = "continue"
			resp.Content = base64.StdEncoding.EncodeToString(buf[:cnt])
			respByte, err := json.Marshal(resp)
			Wigets.ErrHandle(err)
			Wigets.SendBuf(connection, respByte)
		}

	}

}
