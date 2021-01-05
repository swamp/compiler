/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package pathutil

import (
	"os"
	"path/filepath"
)

func TryToMakeRelativePath(filename string) string {
	cwd, cwdErr := os.Getwd()
	if cwdErr == nil {
		relativeFilename, relativeErr := filepath.Rel(cwd, filename)
		if relativeErr == nil {
			return relativeFilename
		}
	}
	return filename
}
