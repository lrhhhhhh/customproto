package main

import (
	"bufio"
	"customproto/model"
	"customproto/protocol"
	"customproto/server"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	UploadFolder = "/home/lrhaoo/Desktop/"
)

func handle(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		err := protocol.PeekPacketKind(reader)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Println("unknown err:", err)
			}
		}

		p := new(protocol.Packet)
		if err := p.UnPack(reader); err != nil {
			log.Println("recv err: ", err)
		}

		switch p.Kind {
		case protocol.TEXT:
			log.Println(string(p.Data))
			break
		case protocol.JSON:
			var msg model.Message
			err = json.Unmarshal(p.Data, &msg)
			if err != nil {
				log.Println("unmarshal err:", err)
			}
			log.Printf("%+v\n", msg)
			break
		case protocol.FileMeta:
			filemeta := string(p.Data)
			arr := strings.Split(filemeta, "_")
			filename := arr[0]
			filesize, _ := strconv.Atoi(arr[1])
			log.Println("file meta: ", filename, filesize)

			fp, err := os.OpenFile(UploadFolder+"test.mp4", os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				log.Println(err)
				break
			}

			for size := 0; size < filesize; {
				tp := new(protocol.Packet)
				err = tp.UnPack(reader)
				if err != nil {
					if err == io.EOF {
						continue
					} else {
						fmt.Println("recv file err: ", err)
					}
				}
				if p.Kind == protocol.FILE {
					_, err = fp.Write(p.Data)
					if err != nil {
						log.Println("write file err: ", err)
					}
					size += len(p.Data)
				} else {
					panic(errors.New("invalid packet kind"))
				}
			}
			break
		}
	}
}

func main() {
	addr := "localhost:12345"
	svr, err := server.New(addr, handle)
	if err != nil {
		panic(err)
	}
	svr.Run()
}
