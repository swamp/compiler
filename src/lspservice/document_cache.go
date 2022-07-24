package lspservice

import (
	"fmt"
	"log"

	"github.com/swamp/compiler/src/loader"
)

type DocumentCache struct {
	documents        map[LocalFileSystemPath]*InMemoryDocument
	fallbackProvider loader.DocumentProvider
}

type LocalFileSystemPath string

func NewDocumentCache(fallbackProvider loader.DocumentProvider) *DocumentCache {
	if fallbackProvider == nil {
		panic("must have fallback provder")
	}
	return &DocumentCache{documents: make(map[LocalFileSystemPath]*InMemoryDocument), fallbackProvider: fallbackProvider}
}

func (d *DocumentCache) Open(path LocalFileSystemPath, payload string) error {
	log.Printf("documentcache: open %v", path)
	return d.internalOpen(path, payload)
}

func (d *DocumentCache) internalOpen(path LocalFileSystemPath, payload string) error {
	found := d.documents[path]
	if found != nil {
		found.Overwrite(payload)
	} else {
		d.documents[path] = NewInMemoryDocument(payload)
	}
	return nil
}

func (d *DocumentCache) Close(path LocalFileSystemPath) error {
	found := d.documents[path]
	if found != nil {
		return fmt.Errorf("no such file cached and open %v", path)
	}

	delete(d.documents, path)

	return nil
}

func (d *DocumentCache) Get(path LocalFileSystemPath) (*InMemoryDocument, error) {
	found := d.documents[path]
	if found == nil {
		return nil, fmt.Errorf("no such file cached and open %v\n%v", path, d.documents)
	}

	return found, nil
}

func (d *DocumentCache) GetDocumentByVersion(path LocalFileSystemPath, version DocumentVersion) (*InMemoryDocument, error) {
	existingDocument, err := d.Get(path)
	if err != nil {
		return nil, err
	}

	if existingDocument.version != version {
		return nil, fmt.Errorf("wrong version of document, %v vs %v", version, existingDocument.version)
	}

	return existingDocument, nil
}

func (d *DocumentCache) TrackDocument(path LocalFileSystemPath, payload string) error {
	log.Printf("documentcache: StartTracking %v", path)
	return d.internalOpen(path, payload)
}

func (d *DocumentCache) ReadDocument(path loader.LocalFileSystemPath) (string, error) {
	log.Printf("documentcache: ReadDocument %v", path)
	inMemoryDocument, err := d.Get(LocalFileSystemPath(path))
	if err != nil {
		payload, payloadErr := d.fallbackProvider.ReadDocument(path)
		if payloadErr != nil {
			return "", payloadErr
		}
		if trackErr := d.TrackDocument(LocalFileSystemPath(path), payload); trackErr != nil {
			return "", trackErr
		}
		return payload, nil
	}

	return inMemoryDocument.payload, nil
}
