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

func (c *Client) SendFile(filename string) error {
	fp, err := os.Open(filename)
	if err != nil {
		return err
	}

	buf := make([]byte, 1500)
	meta, err := os.Stat(filename)
	if err != nil {
		return err
	}

	d1 := filename + "_" + strconv.Itoa(int(meta.Size()))
	log.Println(d1)

	//filemeta, err := protocol.Pack([]byte(d1), protocol.FileMeta)
	p := protocol.Packet{Data: []byte(d1), Kind: protocol.FileMeta}
	filemeta, err := p.Pack()
	if err != nil {
		return err
	}

	_, err = c.conn.Write(filemeta)
	if err != nil {
		return err
	}

	for {
		n, err := fp.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("read EOF") // EOF ???
			} else {
				log.Println("some err: ", err)
			}
			break
		}

		tp := protocol.Packet{Data: buf[:n], Kind: protocol.FILE}
		r, err := tp.Pack()
		if err != nil {
			return err // todo: ??
		}

		_, err = c.conn.Write(r)
		if err != nil {
			log.Println(err)
			return err // todo: ??
		}
	}
	return nil
}

func (c *Client) SendText(text string) error {
	p := protocol.Packet{Data: []byte(text), Kind: protocol.TEXT}
	data, err := p.Pack()
	if err != nil {
		log.Println("pack error: ", err)
		return err
	}
	_, err = c.conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendJson(jsondata []byte) error {
	// jsondata produced by json.Marshal()
	// jsonData, _ := json.Marshal(p)

	p := protocol.Packet{Data: jsondata, Kind: protocol.JSON}
	data, err := p.Pack()
	_, err = c.conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}
