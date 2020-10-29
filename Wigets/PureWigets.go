package Wigets

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"
)

func ErrHandle(e interface{}) {
	if e != nil {
		log.Fatal(e)
	}
}

func RecvBuf(conn net.Conn, ch chan []byte, control chan int) {
	buf := make([]byte, 4096)
	packageLen := -1
	retStr := ""
	for {
		err := conn.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
		ErrHandle(err)
		cnt, err := conn.Read(buf)
		switch err {
		case io.EOF:
			fmt.Printf("[-] Client %v disconnected \n", conn.RemoteAddr())
			control <- -1
			return
		case nil:
			if packageLen == -1 {
				sPackageLength := string(buf[:5])
				packageLen, err = strconv.Atoi(sPackageLength)
				ErrHandle(err)
				retStr += string(buf[5:cnt])
			} else {
				if len(retStr) == packageLen {
					ch <- []byte(retStr)
					retStr = ""
					packageLen = -1
				} else if len(retStr) < packageLen {
					packageLen += cnt
					retStr += string(buf[:cnt])
				}
			}
		default:
			if errors.Is(err, os.ErrDeadlineExceeded) {
				if len(retStr) == packageLen {
					ch <- []byte(retStr)
					retStr = ""
					packageLen = -1
				}
				continue
			}
			switch t := err.(type) {
			case *net.OpError:
				if t.Op == "read" {
					fmt.Printf("[-] Client %v interrupted \n", conn.RemoteAddr())
					control <- -1
					return
				}
			case syscall.Errno:
				if t == syscall.ECONNREFUSED {
					fmt.Printf("[-] Client %v interrupted \n", conn.RemoteAddr())
					control <- -1
					return
				}
			}
		}
	}
}

//TODO: 获取buf长度，然后zero
//
//fill为5加在前面发送
func SendBuf(conn net.Conn, buf []byte) {
	bufLength := 0
	for index, key := range buf {
		if key != 0 {
			bufLength = index
		}
	}
	sBufLength := fmt.Sprintf("%05d", bufLength)
	buf = buf[:bufLength]
	buf = []byte(sBufLength + string(buf))
	_, err := conn.Write(buf)
	ErrHandle(err)
}
