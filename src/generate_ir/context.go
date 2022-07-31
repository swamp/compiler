/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_ir

import (
	"fmt"
	"github.com/llir/llvm/ir"
	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
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
		output.WriteString(fmt.Sprintf("  %v : %v\n", name, param.LLString()))
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
	irTypeRepo         *IrTypeRepo
	irFunctions        *IrFunctions
	parameterContext   *parameterContext
	lookup             typeinfo.TypeLookup
	resourceNameLookup resourceid.ResourceNameLookup
	fileCache          *assembler_sp.FileUrlCache
	inFunction         *decorated.FunctionValue
}

func (x *generateContext) NewBlock(name string) *generateContext {
	newContext := &generateContext{
		irModule:           x.irModule,
		block:              ir.NewBlock(name),
		irTypeRepo:         x.irTypeRepo,
		parameterContext:   newParameterContext(nil),
		lookup:             x.lookup,
		resourceNameLookup: x.resourceNameLookup,
		fileCache:          x.fileCache,
		irFunctions:        x.irFunctions,
	}

	return newContext
}
