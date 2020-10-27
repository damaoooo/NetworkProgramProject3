package main

import (
	"NPProj3/Account"
	"NPProj3/BroadCast"
	"NPProj3/Chat"
	"NPProj3/File"
	"NPProj3/Utils"
	"NPProj3/Wigets"
	"fmt"
	"io"
	"net"
	"syscall"
)

func Dispatcher(connect net.Conn) {
	if !Utils.ChapAuth(connect) {
		return
	}
	buf := make([]byte, 4096)
OUT:
	for {
		cnt, err := connect.Read(buf)
		switch err {
		case io.EOF:
			fmt.Printf("[-] Client %v disconnected \n", connect.RemoteAddr())
			break OUT
		case nil:
		default:
			switch t := err.(type) {
			case *net.OpError:
				if t.Op == "read" {
					fmt.Printf("[-] Client %v interrupted \n", connect.RemoteAddr())
					Account.InterruptQuit(connect)
					break OUT
				}
			case syscall.Errno:
				if t == syscall.ECONNREFUSED {
					fmt.Printf("[-] Client %v interrupted \n", connect.RemoteAddr())
					Account.InterruptQuit(connect)
					break OUT
				}
			}
		}
		messageType, recvJson := Utils.UnSerialize(buf[:cnt])
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
			continue
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
