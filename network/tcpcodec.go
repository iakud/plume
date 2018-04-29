package network

import (
	"encoding/binary"
	"io"
)

const (
	defaultBufSize = 4096
)

type Codec interface {
	Read(rd io.Reader) ([]byte, error)
	Write(w io.Writer, p []byte) error
}

type defaultCodec struct {
}

func (this *defaultCodec) Read(rd io.Reader) ([]byte, error) {
	p := make([]byte, defaultBufSize)
	n, err := rd.Read(p)
	if err != nil {
		return nil, err
	}
	return p[:n], nil
}

func (this *defaultCodec) Write(w io.Writer, p []byte) error {
	if _, err := w.Write(p); err != nil {
		return err
	}
	return nil
}

type TCPCodec struct {
}

func (this *TCPCodec) Read(rd io.Reader) ([]byte, error) {
	h := make([]byte, 2)
	if _, err := io.ReadFull(rd, h); err != nil {
		return nil, err
	}
	n := binary.BigEndian.Uint16(h)
	p := make([]byte, n)
	if _, err := io.ReadFull(rd, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (this *TCPCodec) Write(w io.Writer, p []byte) error {
	h := make([]byte, 2)
	binary.BigEndian.PutUint16(h, uint16(len(p)))
	if _, err := w.Write(h); err != nil {
		return err
	}
	if _, err := w.Write(p); err != nil {
		return err
	}
	return nil
}
