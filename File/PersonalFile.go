package File

import (
	"NPProj3/ORM"
	"NPProj3/Utils"
	"NPProj3/Wigets"
	"encoding/json"
	"net"
)

// 接收消息 --> 转发给客户端等待客户端的yes/no --> yes就开始转发文件块/no就下一个消息
func PersonalFileAsk(connection net.Conn, req ORM.MessageBlock) {
	sendTo := req.SendTo
	waitForConfirmJson := ORM.FileSendToClient{
		MessageType: "personal_file_ask",
		From:        req.Username,
		Uuid:        req.Uuid,
		FileInfo:    req.FileInfo,
	}
	if target := Utils.FindConnByUsername(sendTo); target != nil {
		respByte, err := json.Marshal(waitForConfirmJson)
		Wigets.ErrHandle(err)
		Wigets.SendBuf(target, respByte)

	} else {
		respJson := ORM.CommonResponse{
			Result: "No Such User",
			Uuid:   req.Uuid,
		}
		respByte, err := json.Marshal(respJson)
		Wigets.ErrHandle(err)
		Wigets.SendBuf(connection, respByte)
	}
}

func PersonalFileAskResponse(connection net.Conn, req ORM.MessageBlock) {
	confirmJson := ORM.MessageBlock{
		MessageType: req.MessageType,
		Uuid:        req.Uuid,
		Plain:       req.Plain,
		SendTo:      req.SendTo,
	}
	respByte, err := json.Marshal(confirmJson)
	Wigets.ErrHandle(err)
	if target := Utils.FindConnByUsername(req.SendTo); target != nil {
		Wigets.SendBuf(target, respByte)
	} else {
		respJson := ORM.CommonResponse{
			Result: "No Such User",
			Uuid:   req.Uuid,
		}
		respByte, err := json.Marshal(respJson)
		Wigets.ErrHandle(err)
		Wigets.SendBuf(connection, respByte)
	}
}

func PersonalFileBlockTransfer(connection net.Conn, req ORM.MessageBlock) {
	sendTo := req.SendTo
	if target := Utils.FindConnByUsername(sendTo); target != nil {
		req.Session = ""
		responseByte, err := json.Marshal(req)
		Wigets.ErrHandle(err)
		Wigets.SendBuf(target, responseByte)
	} else {
		respJson := ORM.CommonResponse{
			Result: "No Such user",
			Uuid:   req.Uuid,
		}
		respByte, err := json.Marshal(respJson)
		Wigets.ErrHandle(err)
		Wigets.SendBuf(connection, respByte)
	}
}

func PersonalAckTransfer(connection net.Conn, req ORM.MessageBlock) {
	sendTo := req.SendTo
	if target := Utils.FindConnByUsername(sendTo); target != nil {
		req.Session = ""
		responseByte, err := json.Marshal(req)
		Wigets.ErrHandle(err)
		Wigets.SendBuf(target, responseByte)
	} else {
		respJson := ORM.CommonResponse{
			Result: "No Such user",
			Uuid:   req.Uuid,
		}
		respByte, err := json.Marshal(respJson)
		Wigets.ErrHandle(err)
		Wigets.SendBuf(connection, respByte)
	}
}
