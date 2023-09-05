package protocol

import "io"

type IPacket interface {
	Pack() ([]byte, error)
	UnPack(io.Reader) error
}
