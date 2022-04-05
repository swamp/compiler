package generate_ir

import (
	"fmt"
	"github.com/llir/llvm/ir"
	"github.com/swamp/assembler/lib/assembler_sp"
	"github.com/swamp/compiler/src/resourceid"
	"github.com/swamp/compiler/src/typeinfo"
	"strings"
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

func (c *parameterContext) String() string {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("parameterContext\n"))

	for name, param := range c.lookup {
		output.WriteString(fmt.Sprintf(  "  %v : %v\n", name, param.LLString()))
	}

	return output.String()
}

func (c *parameterContext) Find(name string) *ir.Param {
	return c.lookup[name]
}

func (c *parameterContext) AddParam(param *ir.Param) {
	lookupName := param.Ident()[1:]
	c.lookup[lookupName] = param
}

type generateContext struct {
	irModule           *ir.Module
	block              *ir.Block
	parameterContext   *parameterContext
	lookup             typeinfo.TypeLookup
	resourceNameLookup resourceid.ResourceNameLookup
	fileCache          *assembler_sp.FileUrlCache
}
