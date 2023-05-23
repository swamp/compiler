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
	functionParameter *ast.FunctionParameter
	generatedType     dtype.Type `debug:"true"`
	references        []*FunctionParameterReference
}

func NewFunctionParameterDefinition(identifier *ast.FunctionParameter,
	convertedType dtype.Type) *FunctionParameterDefinition {
	if identifier == nil {
		panic(fmt.Errorf("functionParameter must be set"))
	}
	return &FunctionParameterDefinition{functionParameter: identifier, generatedType: convertedType}
}

func (a *FunctionParameterDefinition) Parameter() *ast.FunctionParameter {
	return a.functionParameter
}

func (a *FunctionParameterDefinition) Type() dtype.Type {
	return a.generatedType
}

func (a *FunctionParameterDefinition) String() string {
	return fmt.Sprintf("[Param %v : %v]", a.functionParameter.SymbolName(), a.generatedType)
}

func (a *FunctionParameterDefinition) HumanReadable() string {
	return "Function Parameter"
}

func (a *FunctionParameterDefinition) FetchPositionLength() token.SourceFileReference {
	return a.functionParameter.FetchPositionLength()
}

func (a *FunctionParameterDefinition) AddReferee(ref *FunctionParameterReference) {
	a.references = append(a.references, ref)
}

func (a *FunctionParameterDefinition) WasReferenced() bool {
	return len(a.references) > 0
}

func (a *FunctionParameterDefinition) References() []*FunctionParameterReference {
	return a.references
}

type FunctionValue struct {
	declaredFunctionType     dtype.Type                     `debug:"true"`
	decoratedExpression      Expression                     `debug:"true"`
	parameters               []*FunctionParameterDefinition `debug:"true"`
	commentBlock             *ast.MultilineComment
	astFunction              *ast.FunctionValue
	sourceFileReference      token.SourceFileReference
	references               []*FunctionReference
	declaredFunctionTypeAtom *dectype.FunctionAtom
}

func NewPrepareFunctionValue(astFunction *ast.FunctionValue, declaredFunctionType dtype.Type,
	parameters []*FunctionParameterDefinition, declaredFunctionTypeAtom *dectype.FunctionAtom,
	commentBlock *ast.MultilineComment) *FunctionValue {
	if len(parameters) != (declaredFunctionTypeAtom.ParameterCount() - 1) {
		panic(fmt.Errorf("not great. different number of parameters %d vs %v", len(parameters),
			declaredFunctionTypeAtom))
	}
	for _, param := range parameters {
		log.Printf("param %v %T", param.Parameter().Name(), param.Type().Next())
	}
	if declaredFunctionType == nil {
		panic("must provide forced function type")
	}
	return &FunctionValue{
		astFunction: astFunction, declaredFunctionType: declaredFunctionType,
		declaredFunctionTypeAtom: declaredFunctionTypeAtom,
		parameters:               parameters, decoratedExpression: nil, commentBlock: commentBlock,
		sourceFileReference: astFunction.DebugFunctionIdentifier().SourceFileReference,
	}
}

func (f *FunctionValue) DeclaredFunctionTypeAtom() *dectype.FunctionAtom {
	return f.declaredFunctionTypeAtom
}

func (f *FunctionValue) DefineExpression(decoratedExpression Expression) {
	f.sourceFileReference = token.MakeInclusiveSourceFileReference(
		f.astFunction.DebugFunctionIdentifier().SourceFileReference, decoratedExpression.FetchPositionLength(),
	)
	f.decoratedExpression = decoratedExpression
}

func (f *FunctionValue) AstFunctionValue() *ast.FunctionValue {
	return f.astFunction
}

func (f *FunctionValue) IsSomeKindOfExternal() bool {
	decl, wasDecl := f.AstFunctionValue().Expression().(*ast.FunctionDeclarationExpression)
	if !wasDecl {
		return false
	}

	return decl.IsSomeKindOfExternal()
}

func (f *FunctionValue) Parameters() []*FunctionParameterDefinition {
	return f.parameters
}

func (f *FunctionValue) DeclaredFunctionTypeAtom2() *dectype.FunctionAtom {
	return dectype.DerefFunctionType(f.declaredFunctionType)
}

func (f *FunctionValue) String() string {
	return fmt.Sprintf("[FunctionValue %v (%v) -> %v]", f.declaredFunctionType, f.parameters, f.decoratedExpression)
}

func (f *FunctionValue) HumanReadable() string {
	return fmt.Sprintf("[FunctionValue (%v) -> %v]", f.parameters, f.decoratedExpression)
}

func (f *FunctionValue) Type() dtype.Type {
	return f.declaredFunctionType
}

func (f *FunctionValue) Next() dtype.Type {
	return f.declaredFunctionType
}

func (f *FunctionValue) Resolve() (dtype.Atom, error) {
	return f.declaredFunctionType.Resolve()
}

func (f *FunctionValue) Expression() Expression {
	return f.decoratedExpression
}

func (f *FunctionValue) FetchPositionLength() token.SourceFileReference {
	return f.sourceFileReference
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
