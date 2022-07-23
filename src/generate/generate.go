package generate

import (
	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/compiler/src/resourceid"
	"github.com/swamp/compiler/src/typeinfo"
	"github.com/swamp/compiler/src/verbosity"
)

type Generator interface {
	GenerateModule(module *decorated.Module,
		lookup typeinfo.TypeLookup, resourceNameLookup resourceid.ResourceNameLookup, fileUrlCache *assembler_sp.FileUrlCache, verboseFlag verbosity.Verbosity) error
}
