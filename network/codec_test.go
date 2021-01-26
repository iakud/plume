package network

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"testing"
)

type codecTest struct {
}

func (this *codecTest) Read(r io.Reader) ([]byte, error) {
	h := make([]byte, 2)
	if _, err := io.ReadFull(r, h); err != nil {
		return nil, err
	}
	n := binary.BigEndian.Uint16(h)
	b := make([]byte, n)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (this *codecTest) Write(w io.Writer, b []byte) error {
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

func TestCodec(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	var c codecTest
	message := "hello"
	if err := c.Write(buffer, []byte(message)); err != nil {
		log.Fatalln(err)
	}
	log.Println("codec write", message)
	b, err := c.Read(buffer)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("codec read", string(b))
}
