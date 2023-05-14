/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"
	"strings"

	"github.com/swamp/compiler/src/token"
)

type FunctionValue struct {
	parameters           []*FunctionParameter
	expression           Expression
	debugAssignedValue   token.VariableSymbolToken
	commentBlock         *MultilineComment
	inclusive            token.SourceFileReference
	declaredFunctionType Type
}

func NewFunctionValue(debugAssignedValue token.VariableSymbolToken, parameters []*FunctionParameter,
	declaredFunctionType Type, expression Expression, commentBlock *MultilineComment) *FunctionValue {
	inclusive := token.MakeInclusiveSourceFileReference(debugAssignedValue.FetchPositionLength(), expression.FetchPositionLength())

	return &FunctionValue{
		declaredFunctionType: declaredFunctionType,
		debugAssignedValue:   debugAssignedValue, parameters: parameters,
		expression: expression, commentBlock: commentBlock, inclusive: inclusive,
	}
}

func (i *FunctionValue) Type() Type {
	return i.declaredFunctionType
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

func (i *FunctionValue) parametersString() string {
	var typeParams []string
	if len(i.parameters) == 0 {
		return ""
	}
	for _, p := range i.parameters {
		typeParams = append(typeParams, p.Name())
	}
	return strings.Join(typeParams, ", ")
}

func (i *FunctionValue) String() string {
	return fmt.Sprintf("[Fn %v (%s) = %v]", i.declaredFunctionType, i.parametersString(), i.expression)
}

func (i *FunctionValue) DebugString() string {
	return "[function]"
}
