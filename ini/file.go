package ini

type File struct {
	sections   []*Section
	sectionMap map[string]*Section
}
