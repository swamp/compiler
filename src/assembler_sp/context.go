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
	scopeVariables *ScopeVariables
	constants      *Constants
}

func (c *Context) MakeScopeContext() *Context {
	newContext := &Context{
		scopeVariables: NewFunctionVariables(),
		constants:      c.constants,
	}
	return newContext
}

func (c *Context) ScopeVariables() *ScopeVariables {
	return c.scopeVariables
}

func (r *Context) Constants() *Constants {
	return r.constants
}

func (c *Context) Free() {
}

func (c *Context) String() string {
	s := "\n"
	s += fmt.Sprintf("%v\n", c.scopeVariables)
	return strings.TrimSpace(s)
}

func (c *Context) ShowSummary() {
	fmt.Printf("---------- Variables ------------\n")
	fmt.Printf("%v\n", c.scopeVariables)
	fmt.Printf("---------- Constants ------------\n")
	fmt.Printf("%v\n", c.constants)
	fmt.Printf("---------------------------------\n")
}
