package dirdrill

type FSObject interface {
	GetName() string
	GetSize() int64
	IsDirectory() bool
}

type FSFile struct {
	name string
	size int64
}

func (f *FSFile) GetName() string {
	return f.name
}

func (f *FSFile) GetSize() int64 {
	return f.size
}

func (f *FSFile) IsDirectory() bool {
	return false
}

type FSDirectory struct {
	name string
	size int64
	filesCount int64

	files []FSObject
}

func (f *FSDirectory) GetName() string {
	return f.name
}

func (f *FSDirectory) GetSize() int64 {
	return f.size
}

func (f *FSDirectory) GetFilesCount() int64 {
	return f.filesCount
}

func (f *FSDirectory) IsDirectory() bool {
	return true
}

func (f *FSDirectory) AddObject(fsObject FSObject) {
	f.files = append(f.files, fsObject)
	if !fsObject.IsDirectory() {
		fileSize := fsObject.GetSize()
		f.AddFileSize(fileSize)
	}
}

func (f *FSDirectory) AddFileSize(fileSize int64) {
	f.size += fileSize
	f.filesCount += 1
}