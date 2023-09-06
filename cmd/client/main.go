package main

import (
	"bufio"
	"customproto/client"
	"customproto/model"
	"encoding/json"
	"log"
	"os"
	"strings"
)

func main() {
	addr := "localhost:12345"
	c, err := client.New(addr)
	if err != nil {
		panic(err)
	}
	log.Println("Client is running, connect to ", addr)
	var cmd, data, line string
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		args := strings.Split(line, " ")
		cmd = args[0]
		data = args[1]
		switch cmd {
		case "text":
			if err := c.SendText(data); err != nil {
				log.Println("Send Text err:", err)
			}
		case "json":
			msg := model.Message{Content: data}
			dump, err := json.Marshal(msg)
			if err != nil {
				log.Println("json marshal err:", err)
			}
			if err := c.SendJson(dump); err != nil {
				log.Println("Send Json err:", err)
			}
		case "file":
			filename := strings.Trim(data, " \n")
			if err := c.SendFile(filename); err != nil {
				log.Println("Send file err:", err)
			}
		default:
			log.Println("invalid cmd: ", cmd)
			break
		}
	}
}
