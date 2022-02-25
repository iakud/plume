package csvdata

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/gocarina/gocsv"
)

type Table[T any] struct {
	items []*T
}

func (t *Table[T]) GetItems() []*T {
	return t.items
}

func (t *Table[T]) Load(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	reader := csv.NewReader(f)
	um, err := gocsv.NewUnmarshaller(reader, (*T)(nil))
	if err != nil {
		return err
	}
	var items []*T
	for {
		obj, err := um.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		item, ok := obj.(*T)
		if !ok {
			continue
		}
		items = append(items, item)
	}
	t.items = items
	return nil
}