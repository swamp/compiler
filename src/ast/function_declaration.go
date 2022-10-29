/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type FunctionDeclaration struct {
	parameters             []*FunctionParameter
	expression             Expression
	debugAssignedValue     token.VariableSymbolToken
	commentBlock           *MultilineComment
	inclusive              token.SourceFileReference
	returnType             Type
	functionType           Type
	annotationFunctionType token.AnnotationFunctionType
}

func NewFunctionDeclaration(debugAssignedValue token.VariableSymbolToken, parameters []*FunctionParameter,
	returnType Type, expression Expression, annotationFunctionType token.AnnotationFunctionType, commentBlock *MultilineComment) *FunctionDeclaration {
	inclusive := token.MakeInclusiveSourceFileReference(debugAssignedValue.FetchPositionLength(), expression.FetchPositionLength())
	if inclusive.Range.End().Line() == 0 && inclusive.Range.End().Column() == 0 {
		panic("problem")
	}

	var types []Type
	for _, param := range parameters {
		types = append(types, param.parameterType)
	}
	types = append(types, returnType)

	return &FunctionDeclaration{
		returnType:             returnType,
		annotationFunctionType: annotationFunctionType,
		functionType:           NewFunctionType(types),
		debugAssignedValue:     debugAssignedValue, parameters: parameters,
		expression: expression, commentBlock: commentBlock, inclusive: inclusive,
	}
}

func (i *FunctionDeclaration) AnnotationFunctionType() token.AnnotationFunctionType {
	return i.annotationFunctionType
}

func (i *FunctionDeclaration) IsSomeKindOfExternal() bool {
	return i.annotationFunctionType == token.AnnotationFunctionTypeExternal || i.annotationFunctionType == token.AnnotationFunctionTypeExternalVarEx || i.annotationFunctionType == token.AnnotationFunctionTypeExternalVar
}

func (f *FunctionDeclaration) IsExternalVarFunction() bool {
	return f.annotationFunctionType == token.AnnotationFunctionTypeExternalVar
}

func (f *FunctionDeclaration) IsExternalVarExFunction() bool {
	return f.annotationFunctionType == token.AnnotationFunctionTypeExternalVarEx
}

func (i *FunctionDeclaration) Type() Type {
	return i.functionType
}

func (i *FunctionDeclaration) Parameters() []*FunctionParameter {
	return i.parameters
}

func (i *FunctionDeclaration) CommentBlock() *MultilineComment {
	return i.commentBlock
}

func (i *FunctionDeclaration) Expression() Expression {
	return i.expression
}

func (i *FunctionDeclaration) FetchPositionLength() token.SourceFileReference {
	return i.inclusive
}

func (i *FunctionDeclaration) DebugFunctionIdentifier() token.VariableSymbolToken {
	return i.debugAssignedValue
}

func (i *FunctionDeclaration) String() string {
	return fmt.Sprintf("[FnDecl (%v) => %v = %v]", i.parameters, i.returnType, i.expression)
}

func (i *FunctionDeclaration) DebugString() string {
	return "[function-declaration]"
}
