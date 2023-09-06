package main

import (
	"bufio"
	"customproto/model"
	"customproto/protocol"
	"customproto/server"
	"encoding/json"
	"errors"
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
		case protocol.JSON:
			var msg model.Message
			err = json.Unmarshal(p.Data, &msg)
			if err != nil {
				log.Println("unmarshal err:", err)
			}
			log.Printf("%+v\n", msg)
		case protocol.FILE:
			filemeta := string(p.Data)
			arr := strings.Split(filemeta, "_")
			filename := arr[0]
			filesize, _ := strconv.Atoi(arr[1])
			saveTo := UploadFolder + filename
			log.Printf("recv file, filename=%s, size=%d, saveTo=%s\n", filename, filesize, saveTo)

			fp, err := os.OpenFile(saveTo, os.O_CREATE|os.O_WRONLY, 0777)
			if err != nil {
				log.Println(err)
				continue
			}

			// 处理被分块的文件
			var size int
			filePacket := new(protocol.Packet)
			for size = 0; size < filesize; {
				err = filePacket.UnPack(reader)
				if err != nil {
					if err == io.EOF {
						log.Println(err)
						continue
					} else {
						log.Println("recv file err: ", err)
					}
				}
				if filePacket.Kind == protocol.FILE {
					if n, err := fp.Write(filePacket.Data); err != nil || n != len(filePacket.Data) {
						log.Println("write file err: ", err)
					}
					size += len(filePacket.Data)
				} else {
					panic(errors.New("invalid packet kind"))
				}
			}
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
