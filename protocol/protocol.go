package protocol

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

const (
	INVALID  uint16 = 0
	TEXT     uint16 = 1
	JSON     uint16 = 2
	FileMeta uint16 = 3
	FILE     uint16 = 4
)

var (
	errLargeData = errors.New("data larger than 65535")
)

type Packet struct {
	//Id      int    `json:"id"`
	//Content string `json:"content"`
	Kind uint16
	Len  uint16
	Data []byte
}

func PeekPacketKind(reader *bufio.Reader) error {
	_, err := reader.Peek(4)
	if err != nil {
		return err
	}
	return nil
}

func (p *Packet) Pack() ([]byte, error) {
	length := int16(len(p.Data))
	data := &bytes.Buffer{}
	if err := binary.Write(data, binary.LittleEndian, length); err != nil {
		return nil, err
	}
	if err := binary.Write(data, binary.LittleEndian, p.Kind); err != nil {
		return nil, err
	}
	if err := binary.Write(data, binary.LittleEndian, p.Data); err != nil {
		return nil, err
	}
	return data.Bytes(), nil
}

func (p *Packet) UnPack(rd io.Reader) error {
	reader := bufio.NewReader(rd)
	header, err := reader.Peek(4)
	if err != nil {
		return err
	}
	length, _ := bytesToUint16(header[0:2])
	kind, _ := bytesToUint16(header[2:4])
	if length > 1<<16-1 {
		return errLargeData
	}

	packet := make([]byte, int(4+length))
	_, err = io.ReadFull(reader, packet)
	if err != nil {
		return err
	}

	p.Data = packet[4:]
	p.Kind = kind
	return nil
}

func bytesToUint16(buf []byte) (uint16, error) {
	if len(buf) != 2 {
		return 0, errors.New("len(buf) must equal 2")
	} else {
		return binary.LittleEndian.Uint16(buf), nil
	}
}
