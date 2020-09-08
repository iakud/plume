package csvtable

import (
	"testing"
)

type sample struct {
	Id     int    `csv:"id"`
	FieldA string `csv:"field_a"`
	FieldB string `csv:"field_b"`
}

func (s *sample) Key() interface{} {
	return s.Id
}

func TestCSVTable(t *testing.T) {
	table := AddTable("test", (*sample)(nil))
	if err := Open("./"); err != nil {
		t.Fatal(err)
	}
	t.Log(table.GetItems().([]*sample)[0])
	item, ok := table.GetItem(2)
	if !ok {
		t.Fatal("get item failed")
	}
	t.Log(item.(*sample))
}
