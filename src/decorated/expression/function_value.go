/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"
	"log"

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

func (a *FunctionParameterDefinition) FetchPositionLength() token.SourceFileReference {
	return a.identifier.Symbol().SourceFileReference
}

type FunctionValue struct {
	forcedFunctionType  *dectype.FunctionAtom
	decoratedExpression Expression
	parameters          []*FunctionParameterDefinition
	commentBlock        token.CommentBlock
	astFunction         *ast.FunctionValue
	sourceFileReference token.SourceFileReference
}

func NewFunctionValue(astFunction *ast.FunctionValue, forcedFunctionType *dectype.FunctionAtom, parameters []*FunctionParameterDefinition, decoratedExpression Expression, commentBlock token.CommentBlock) *FunctionValue {
	if len(parameters) != (forcedFunctionType.ParameterCount() - 1) {
		panic("not great. different parameters")
	}
	for _, parameter := range parameters {
		log.Printf("param %v %v\n", parameter.FetchPositionLength(), parameter)
	}

	sourceFileReference := token.MakeInclusiveSourceFileReference(astFunction.DebugFunctionIdentifier().SourceFileReference, decoratedExpression.FetchPositionLength())

	return &FunctionValue{astFunction: astFunction, forcedFunctionType: forcedFunctionType, parameters: parameters, decoratedExpression: decoratedExpression, commentBlock: commentBlock, sourceFileReference: sourceFileReference}
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

func (f *FunctionValue) Expression() Expression {
	return f.decoratedExpression
}

func (f *FunctionValue) FetchPositionLength() token.SourceFileReference {
	return f.sourceFileReference
}

func (f *FunctionValue) CommentBlock() token.CommentBlock {
	return f.commentBlock
}
