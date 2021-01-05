/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type LambdaFunctionValue struct {
	parameters  []*VariableIdentifier
	expression  Expression
	lambdaToken token.Token
}

func NewLambdaFunctionValue(lambdaToken token.Token, parameters []*VariableIdentifier, expression Expression) *LambdaFunctionValue {
	return &LambdaFunctionValue{lambdaToken: lambdaToken, parameters: parameters, expression: expression}
}

func (i *LambdaFunctionValue) Parameters() []*VariableIdentifier {
	return i.parameters
}

func (i *LambdaFunctionValue) Expression() Expression {
	return i.expression
}

func (i *LambdaFunctionValue) Token() token.Token {
	return i.lambdaToken
}

func (i * LambdaFunctionValue) PositionLength() token.PositionLength {
	return i.lambdaToken.FetchPositionLength()
}

func (i *LambdaFunctionValue) String() string {
	return fmt.Sprintf("[lambda (%v) -> %v]", i.parameters, i.expression)
}

func (i *LambdaFunctionValue) DebugString() string {
	return fmt.Sprintf("[lambda]")
}
