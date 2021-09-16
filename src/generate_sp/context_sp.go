package generate_sp

import (
	"github.com/swamp/compiler/src/assembler_sp"
)

type Context struct {
	startMemoryConstants *StartMemoryConstants
	constants            *assembler_sp.Constants
	functionVariables    *assembler_sp.FunctionVariables
	stackMemory          *assembler_sp.StackMemoryMapper
}

func NewContext() *Context {
	return &Context{
		startMemoryConstants: NewStartMemoryConstants(),
		constants:            assembler_sp.NewConstants(),
		functionVariables:    assembler_sp.NewFunctionVariables(),
		stackMemory:          assembler_sp.NewStackMemoryMapper(32 * 1024),
	}
}

func (c *Context) StartMemoryConstants() *StartMemoryConstants {
	return c.startMemoryConstants
}

func (c *Context) MakeScopeContext() *Context {
	newContext := &Context{
		startMemoryConstants: c.startMemoryConstants,
		constants:            c.constants,
		functionVariables:    assembler_sp.NewFunctionVariablesWithParent(c.functionVariables),
		stackMemory:          c.stackMemory,
	}

	return newContext
}

func (c *Context) MakeFunctionContext() *Context {
	newContext := &Context{
		startMemoryConstants: c.startMemoryConstants,
		constants:            c.constants,
		functionVariables:    assembler_sp.NewFunctionVariables(),
		stackMemory:          assembler_sp.NewStackMemoryMapper(32 * 1024),
	}

	return newContext
}

func (c *Context) Constants() *assembler_sp.Constants {
	return c.constants
}
