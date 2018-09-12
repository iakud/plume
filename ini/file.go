package ini

type File struct {
	sections   []*Section
	sectionMap map[string]*Section
}

func (this *File) addSection(name string) *Section {
	if section, ok := this.sectionMap[name]; ok {
		return section
	}
	section := newSection(name)
	this.sections = append(this.sections, section)
	this.sectionMap[name] = section
	return nil
}

func (this *File) GetSection(name string) *Section {
	if section, ok := this.sectionMap[name]; ok {
		return section
	}

	return nil
}
