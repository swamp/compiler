package loader

type LocalFileSystemPath string

type DocumentProvider interface {
	ReadDocument(uri LocalFileSystemPath) (string, error)
}
