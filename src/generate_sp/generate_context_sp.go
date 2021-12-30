package generate_sp

import (
	"github.com/swamp/assembler/lib/assembler_sp"
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/compiler/src/typeinfo"
	"github.com/swamp/opcodes/opcode_sp"
)

type generateContext struct {
	context   *Context
	lookup    typeinfo.TypeLookup
	fileCache *assembler_sp.FileUrlCache
}

func (c *generateContext) MakeScopeContext(debugString string) *generateContext {
	newGenContext := &generateContext{
		context:   c.context.MakeScopeContext(debugString),
		lookup:    c.lookup,
		fileCache: c.fileCache,
	}

	return newGenContext
}

func (c *generateContext) toFilePosition(source token.SourceFileReference) opcode_sp.FilePosition {
	return toFilePosition(c.fileCache, source)
}
