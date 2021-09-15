package generate_sp

import (
	decorator "github.com/swamp/compiler/src/decorated/convert"
	"github.com/swamp/compiler/src/typeinfo"
)

type generateContext struct {
	context     *Context
	definitions *decorator.VariableContext
	lookup      typeinfo.TypeLookup
}

func (c *generateContext) MakeScopeContext() *generateContext {
	newGenContext := &generateContext{
		context:     c.context.MakeScopeContext(),
		definitions: c.definitions,
		lookup:      c.lookup,
	}

	return newGenContext
}
