/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

type VariableContext struct {
	parent            *VariableContext
	lookup            map[string]*decorated.NamedDecoratedExpression
	parentDefinitions *decorated.ModuleDefinitionsCombine
}

func NewVariableContext(parentDefinitions *decorated.ModuleDefinitionsCombine) *VariableContext {
	return &VariableContext{parent: nil, parentDefinitions: parentDefinitions, lookup: make(map[string]*decorated.NamedDecoratedExpression)}
}

func (c *VariableContext) ResolveVariable(name *ast.VariableIdentifier) decorated.Expression {
	def := c.FindNamedDecoratedExpression(name)
	if def != nil {
		def.SetReferenced()
	}

	switch t := def.Expression().(type) {
	case *decorated.FunctionValue:
		return decorated.NewFunctionReference(name, t)
	case *decorated.FunctionParameterDefinition:
		return decorated.NewFunctionParameterReference(name, t)
	case *decorated.LetAssignment:
		return decorated.NewLetVariableReference(name, t)
	default:
		panic(fmt.Errorf("what to do with '%v' => %T", name, t))
	}
}

func (c *VariableContext) InternalLookups() map[string]*decorated.NamedDecoratedExpression {
	return c.lookup
}

func (c *VariableContext) FindNamedDecoratedExpression(name *ast.VariableIdentifier) *decorated.NamedDecoratedExpression {
	def := c.lookup[name.Name()]
	if def == nil {
		if c.parent != nil {
			return c.parent.FindNamedDecoratedExpression(name)
		}
		mDef := c.parentDefinitions.FindDefinitionExpression(name)
		if mDef == nil {
			return nil
		}

		def = decorated.NewNamedDecoratedExpression(mDef.FullyQualifiedVariableName().String(), mDef, mDef.Expression())
	}

	if def != nil {
		def.SetReferenced()
	}
	return def
}

func (c *VariableContext) Add(name *ast.VariableIdentifier, namedExpression *decorated.NamedDecoratedExpression) {
	c.lookup[name.Name()] = namedExpression
}

func (c *VariableContext) String() string {
	s := "[context \n"
	for name, contextType := range c.lookup {
		s += fmt.Sprintf("   %v = %v\n", name, contextType)
	}
	if c.parent != nil {
		s += c.parent.String()
	}
	s += "\n]"
	return s
}

func (c *VariableContext) MakeVariableContext() *VariableContext {
	return &VariableContext{parent: c, lookup: make(map[string]*decorated.NamedDecoratedExpression)}
}
