/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

type LambdaToken struct {
	Range
	debugString string
}

func NewLambdaToken(startPosition Range, debugString string) LambdaToken {
	return LambdaToken{Range: startPosition, debugString: debugString}
}

func (s LambdaToken) Type() Type {
	return Lambda
}

func (s LambdaToken) String() string {
	return "lambda"
}
