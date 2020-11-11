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

func RecvBuf(conn net.Conn, ch chan string, control chan int) {
	buf := make([]byte, 40960)
	packageLen := -1
	retStr := ""
	for {
		err := conn.SetReadDeadline(time.Now().Add(1 * time.Millisecond))
		ErrHandle(err)

		cnt, err := conn.Read(buf)
		switch err {
		case io.EOF:
			fmt.Printf("[-] Client %v disconnected \n", conn.RemoteAddr())
			control <- 0
			return
		case nil:
			// 解决半包和粘包问题
		RECV:
			for {
				if packageLen == -1 {
					sPackageLength := string(buf[:5])
					packageLen, err = strconv.Atoi(sPackageLength)
					ErrHandle(err)
					retStr += string(buf[5:cnt])
				}
				if len(retStr) == packageLen {
					ch <- retStr
					retStr = ""
					packageLen = -1
					break
				} else if len(retStr) > packageLen {
					// 如果粘包
					for len(retStr) > packageLen {
						ch <- retStr[:packageLen]
						retStr = retStr[packageLen:]
						packageLen, err = strconv.Atoi(retStr[:5])
						ErrHandle(err)
						retStr = retStr[5:]
					}
					break RECV
				} else if len(retStr) < packageLen {
					retStr += string(buf[:cnt])
					if len(retStr) < packageLen {
						break
					} else {
						continue
					}
				}
			}
		default:
			if errors.Is(err, os.ErrDeadlineExceeded) {
				if len(retStr) == packageLen {
					ch <- retStr
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

//fill为5加在前面发送
func SendBuf(conn net.Conn, buf []byte) {
	bufLength := 0
	for _, value := range buf {
		if value != 0 {
			bufLength++
		}
	}
	sBufLength := fmt.Sprintf("%05d", bufLength)
	buf = buf[:bufLength]
	buf = []byte(sBufLength + string(buf))
	_, err := conn.Write(buf)
	ErrHandle(err)
}
