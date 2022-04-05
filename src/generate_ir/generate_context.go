package generate_ir

import (
	"github.com/llir/llvm/ir"
	"github.com/swamp/assembler/lib/assembler_sp"
	"github.com/swamp/compiler/src/resourceid"
	"github.com/swamp/compiler/src/typeinfo"
)

type parameterContext struct {
	lookup map[string]*ir.Param
}

func newParameterContext(params []*ir.Param) *parameterContext {
	self := &parameterContext{lookup: make(map[string]*ir.Param)}
	for _, param := range params {
		self.AddParam(param)
	}

	return self
}

func (c* parameterContext) Find(name string) *ir.Param {
	return c.lookup[name]
}

func (c *parameterContext) AddParam(param* ir.Param) {
	lookupName := param.Ident()[1:]
	c.lookup[lookupName] = param
}

type generateContext struct {
	irModule           *ir.Module
	block *ir.Block
	parameterContext* parameterContext
	lookup             typeinfo.TypeLookup
	resourceNameLookup resourceid.ResourceNameLookup
	fileCache          *assembler_sp.FileUrlCache
}
