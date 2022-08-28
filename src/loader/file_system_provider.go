/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package loader

import (
	"fmt"
	"os"

	decorated "github.com/swamp/compiler/src/decorated/expression"
)

type FileSystemDocumentProvider struct{}

func NewFileSystemDocumentProvider() *FileSystemDocumentProvider {
	return &FileSystemDocumentProvider{}
}

func (s *FileSystemDocumentProvider) ReadDocument(completeFilename LocalFileSystemPath) (string, error) {
	octets, readFileErr := os.ReadFile(string(completeFilename))
	if readFileErr != nil {
		switch readFileErr {
		case os.ErrInvalid:
		case os.ErrPermission:
		case os.ErrNotExist:
			return "", decorated.NewInternalError(fmt.Errorf("file '%s' didn't exist (%v)", completeFilename, readFileErr))
		default:
			switch v := readFileErr.(type) {
			case *os.PathError:
				return "", decorated.NewInternalError(fmt.Errorf("couldn't find file '%s' (%v)", completeFilename, v))
			default:
				return "", decorated.NewInternalError(fmt.Errorf("couldn't open file '%s'", completeFilename))
			}
		}
	}

	str := string(octets)
	return str, nil
}
