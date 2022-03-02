package dataset

type Collection[T any] struct {
	a []*T
}

func (s *Collection[T]) Len() int {
	return len(s.a)
}

func (s *Collection[T]) Index(i int) *T {
	return s.a[i]
}

func (s *Collection[T]) GetAll() []*T {
	return s.a
}