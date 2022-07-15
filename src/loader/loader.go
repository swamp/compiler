/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package loader

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/verbosity"
)

type Loader struct {
	rootPath         string
	documentProvider DocumentProvider
}

func moduleNameToRelativeFilePath(moduleName dectype.PackageRelativeModuleName) string {
	fixed := strings.Replace(moduleName.String(), ".", "/", -1)
	if fixed == "" {
		fixed += "Main"
	}
	if fixed != "" {
		fixed = fixed + ".swamp"
	}

	return fixed
}

func NewLoader(rootPath string, documentProvider DocumentProvider) *Loader {
	return &Loader{rootPath: rootPath, documentProvider: documentProvider}
}

func (l *Loader) Load(relativeModuleName dectype.PackageRelativeModuleName, verboseFlag verbosity.Verbosity) (string, string, decshared.DecoratedError) {
	relativePath := moduleNameToRelativeFilePath(relativeModuleName)
	if path.IsAbs(relativePath) {
		return "", "", decorated.NewInternalError(fmt.Errorf("loader wants relative paths, can not use absolute ones '%s'", relativeModuleName))
	}
	fullPath := path.Join(l.rootPath, relativePath)
	fullPath = filepath.ToSlash(fullPath)
	completeFilename, completeFilenameErr := filepath.Abs(fullPath)
	if completeFilenameErr != nil {
		return "", "", decorated.NewInternalError(completeFilenameErr)
	}
	completeFilename = filepath.ToSlash(completeFilename)
	if verboseFlag >= verbosity.Low {
		log.Printf("* loading file %v\n", completeFilename)
	}

	completeDocumentFilename := LocalFileSystemPath(completeFilename)

	str, err := l.documentProvider.ReadDocument(completeDocumentFilename)
	if err != nil {
		return "", "", decorated.NewModuleNotFoundInDocumentProvider(relativeModuleName, string(completeDocumentFilename), err)
	}

	return fullPath, str, nil
}
