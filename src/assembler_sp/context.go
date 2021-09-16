/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import (
	"fmt"
	"strings"
)

type Context struct {
	functionVariables *FunctionVariables
	parent            *Context
	root              *FunctionRootContext
	constants         *Constants
}

func (c *Context) MakeScopeContext() *Context {
	newContext := &Context{
		functionVariables: NewFunctionVariables(),
		root:              c.root, parent: c, constants: c.constants,
	}
	newContext.parent = c
	return newContext
}

func (c *Context) Parent() *Context {
	return c.parent
}

func (c *Context) ScopeVariables() *FunctionVariables {
	return c.functionVariables
}

func (r *Context) Constants() *Constants {
	return r.constants
}

func (c *Context) Free() {
}

func (c *Context) String() string {
	s := "\n"
	s += fmt.Sprintf("%v\n", c.functionVariables)
	return strings.TrimSpace(s)
}

func (c *Context) ShowSummary() {
	fmt.Printf("---------- Variables ------------\n")
	fmt.Printf("%v\n", c.functionVariables)
	fmt.Printf("---------- Constants ------------\n")
	fmt.Printf("%v\n", c.constants)
	fmt.Printf("---------------------------------\n")
}

type FunctionRootContext struct {
	constants    *Constants
	scopeContext *Context
}

func NewFunctionRootContext() *FunctionRootContext {
	c := &FunctionRootContext{constants: &Constants{}}
	bootstrap := &Context{root: c, constants: c.constants}
	c.scopeContext = bootstrap.MakeScopeContext()
	return c
}

func (r *FunctionRootContext) ScopeContext() *Context {
	return r.scopeContext
}

func (r *FunctionRootContext) Constants() *Constants {
	return r.constants
}

func (r *FunctionRootContext) String() string {
	return r.constants.String() + r.scopeContext.String()
}
