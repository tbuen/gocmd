package fs

type File interface {
	Name() string
}

type file struct {
	name string
}

func newFile(name string) File {
	return &file{name: name}
}

func (f *file) Name() string {
	return f.name
}
