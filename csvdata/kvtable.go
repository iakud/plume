package csvdata

type Keyer[T any, K comparable] interface {
	Key() K
	*T
}

type KVTable[T any, K comparable, P Keyer[T, K]] struct {
	Table[T]
	itemMap map[K]*T
}

func (t *KVTable[T, K, P]) Load(filename string) error {
	if err := t.Table.Load(filename); err != nil {
		return err
	}
	itemMap := make(map[K]*T)
	for _, item := range t.items {
		itemMap[P(item).Key()] = item
	}
	t.itemMap = itemMap
	return nil
}

func (t *KVTable[T, K, P]) GetItem(key K) (*T, bool) {
	item, ok := t.itemMap[key]
	return item, ok
}