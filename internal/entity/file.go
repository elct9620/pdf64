package entity

type File struct {
	id          string
	path        string
	isEncrypted bool
}

func NewFile(id, path string) *File {
	return &File{
		id:          id,
		path:        path,
		isEncrypted: false,
	}
}

func (f *File) Id() string {
	return f.id
}

func (f *File) Path() string {
	return f.path
}

func (f *File) IsEncrypted() bool {
	return f.isEncrypted
}

func (f *File) Encrypt() {
	f.isEncrypted = true
}

func (f *File) Decrypt() {
	f.isEncrypted = false
}
