package main

import (
	"NPProj3/Account"
	"NPProj3/BroadCast"
	"NPProj3/Chat"
	"NPProj3/Utils"
	"fmt"
	"io"
	"net"
	"syscall"
)

func Dispatcher(connect net.Conn) {
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
					break OUT
				}
			case syscall.Errno:
				if t == syscall.ECONNREFUSED {
					fmt.Printf("[-] Client %v interrupted \n", connect.RemoteAddr())
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
		Utils.ErrHandle(err)
		go Dispatcher(conn)
	}
	//TODO: 未验证身份的用户调用API的处理，即session
}
