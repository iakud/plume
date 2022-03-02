package dataset

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

func TestDict(t *testing.T) {
	dataSet, err := LoadDictCSV[sample, sampleKey]("./data.csv")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("count:", dataSet.Len())
	// get
	if item, ok := dataSet.Get(sampleKey{1, 1}); ok {
		t.Log("get <1,1>:", item)
	}
	if item, ok := dataSet.Get(sampleKey{2, 1}); ok {
		t.Log("get <2,1>:", item)
	}
	if item, ok := dataSet.Get(sampleKey{4, 1}); ok {
		t.Log("get <4,1>:", item)
	}
}

func TestCollection(t *testing.T) {
	dataSet, err := LoadCollectionCSV[sample]("./data.csv")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("count:", dataSet.Len())
	// get all
	for i, v := range dataSet.GetAll() {
		t.Log(i+1, ":", v)
	}
}
