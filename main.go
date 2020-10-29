package main

import (
	"NPProj3/Account"
	"NPProj3/BroadCast"
	"NPProj3/Chat"
	"NPProj3/File"
	"NPProj3/Utils"
	"NPProj3/Wigets"
	"fmt"
	"net"
)

func Dispatcher(connect net.Conn) {
	if !Utils.ChapAuth(connect) {
		return
	}
	recv := make(chan []byte)
	control := make(chan int)
	go Wigets.RecvBuf(connect, recv, control)
OUT:

	for {
		select {
		case buf := <-recv:
			messageType, recvJson := Utils.UnSerialize(buf)
			if !Utils.SessionValidate(*recvJson, connect) {
				continue
			}
			switch messageType {
			case "login":
				Account.Login(connect, *recvJson)
			case "offline":
				Account.Logout(connect, *recvJson)
			case "get_members":
				Account.GetMembers(connect, *recvJson)
			case "group":
				Chat.GroupChat(connect, *recvJson)
			case "person":
				Chat.PersonalChat(connect, *recvJson)
			case "group_file":
				File.GroupFileUpload(connect, *recvJson)
			case "group_file_list":
				File.GroupFileList(connect, *recvJson)
			case "file_send":
				File.RecvFileMeta(connect, *recvJson)
			case "download_group_file":
				File.GroupFileDownload(connect, *recvJson)
			case "person_file":
				File.PersonalFile(connect, *recvJson)

			case "ack":
				continue OUT
			}
		case status := <-control:
			if status == -1 {
				Account.InterruptQuit(connect)
				break OUT
			}
		}
	}
}

func main() {
	server, err := net.Listen("tcp", ":5123")
	if err != nil {
		fmt.Printf("[-] Server start Failed due to %v\n", err)
		return
	}
	fmt.Println("[+] Server start at port 5123")
	go BroadCast.MessageListen()
	fmt.Println("[+] Start Event Listening")
	for {
		conn, err := server.Accept()
		fmt.Printf("[+] Client %v connect\n", conn.RemoteAddr())
		Wigets.ErrHandle(err)
		go Dispatcher(conn)
	}

}
