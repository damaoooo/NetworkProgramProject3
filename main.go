package main

import (
	"NPProj3/Account"
	"NPProj3/ORM"
	"NPProj3/Utils"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

func Dispatcher(connect net.Conn) {
	buf := make([]byte, 4096)
	for {
		cnt, err := connect.Read(buf)
		if err != nil || cnt == 0 {
			if err == io.EOF {
				fmt.Printf("[-] Client %v disconnected \n", connect.RemoteAddr())
				break
			}
			Utils.Errhandle(err)
			err = connect.Close()
			Utils.Errhandle(err)
			break
		}
		recvJson := new(ORM.MessageBlock)
		err = json.Unmarshal(buf[:cnt], recvJson)
		Utils.Errhandle(err)
		messageType := recvJson.MessageType
		switch messageType {
		case "login":
			Account.Login(connect, *recvJson)
		case "offline":
			Account.Logout(connect, *recvJson)
		case "get_members":
			Account.GetMembers(connect, *recvJson)
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
	for {
		conn, err := server.Accept()
		fmt.Printf("[+] Client %v connect\n", conn.RemoteAddr())
		Utils.Errhandle(err)
		go Dispatcher(conn)
		// TODO: 怎么处理全局广播?
	}
}
