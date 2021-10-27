package generate_sp

import (
	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

type Context struct {
	constants      *assembler_sp.PackageConstants
	scopeVariables *assembler_sp.ScopeVariables
	stackMemory    *assembler_sp.StackMemoryMapper
	inFunction     *decorated.FunctionValue
}

func NewContext(packageConstants *assembler_sp.PackageConstants) *Context {
	return &Context{
		constants:      packageConstants,
		scopeVariables: assembler_sp.NewFunctionVariables(),
		stackMemory:    assembler_sp.NewStackMemoryMapper(32 * 1024),
	}
}

func (c *Context) MakeScopeContext() *Context {
	newContext := &Context{
		constants:      c.constants,
		inFunction:     c.inFunction,
		scopeVariables: assembler_sp.NewFunctionVariablesWithParent(c.scopeVariables),
		stackMemory:    c.stackMemory,
	}

	return newContext
}

func (c *Context) MakeFunctionContext(inFunction *decorated.FunctionValue) *Context {
	newContext := &Context{
		constants:      c.constants,
		inFunction:     inFunction,
		scopeVariables: assembler_sp.NewFunctionVariables(),
		stackMemory:    assembler_sp.NewStackMemoryMapper(32 * 1024),
	}

	return newContext
}

func (c *Context) Constants() *assembler_sp.PackageConstants {
	return c.constants
}
