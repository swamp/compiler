package loader

type LocalFileSystemPath string

type LocalFileSystemRoot string

type DocumentProvider interface {
	ReadDocument(uri LocalFileSystemPath) (string, error)
}
