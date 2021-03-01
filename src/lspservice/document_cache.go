package lspservice

import (
	"fmt"

	"github.com/swamp/compiler/src/loader"
)

type InMemoryDocument struct {
	payload string
}

func (c *InMemoryDocument) Overwrite(payload string) {
	c.payload = payload
}

type DocumentCache struct {
	documents        map[LocalFileSystemPath]*InMemoryDocument
	fallbackProvider loader.DocumentProvider
}

type LocalFileSystemPath string

func NewDocumentCache(fallbackProvider loader.DocumentProvider) *DocumentCache {
	return &DocumentCache{documents: make(map[LocalFileSystemPath]*InMemoryDocument), fallbackProvider: fallbackProvider}
}

func (d *DocumentCache) Open(path LocalFileSystemPath, payload string) error {
	found, _ := d.documents[path]
	if found != nil {
		found.Overwrite(payload)
	} else {
		d.documents[path] = &InMemoryDocument{payload: payload}
	}
	return nil
}

func (d *DocumentCache) Close(path LocalFileSystemPath) error {
	found, _ := d.documents[path]
	if found != nil {
		return fmt.Errorf("no such file cached and open %v", path)
	}

	delete(d.documents, path)

	return nil
}

func (d *DocumentCache) Get(path LocalFileSystemPath) (*InMemoryDocument, error) {
	found, _ := d.documents[path]
	if found != nil {
		return nil, fmt.Errorf("no such file cached and open %v", path)
	}

	return found, nil
}

func (d *DocumentCache) ReadDocument(path loader.LocalFileSystemPath) (string, error) {
	inMemoryDocument, err := d.Get(LocalFileSystemPath(path))
	if err != nil {
		return d.fallbackProvider.ReadDocument(path)
	}

	return inMemoryDocument.payload, nil
}
