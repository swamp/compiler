/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
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

func (c *VariableContext) ResolveVariable(name *ast.VariableIdentifier) (decorated.Expression, decshared.DecoratedError) {
	def := c.FindNamedDecoratedExpression(name)
	if def == nil {
		return nil, decorated.NewInternalError(fmt.Errorf("couldn't find %v", name))
	}

	def.SetReferenced()

	if constantExpression, wasConstant := isConstant(def.Expression()); wasConstant {
		switch t := constantExpression.(type) {
		case *decorated.IntegerLiteral:
			return t, nil
		case *decorated.StringLiteral:
			return t, nil
		case *decorated.CharacterLiteral:
			return t, nil
		case *decorated.TypeIdLiteral:
			return t, nil
		case *decorated.ResourceNameLiteral:
			return t, nil
		case *decorated.RecordLiteral:
			return t, nil
		case *decorated.ListLiteral:
			return t, nil
		case *decorated.FixedLiteral:
			return t, nil
		}
	}
	switch t := def.Expression().(type) {
	case *decorated.FunctionValue:
		return decorated.NewFunctionReference(name, t), nil
	case *decorated.FunctionParameterDefinition:
		return decorated.NewFunctionParameterReference(name, t), nil
	case *decorated.LetAssignment:
		return decorated.NewLetVariableReference(name, t), nil
	default:
		return nil, decorated.NewInternalError(fmt.Errorf("what to do with '%v' => %T", name, t))
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
