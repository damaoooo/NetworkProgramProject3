package File

import (
	"NPProj3/ORM"
	"NPProj3/Utils"
	"NPProj3/Wigets"
	"encoding/base64"
	"encoding/json"
	"io"
	"net"
	"os"
)

func SendFileMeta(connection net.Conn, filePath string, req ORM.MessageBlock) {
	buf := make([]byte, 10240)
	resp := ORM.SendFileResponse{
		Uuid:        req.Uuid,
		MessageType: "file_send",
		Plain:       "",
		Content:     "",
	}
	confirm := "success"
	file, err := os.Open(filePath)
	Wigets.ErrHandle(err)
	ch, err := Utils.FileChanManager.FindChanByUuid(req.Uuid)
	Wigets.ErrHandle(err)
OUT:
	for {
		switch confirm {
		case "success":
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
				resp.Content = ""
				respByte, err := json.Marshal(resp)
				Wigets.ErrHandle(err)
				Wigets.SendBuf(connection, respByte)
				confirm = <-ch
				break OUT

			} else {
				resp.Plain = "continue"
				resp.Content = base64.StdEncoding.EncodeToString(buf[:cnt])
				respByte, err := json.Marshal(resp)
				Wigets.ErrHandle(err)
				Wigets.SendBuf(connection, respByte)
				confirm = <-ch
			}
		}

	}

}

func GroupFileConfirm(req ORM.MessageBlock) {
	msgType := req.Plain
	ch, err := Utils.FileChanManager.FindChanByUuid(req.Uuid)
	Wigets.ErrHandle(err)
	ch <- msgType
}
