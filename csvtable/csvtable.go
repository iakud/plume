package csvtable

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"reflect"

	"github.com/gocarina/gocsv"
)

var (
	csvItems  []csvItem
	csvTables map[string]CSVTable
	stdTables []*stdTable
)

func init() {
	csvTables = make(map[string]CSVTable)
}

type Item interface {
	Key() interface{}
}

type CSVTable interface {
	GetItems() interface{}
	GetItem(key interface{}) (interface{}, bool)
}

type csvTable int

func (t *csvTable) GetItems() interface{} {
	return stdTables[*t].items
}

func (t *csvTable) GetItem(key interface{}) (interface{}, bool) {
	if item, ok := stdTables[*t].itemMap[key]; ok {
		return item, true
	}
	return nil, false
}

type csvItem struct {
	name string
	item Item
}

func AddTable(name string, item Item) CSVTable {
	if t, ok := csvTables[name]; ok {
		return t
	}
	t := new(csvTable)
	*t = csvTable(len(csvItems))
	csvItems = append(csvItems, csvItem{name, item})
	csvTables[name] = t
	return t
}

type stdTable struct {
	items   interface{}
	itemMap map[interface{}]interface{}
}

func Open(dir string) error {
	tables := make([]*stdTable, len(csvItems))
	for i, csvItem := range csvItems {
		filename := filepath.Join(dir, csvItem.name) + ".csv"
		table, err := readTable(filename, csvItem.item)
		if err != nil {
			return err
		}
		tables[i] = table
	}
	stdTables = tables
	return nil
}

func readTable(filename string, out interface{}) (*stdTable, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader := csv.NewReader(f)
	um, err := gocsv.NewUnmarshaller(reader, out)
	if err != nil {
		return nil, err
	}

	itemsValue := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(out)), 0, 0)
	itemMap := make(map[interface{}]interface{})
	for {
		obj, err := um.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		item, ok := obj.(Item)
		if !ok {
			continue
		}
		itemsValue = reflect.Append(itemsValue, reflect.ValueOf(obj))
		itemMap[item.Key()] = item
	}

	table := &stdTable{itemsValue.Interface(), itemMap}
	return table, nil
}
