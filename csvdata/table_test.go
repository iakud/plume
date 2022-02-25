package csvdata

import (
	"testing"
)

type sampleKey struct {
	Id1 int
	Id2 int
}

type sample struct {
	Id1 int `csv:"id1"`
	Id2 int `csv:"id2"`
	FieldA string `csv:"field_a"`
	FieldB string `csv:"field_b"`
}

func (s *sample) Key() sampleKey {
	return sampleKey{s.Id1, s.Id2}
}

func TestCTable(t *testing.T) {
	ct := &KVTable[sample, sampleKey]{}
	if err := ct.Load("./sample.csv"); err != nil {
		t.Fatal(err)
	}
	// all items
	t.Log("------all items--------")
	for _, item := range ct.GetItems() {
		t.Log(item)
	}
	// get
	t.Log("-----------------------")
	if item, ok := ct.GetItem(sampleKey{1, 1}); ok {
		t.Log("get item:", item)
	}
	if item, ok := ct.GetItem(sampleKey{2, 1}); ok {
		t.Log("get item:", item)
	}
	if item, ok := ct.GetItem(sampleKey{4, 1}); ok {
		t.Log("get item:", item)
	}
}
