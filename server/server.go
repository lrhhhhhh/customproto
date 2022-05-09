package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"custom-protocol-over-tcp/protocol"
)

const (
	basePath = "/home/lrhaoo/GolandProjects/gmc/server/"
)

func handle(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		err := protocol.PeekPacketKind(reader)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("fuckyou: ", err)
			}
		}

		recv, kind, err := protocol.UnPack(reader); if err != nil {
			fmt.Println("recv err: ", err)
		}

		switch kind {
		case protocol.TEXT:
			fmt.Println(string(recv))
			break
		case protocol.JSON:
			p := &protocol.Packet{}
			err = json.Unmarshal(recv, p); if err != nil {
				fmt.Println("unmarshal fail: ", err)
			}
			fmt.Println(p.Id, p.Content)
			break
		case protocol.FILE_META:
			fileMeta := string(recv)
			arr := strings.Split(fileMeta, "_")
			filename := arr[0]
			filesize, _ := strconv.Atoi(arr[1])
			size := 0
			fmt.Println("file meta: ", filename, filesize)

			fp, err := os.OpenFile(basePath+"test.mp4", os.O_CREATE|os.O_WRONLY, 0666); if err != nil {
				fmt.Println(err)
			}

			for size < filesize {
				recv, kind, err = protocol.UnPack(reader); if err != nil {
					if err == io.EOF {
						continue
					} else {
						fmt.Println("recv file err: ", err)
					}
				}
				if kind == protocol.FILE {
					_, err = fp.Write(recv); if err != nil {
						fmt.Println("write file err: ", err)
					}
					size += len(recv)
				} else {
					fmt.Println("something err")
				}
			}
			break
		}
	}
}

func main() {
	addr := "localhost:12345"
	server, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("connect fail")
	}
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("accept err: ", err.Error())
			os.Exit(1)
		}
		go handle(conn)
	}
}
