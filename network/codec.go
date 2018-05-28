package network

import (
	"encoding/binary"
	"io"
)

type Codec interface {
	Read(rd io.Reader) ([]byte, error)
	Write(w io.Writer, b []byte) error
}

type StdCodec struct {
}

func (this *StdCodec) Read(rd io.Reader) ([]byte, error) {
	h := make([]byte, 2)
	if _, err := io.ReadFull(rd, h); err != nil {
		return nil, err
	}
	n := binary.BigEndian.Uint16(h)
	b := make([]byte, n)
	if _, err := io.ReadFull(rd, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (this *StdCodec) Write(w io.Writer, b []byte) error {
	h := make([]byte, 2)
	binary.BigEndian.PutUint16(h, uint16(len(b)))
	if _, err := w.Write(h); err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}
