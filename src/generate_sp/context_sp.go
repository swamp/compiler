/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

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

func NewContext(packageConstants *assembler_sp.PackageConstants, debugString string) *Context {
	return &Context{
		constants:      packageConstants,
		inFunction:     nil,
		scopeVariables: assembler_sp.NewFunctionVariables(debugString),
		stackMemory:    assembler_sp.NewStackMemoryMapper(32 * 1024),
	}
}

func (c *Context) MakeScopeContext(debugString string) *Context {
	newContext := &Context{
		constants:      c.constants,
		inFunction:     c.inFunction,
		scopeVariables: assembler_sp.NewFunctionVariablesWithParent(c.scopeVariables, debugString),
		stackMemory:    c.stackMemory,
	}

	return newContext
}

func (c *Context) MakeFunctionContext(inFunction *decorated.FunctionValue, debugString string) *Context {
	newContext := &Context{
		constants:      c.constants,
		inFunction:     inFunction,
		scopeVariables: assembler_sp.NewFunctionVariables(debugString),
		stackMemory:    assembler_sp.NewStackMemoryMapper(32 * 1024),
	}

	return newContext
}

func (c *Context) Constants() *assembler_sp.PackageConstants {
	return c.constants
}
