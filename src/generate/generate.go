/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate

import (
	"github.com/swamp/compiler/src/loader"
	"github.com/swamp/compiler/src/resourceid"
	"github.com/swamp/compiler/src/verbosity"
)

type Generator interface {
	GenerateFromPackage(module *loader.Package, resourceNameLookup resourceid.ResourceNameLookup, outputDirectory string, packageSubDirectory string, verboseFlag verbosity.Verbosity) error
}
