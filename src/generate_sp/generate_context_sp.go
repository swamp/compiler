package generate_sp

import (
	"github.com/swamp/compiler/src/typeinfo"
)

type generateContext struct {
	context *Context
	lookup  typeinfo.TypeLookup
}

func (c *generateContext) MakeScopeContext() *generateContext {
	newGenContext := &generateContext{
		context: c.context.MakeScopeContext(),
		lookup:  c.lookup,
	}

	return newGenContext
}
