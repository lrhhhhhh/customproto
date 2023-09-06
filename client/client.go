package client

import (
	"customproto/protocol"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

type Client struct {
	addr string
	conn net.Conn
}

func New(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{addr: addr, conn: conn}, nil
}

// SendFile 先发送文件的meta信息（大小，名字），然后再发送文件内容
func (c *Client) SendFile(filename string) error {
	meta, err := os.Stat(filename)
	if err != nil {
		return err
	}

	fileMetaStr := meta.Name() + "_" + strconv.Itoa(int(meta.Size()))
	log.Printf("Send file, filename=%s, size=%d\n", meta.Name(), meta.Size())

	p := protocol.Packet{Data: []byte(fileMetaStr), Kind: protocol.FILE}
	fileMetaData, err := p.Pack()
	if err != nil {
		return err
	}

	_, err = c.conn.Write(fileMetaData)
	if err != nil {
		return err
	}

	buf := make([]byte, 1500)
	fp, err := os.Open(filename)
	if err != nil {
		return err
	}

	var size int
	for {
		n, err := fp.Read(buf)
		size += n
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err // todo: unknown error
			}
		}

		filePacket := protocol.Packet{Data: buf[:n], Kind: protocol.FILE}
		r, err := filePacket.Pack()
		if err != nil {
			return err
		}

		_, err = c.conn.Write(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) SendText(text string) error {
	p := protocol.Packet{Data: []byte(text), Kind: protocol.TEXT}
	data, err := p.Pack()
	if err != nil {
		return err
	}
	_, err = c.conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendJson(jsonData []byte) error {
	p := protocol.Packet{Data: jsonData, Kind: protocol.JSON}
	data, err := p.Pack()
	_, err = c.conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}
