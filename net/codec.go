package net

import (
	"encoding/binary"
	"io"
)

type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

func (this *Decoder) Decode() ([]byte, error) {
	h := make([]byte, 2)
	if _, err := io.ReadFull(this.r, h); err != nil {
		return nil, err
	}
	n := binary.BigEndian.Uint16(h)
	b := make([]byte, n)
	if _, err := io.ReadFull(this.r, b); err != nil {
		return nil, err
	}
	return b, nil
}

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (this *Encoder) Encode(b []byte) error {
	h := make([]byte, 2)
	binary.BigEndian.PutUint16(h, uint16(len(b)))
	if _, err := this.w.Write(h); err != nil {
		return err
	}
	if _, err := this.w.Write(b); err != nil {
		return err
	}
	return nil
}
