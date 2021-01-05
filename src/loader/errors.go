/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package loader

import (
	"fmt"

	dectype "github.com/swamp/compiler/src/decorated/types"
)

type CircularDependencyDetected struct {
	loadedModules []dectype.PackageRelativeModuleName
	lastModule    dectype.ArtifactFullyQualifiedModuleName
}

func NewCircularDependencyDetected(loadedModules []dectype.PackageRelativeModuleName, lastModule dectype.ArtifactFullyQualifiedModuleName) *CircularDependencyDetected {
	return &CircularDependencyDetected{loadedModules: loadedModules, lastModule: lastModule}
}

func (e *CircularDependencyDetected) Error() string {
	return fmt.Sprintf("Circular dependency %v %v", e.loadedModules, e.lastModule)
}
