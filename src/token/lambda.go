/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

type LambdaToken struct {
	SourceFileReference
	debugString string
}

func NewLambdaToken(startPosition SourceFileReference, debugString string) LambdaToken {
	return LambdaToken{SourceFileReference: startPosition, debugString: debugString}
}

func (s LambdaToken) Type() Type {
	return Lambda
}

func (s LambdaToken) String() string {
	return "lambda"
}

func (s LambdaToken) FetchPositionLength() SourceFileReference {
	return s.SourceFileReference
}
