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
	references    []*FunctionParameterReference
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

func (a *FunctionParameterDefinition) HumanReadable() string {
	return "Function Parameter"
}

func (a *FunctionParameterDefinition) FetchPositionLength() token.SourceFileReference {
	return a.identifier.Symbol().SourceFileReference
}

func (a *FunctionParameterDefinition) AddReferee(ref *FunctionParameterReference) {
	a.references = append(a.references, ref)
}

func (a *FunctionParameterDefinition) References() []*FunctionParameterReference {
	return a.references
}

type FunctionValue struct {
	forcedFunctionType  dectype.FunctionTypeLike
	decoratedExpression Expression
	parameters          []*FunctionParameterDefinition
	commentBlock        *ast.MultilineComment
	astFunction         *ast.FunctionValue
	sourceFileReference token.SourceFileReference
	references          []*FunctionReference
	annotation          *AnnotationStatement
}

func NewFunctionValue(annotation *AnnotationStatement, astFunction *ast.FunctionValue, forcedFunctionType dectype.FunctionTypeLike, parameters []*FunctionParameterDefinition, decoratedExpression Expression, commentBlock *ast.MultilineComment) *FunctionValue {
	if len(parameters) != (forcedFunctionType.ParameterCount() - 1) {
		panic("not great. different parameters")
	}

	sourceFileReference := token.MakeInclusiveSourceFileReference(astFunction.DebugFunctionIdentifier().SourceFileReference, decoratedExpression.FetchPositionLength())

	return &FunctionValue{annotation: annotation, astFunction: astFunction, forcedFunctionType: forcedFunctionType, parameters: parameters, decoratedExpression: decoratedExpression, commentBlock: commentBlock, sourceFileReference: sourceFileReference}
}

func (f *FunctionValue) AstFunctionValue() *ast.FunctionValue {
	return f.astFunction
}

func (f *FunctionValue) Annotation() *AnnotationStatement {
	return f.annotation
}

func (f *FunctionValue) Parameters() []*FunctionParameterDefinition {
	return f.parameters
}

func (f *FunctionValue) ForcedFunctionType() dectype.FunctionTypeLike {
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
	return f.annotation.Type()
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
	return f.astFunction.FetchPositionLength()
}

func (f *FunctionValue) CommentBlock() *ast.MultilineComment {
	return f.commentBlock
}

func (f *FunctionValue) AddReferee(ref *FunctionReference) {
	f.references = append(f.references, ref)
}

func (f *FunctionValue) References() []*FunctionReference {
	return f.references
}
