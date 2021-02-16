/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type FunctionParameterDefinition struct {
	identifier    *ast.VariableIdentifier
	generatedType dtype.Type
}

func NewFunctionParameterDefinition(identifier *ast.VariableIdentifier, convertedType dtype.Type) *FunctionParameterDefinition {
	return &FunctionParameterDefinition{identifier: identifier, generatedType: convertedType}
}

func (a *FunctionParameterDefinition) Identifier() *ast.VariableIdentifier {
	return a.identifier
}

func (a *FunctionParameterDefinition) Type() dtype.Type {
	return a.generatedType
}

func (a *FunctionParameterDefinition) String() string {
	return fmt.Sprintf("[arg %v = %v]", a.identifier, a.generatedType)
}

func (a *FunctionParameterDefinition) FetchPositionAndLength() token.PositionLength {
	return a.identifier.Symbol().FetchPositionLength()
}

type FunctionValue struct {
	forcedFunctionType  *dectype.FunctionAtom
	decoratedExpression DecoratedExpression
	parameters          []*FunctionParameterDefinition
}

func NewFunctionValue(forcedFunctionType *dectype.FunctionAtom, parameters []*FunctionParameterDefinition, decoratedExpression DecoratedExpression) *FunctionValue {
	if len(parameters) != (forcedFunctionType.ParameterCount() - 1) {
		panic("not great. different parameters")
	}
	return &FunctionValue{forcedFunctionType: forcedFunctionType, parameters: parameters, decoratedExpression: decoratedExpression}
}

func (f *FunctionValue) Parameters() []*FunctionParameterDefinition {
	return f.parameters
}

func (f *FunctionValue) ForcedFunctionType() *dectype.FunctionAtom {
	return f.forcedFunctionType
}

func (f *FunctionValue) String() string {
	return fmt.Sprintf("[functionvalue (%v) -> %v]", f.parameters, f.decoratedExpression)
}

func (f *FunctionValue) DecoratedName() string {
	return fmt.Sprintf("[functionvalue (%v) -> %v]", f.parameters, f.decoratedExpression)
}

func (f *FunctionValue) HumanReadable() string {
	return fmt.Sprintf("[functionvalue (%v) -> %v]", f.parameters, f.decoratedExpression)
}

func (f *FunctionValue) ShortName() string {
	return fmt.Sprintf("[functionvalue (%v) -> %v]", f.parameters, f.decoratedExpression)
}

func (f *FunctionValue) ShortString() string {
	return fmt.Sprintf("[functionvalue (%v) -> %v]", f.parameters, f.decoratedExpression)
}

func (f *FunctionValue) DebugString() string {
	return fmt.Sprintf("[functionval]")
}

func (f *FunctionValue) Type() dtype.Type {
	return f.forcedFunctionType
}

func (f *FunctionValue) Next() dtype.Type {
	return f.forcedFunctionType
}

func (f *FunctionValue) Resolve() (dtype.Atom, error) {
	return f.forcedFunctionType.Resolve()
}

func (f *FunctionValue) Expression() DecoratedExpression {
	return f.decoratedExpression
}

func (f *FunctionValue) FetchPositionAndLength() token.PositionLength {
	return f.decoratedExpression.FetchPositionAndLength()
}
