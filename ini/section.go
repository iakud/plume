package ini

type Section struct {
	name string

	keys   []*Key
	keyMap map[string]*Key
}

func newSection(name string) *Section {
	section := new(Section)
	section.name = name
	section.keyMap = make(map[string]*Key)
	return section
}

func (this *Section) addKey(name, value string) *Key {
	if key, ok := this.keyMap[name]; ok {
		key.value = value
		return key
	}
	key := newKey(name, value)
	this.keys = append(this.keys, key)
	this.keyMap[name] = key
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
