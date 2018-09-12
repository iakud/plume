package ini

type Section struct {
	name string

	keys   []*Key
	keyMap map[string]*Key
}

func (this *Section) newKey(name string) *Key {
	key := newKey(name)
	return key
}

func (this *Section) Name() string {
	return this.name
}

func (this *Section) GetKey(name string) *Key {
	if key, ok := this.keyMap[name]; ok {
		return key
	}
	return nil
}

func (this *Section) Keys() []*Key {
	keys := make([]*Key, len(this.keys))
	copy(keys, this.keys)
	return keys
}
