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
	TEXT    uint16 = 1
	JSON    uint16 = 2
	FILE    uint16 = 3
)

var (
	errPacketSize = errors.New("packet Len larger than 65535")
)

type Packet struct {
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
	if len(p.Data) >= 1<<16 {
		return nil, errPacketSize
	}
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
	length := bytesToUint16(header[0:2])
	kind := bytesToUint16(header[2:4])

	packet := make([]byte, int(4+length))
	if _, err = io.ReadFull(reader, packet); err != nil {
		return err
	}

	p.Data = packet[4:]
	p.Kind = kind

	return nil
}

func bytesToUint16(buf []byte) uint16 {
	// len(buf) must equal 2
	return binary.LittleEndian.Uint16(buf)
}
