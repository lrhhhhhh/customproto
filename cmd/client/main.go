package main

import (
	"bufio"
	"customproto/client"
	"log"
	"os"
	"strings"
)

//var data1 = "response a short message"
//var data2 = "response a very long text: "
//
//func init() {
//	for i := 0; i < 4096; i++ {
//		data2 += string(byte('a' + i%26))
//	}
//}

func testLargeSizeText() {

}

func main() {
	addr := "localhost:12345"
	c, err := client.New(addr)
	if err != nil {
		panic(err)
	}

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
			if err := c.SendJson([]byte{}); err != nil {
				log.Println("Send Json err:", err)
			}
		case "file":
			// filename = "/home/lrhaoo/Desktop/test.mp4"
			if err := c.SendFile(data); err != nil {
				log.Println("Send file err:", err)
			}
		default:
			log.Println("invalid cmd: ", cmd)
			break
		}
	}
}
