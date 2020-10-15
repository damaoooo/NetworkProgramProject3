package main

import (
	"fmt"
	"myTest/Unit"
	"net"
)

func Dispatcher(connect net.Conn) {

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
		Unit.ErrHandle(err)
		go Dispatcher(conn)
	}
}
