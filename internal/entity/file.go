package entity

type File struct {
	id   string
	path string
}

func NewFile(id, path string) *File {
	return &File{
		id:   id,
		path: path,
	}
}

func (f *File) Id() string {
	return f.id
}

func (f *File) Path() string {
	return f.path
}
