package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"sync"
	"custom-protocol-over-tcp/protocol"
)

var data1 = "response a short message"
var data2 = "response a very long text: "

func init() {
	for i:=0; i<4096; i++ {
		data2 += string(byte('a'+i%26))
	}
}

func sendFile(conn net.Conn, filename string) {
	fp, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	buf := make([]byte, 1500)
	meta, err := os.Stat(filename); if err != nil {
		fmt.Println(err)
		return
	}
	d1 := filename + "_" + strconv.Itoa(int(meta.Size()))
	fmt.Println(d1)

	filemeta, err := protocol.Pack([]byte(d1), protocol.FILE_META)
	_, err = conn.Write(filemeta); if err != nil {
		fmt.Println(err)
	}

	for {
		n, err := fp.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("read EOF")   // EOF ???
			} else {
				fmt.Println("some err: ", err)
			}
			break
		}
		r, err := protocol.Pack(buf[:n], protocol.FILE)
		_, err = conn.Write(r)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func consume(addr string, wg *sync.WaitGroup, id int) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("dail fail: ", err)
	}
	defer conn.Close()

	//// short text
	//for i:=0; i<10; i++ {
	//	data, err := protocol.Pack([]byte("hello, motherfucker! " + strconv.Itoa(id)), protocol.TEXT)
	//	if err != nil {
	//		fmt.Println("send text fail: ", err)
	//	}
	//	_, err = conn.Write(data)
	//	time.Sleep(time.Millisecond * 500)
	//}
	//
	//// short json
	//for i:=0; i<10; i++ {
	//	p := protocol.Packet{Id: id, Content: data1}
	//	jsonData, _ := json.Marshal(p)
	//	sendData, err := protocol.Pack(jsonData, protocol.JSON)
	//	_, err = conn.Write(sendData); if err != nil {
	//		fmt.Println("send fail ", err)
	//	}
	//}

	// file
	sendFile(conn, "/home/lrhaoo/GolandProjects/gmc/client/test.mp4")

	wg.Done()
}

func main() {
	addr := "localhost:12345"
	wg := sync.WaitGroup{}
	n := 1
	wg.Add(n)
	for i:=0; i<n; i++ {
		go consume(addr, &wg, i)
	}

	wg.Wait()
}