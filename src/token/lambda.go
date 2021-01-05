/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

type LambdaToken struct {
	PositionLength
	debugString string
}

func NewLambdaToken(startPosition PositionLength, debugString string) LambdaToken {
	return LambdaToken{PositionLength: startPosition, debugString: debugString}
}

func (s LambdaToken) Type() Type {
	return Lambda
}

func (s LambdaToken) String() string {
	return "lambda"
}
