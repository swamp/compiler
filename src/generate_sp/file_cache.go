package generate_sp

type FileUrlID uint

type FileUrl struct {
	File string
	ID   FileUrlID
}

type FileUrlCache struct {
	fileUrl map[string]*FileUrl
}

func NewFileUrlCache() *FileUrlCache {
	return &FileUrlCache{fileUrl: make(map[string]*FileUrl)}
}

func (f *FileUrlCache) GetID(fileUrl string) FileUrlID {
	existingFileUrl, hasIt := f.fileUrl[fileUrl]
	if hasIt {
		return existingFileUrl.ID
	}

	newFileUrl := &FileUrl{
		File: fileUrl,
		ID:   FileUrlID(len(f.fileUrl)),
	}
	f.fileUrl[fileUrl] = newFileUrl

	return newFileUrl.ID
}
