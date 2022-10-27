/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type FunctionValue struct {
	parameters             []*FunctionParameter
	expression             Expression
	debugAssignedValue     token.VariableSymbolToken
	commentBlock           *MultilineComment
	inclusive              token.SourceFileReference
	returnType             Type
	functionType           Type
	annotationFunctionType token.AnnotationFunctionType
}

func NewFunctionValue(debugAssignedValue token.VariableSymbolToken, parameters []*FunctionParameter,
	returnType Type, expression Expression, annotationFunctionType token.AnnotationFunctionType, commentBlock *MultilineComment) *FunctionValue {
	inclusive := token.MakeInclusiveSourceFileReference(debugAssignedValue.FetchPositionLength(), expression.FetchPositionLength())
	if inclusive.Range.End().Line() == 0 && inclusive.Range.End().Column() == 0 {
		panic("problem")
	}

	var types []Type
	for _, param := range parameters {
		types = append(types, param.parameterType)
	}
	types = append(types, returnType)

	return &FunctionValue{
		returnType:             returnType,
		annotationFunctionType: annotationFunctionType,
		functionType:           NewFunctionType(types),
		debugAssignedValue:     debugAssignedValue, parameters: parameters,
		expression: expression, commentBlock: commentBlock, inclusive: inclusive,
	}
}

func (i *FunctionValue) AnnotationFunctionType() token.AnnotationFunctionType {
	return i.annotationFunctionType
}

func (i *FunctionValue) IsSomeKindOfExternal() bool {
	return i.annotationFunctionType == token.AnnotationFunctionTypeExternal || i.annotationFunctionType == token.AnnotationFunctionTypeExternalVarEx || i.annotationFunctionType == token.AnnotationFunctionTypeExternalVar
}

func (f *FunctionValue) IsExternalVarFunction() bool {
	return f.annotationFunctionType == token.AnnotationFunctionTypeExternalVar
}

func (f *FunctionValue) IsExternalVarExFunction() bool {
	return f.annotationFunctionType == token.AnnotationFunctionTypeExternalVarEx
}

func (i *FunctionValue) Type() Type {
	return i.functionType
}

func (i *FunctionValue) Parameters() []*FunctionParameter {
	return i.parameters
}

func (i *FunctionValue) CommentBlock() *MultilineComment {
	return i.commentBlock
}

func (i *FunctionValue) Expression() Expression {
	return i.expression
}

func (i *FunctionValue) FetchPositionLength() token.SourceFileReference {
	return i.inclusive
}

func (i *FunctionValue) DebugFunctionIdentifier() token.VariableSymbolToken {
	return i.debugAssignedValue
}

func (i *FunctionValue) String() string {
	return fmt.Sprintf("[Fn (%v) => %v = %v]", i.parameters, i.returnType, i.expression)
}

func (i *FunctionValue) DebugString() string {
	return "[function]"
}
