package File

import (
	"NPProj3/ORM"
	"NPProj3/Utils"
	"NPProj3/Wigets"
	"encoding/json"
	"net"
)

// 接收请求 ——> 回复ACK, 找到conn, --> 另一边变成接收模式 --> 发送文件
func PersonalFile(connection net.Conn, req ORM.MessageBlock) {
	sendTo := req.SendTo
	if target := Utils.FindConnByUsername(sendTo); target != nil {
		file := Utils.FileManager.FindFileItemByMD5(req.FileInfo.MD5)
		if file != nil {
			respJson := ORM.CommonResponse{
				Result: "success",
				Uuid:   req.Uuid,
			}
			respByte, err := json.Marshal(respJson)
			Wigets.ErrHandle(err)
			_, err = connection.Write(respByte)
			Wigets.ErrHandle(err)
			targetResp := ORM.FileSendToClient{
				MessageType: "personal_file",
				From:        sendTo,
				Uuid:        req.Uuid,
				FileInfo:    req.FileInfo,
			}

			targetByte, err := json.Marshal(targetResp)
			Wigets.ErrHandle(err)
			_, err = target.Write(targetByte)
			Wigets.ErrHandle(err)

			SendFileMeta(connection, file.FileDescriptor, req)
		} else {
			respJson := ORM.CommonResponse{
				Result: "No Such File",
				Uuid:   req.Uuid,
			}
			respByte, err := json.Marshal(respJson)
			Wigets.ErrHandle(err)
			_, err = connection.Write(respByte)
			Wigets.ErrHandle(err)
		}
	} else {
		respJson := ORM.CommonResponse{
			Result: "No Such User",
			Uuid:   req.Uuid,
		}
		respByte, err := json.Marshal(respJson)
		Wigets.ErrHandle(err)
		_, err = connection.Write(respByte)
		Wigets.ErrHandle(err)
	}
}
