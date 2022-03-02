package dataset

import (
	"os"

	"github.com/gocarina/gocsv"
)

func LoadCollectionCSV[T any](name string) (*Collection[T], error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var a []*T
	if err := gocsv.UnmarshalFile(f, &a); err != nil {
		return nil, err
	}
	collection := &Collection[T]{a}
	collection.a = a
	return collection, nil
}

type keyer[K comparable, T any] interface {
	Key() K
	*T
}

func LoadDictCSV[T any, K comparable, PT keyer[K, T]](name string) (*Dict[T, K], error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var a []*T
	if err := gocsv.UnmarshalFile(f, &a); err != nil {
		return nil, err
	}
	m := make(map[K]*T)
	for _, v := range a {
		m[PT(v).Key()] = v
	}
	dict := &Dict[T, K]{}
	dict.a = a
	dict.m = m
	return dict, nil
}