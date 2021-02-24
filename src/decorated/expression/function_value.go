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

func (a *FunctionParameterDefinition) FetchPositionLength() token.Range {
	return a.identifier.Symbol().FetchPositionLength()
}

type FunctionValue struct {
	forcedFunctionType  *dectype.FunctionAtom
	decoratedExpression DecoratedExpression
	parameters          []*FunctionParameterDefinition
	commentBlock        token.CommentBlock
	astFunction         *ast.FunctionValue
}

func NewFunctionValue(astFunction *ast.FunctionValue, forcedFunctionType *dectype.FunctionAtom, parameters []*FunctionParameterDefinition, decoratedExpression DecoratedExpression, commentBlock token.CommentBlock) *FunctionValue {
	if len(parameters) != (forcedFunctionType.ParameterCount() - 1) {
		panic("not great. different parameters")
	}
	start := astFunction.DebugFunctionIdentifier().FetchPositionLength().Start()
	end := decoratedExpression.FetchPositionLength().End()

	return &FunctionValue{astFunction: astFunction, forcedFunctionType: forcedFunctionType, parameters: parameters, decoratedExpression: decoratedExpression, commentBlock: commentBlock}
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

func (f *FunctionValue) FetchPositionLength() token.Range {
	return f.decoratedExpression.FetchPositionLength()
}

func (f *FunctionValue) CommentBlock() token.CommentBlock {
	return f.commentBlock
}
