package generate

import (
	"github.com/swamp/compiler/src/loader"
	"github.com/swamp/compiler/src/resourceid"
	"github.com/swamp/compiler/src/verbosity"
)

type Generator interface {
	GenerateFromPackage(module *loader.Package, resourceNameLookup resourceid.ResourceNameLookup, verboseFlag verbosity.Verbosity) error
}
