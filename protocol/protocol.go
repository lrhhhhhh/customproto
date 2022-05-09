package protocol

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

const (
	INVALID uint16 = 0
	TEXT uint16 = 1
	JSON uint16 = 2
	FILE_META uint16 = 3
	FILE uint16 = 4
)

var (
	errLargeData = errors.New("data larger than 65535")
)

type Packet struct {
	Id int `json:"id"`
	Content string `json:"content"`
}


func PeekPacketKind(reader *bufio.Reader) error {
	_, err := reader.Peek(4); if err != nil {
		return err
	}
	return nil
}

func Pack(buf []byte, kind uint16) ([]byte, error) {
	length := int16(len(buf))
	data := &bytes.Buffer{}
	err := binary.Write(data, binary.LittleEndian, length); if err != nil {
		return nil, err
	}
	err = binary.Write(data, binary.LittleEndian, kind); if err != nil {
		return nil, err
	}

	err = binary.Write(data, binary.LittleEndian, buf); if err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}


func UnPack(reader *bufio.Reader) ([]byte, uint16, error) {
	header, err := reader.Peek(4); if err != nil {
		return nil, 0, err
	}
	length, _ := bytesToUint16(header[0:2])
	kind, _ := bytesToUint16(header[2:4])
	if length > 1<<16-1 {
		return nil, 0, errLargeData
	}

	packet := make([]byte, int(4+length))
	_, err = io.ReadFull(reader, packet)
	if err != nil {
		return nil, 0, err
	}

	return packet[4:], kind, nil
}


func bytesToUint16(buf []byte) (uint16, error) {
	if len(buf) != 2 {
		return 0, errors.New("len(buf) must equal 2")
	} else {
		return binary.LittleEndian.Uint16(buf), nil
	}
}
