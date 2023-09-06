package protocol

import "io"

type iPacket interface {
	Pack() ([]byte, error)
	UnPack(io.Reader) error
}
