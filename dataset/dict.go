package dataset

type Dict[T any, K comparable] struct {
	Collection[T]
	m map[K]*T
}

func (d *Dict[T, K]) Exist(key K) bool {
	_, ok := d.m[key]
	return ok
}

func (d *Dict[T, K]) MustExist(key K) {
	if _, ok := d.m[key]; !ok {
		panic("error")
	}
}

func (d *Dict[T, K]) Get(key K) (*T, bool) {
	v, ok := d.m[key]
	return v, ok
}

func (d *Dict[T, K]) MustGet(key K) *T {
	v, ok := d.m[key]
	if !ok {
		panic("error")
	}
	return v
}