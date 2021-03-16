/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"
	"log"

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
	if parentDefinitions == nil {
		panic("parentDefinitions nil")
	}
	return &VariableContext{parent: nil, parentDefinitions: parentDefinitions, lookup: make(map[string]*decorated.NamedDecoratedExpression)}
}

func ReferenceFromVariable(name ast.ScopedOrNormalVariableIdentifier, expression decorated.Expression, module *decorated.Module) (decorated.Expression, decshared.DecoratedError) {
	switch t := expression.(type) {
	case *decorated.FunctionValue:
		var moduleRef *decorated.ModuleReference
		scoped, wasScoped := name.(*ast.VariableIdentifierScoped)
		if wasScoped {
			moduleRef = decorated.NewModuleReference(scoped.ModuleReference(), module)
		}
		nameWithModuleRef := decorated.NewNamedDefinitionReference(moduleRef, name)
		functionReference := decorated.NewFunctionReference(nameWithModuleRef, t)
		return functionReference, nil
	case *decorated.FunctionParameterDefinition:
		return decorated.NewFunctionParameterReference(name, t), nil
	case *decorated.LetVariable:
		return decorated.NewLetVariableReference(name, t), nil
	case *decorated.CaseConsequenceParameterForCustomType:
		return decorated.NewCaseConsequenceParameterReference(name, t), nil
	default:
		return nil, decorated.NewInternalError(fmt.Errorf("what to do with '%v' => %T", name, t))
	}
}

func createConstantHelper(identifier ast.ScopedOrNormalVariableIdentifier, functionReference *decorated.FunctionReference) (*decorated.Constant, decshared.DecoratedError) {
	return decorated.NewConstant(identifier, functionReference), nil
}

func createConstant(identifier ast.ScopedOrNormalVariableIdentifier, functionReference *decorated.FunctionReference) (*decorated.Constant, decshared.DecoratedError) {
	perhapsConstantExpression := functionReference.FunctionValue().Expression()
	switch perhapsConstantExpression.(type) {
	case *decorated.IntegerLiteral:
		return createConstantHelper(identifier, functionReference)
	case *decorated.StringLiteral:
		return createConstantHelper(identifier, functionReference)
	case *decorated.CharacterLiteral:
		return createConstantHelper(identifier, functionReference)
	case *decorated.TypeIdLiteral:
		return createConstantHelper(identifier, functionReference)
	case *decorated.ResourceNameLiteral:
		return createConstantHelper(identifier, functionReference)
	case *decorated.RecordLiteral:
		return createConstantHelper(identifier, functionReference)
	case *decorated.ListLiteral:
		return createConstantHelper(identifier, functionReference)
	case *decorated.FixedLiteral:
		return createConstantHelper(identifier, functionReference)
	}

	return nil, decorated.NewInternalError(fmt.Errorf("not a constant %T", perhapsConstantExpression))
}

func (c *VariableContext) ResolveVariable(name *ast.VariableIdentifier) (decorated.Expression, decshared.DecoratedError) {
	def := c.FindNamedDecoratedExpression(name)
	if def == nil {
		return nil, decorated.NewInternalError(fmt.Errorf("couldn't find %v", name))
	}

	def.SetReferenced()

	someReference, err := ReferenceFromVariable(name, def.Expression(), nil)
	if err != nil {
		return nil, err
	}
	if functionReference, wasConstant := isConstant(someReference); wasConstant {
		return createConstant(name, functionReference)
	}
	return someReference, nil
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

func (c *VariableContext) FindScopedNamedDecoratedExpression(name *ast.VariableIdentifierScoped) *decorated.NamedDecoratedExpression {
	if c.parentDefinitions == nil {
		log.Printf("it was scoped, but I don't have any parent definitions %v", name)
		return nil
	}
	mDef := c.parentDefinitions.FindScopedDefinitionExpression(name)
	if mDef == nil {
		return nil
	}

	def := decorated.NewNamedDecoratedExpression(mDef.FullyQualifiedVariableName().String(), mDef, mDef.Expression())

	return def
}

func (c *VariableContext) FindScopedNamedDecoratedExpressionScopedOrNormal(name ast.ScopedOrNormalVariableIdentifier) *decorated.NamedDecoratedExpression {
	scoped, wasScoped := name.(*ast.VariableIdentifierScoped)
	if wasScoped {
		return c.FindScopedNamedDecoratedExpression(scoped)
	}
	return c.FindNamedDecoratedExpression(name.(*ast.VariableIdentifier))
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
	return &VariableContext{parent: c, lookup: make(map[string]*decorated.NamedDecoratedExpression), parentDefinitions: c.parentDefinitions}
}
