/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import "fmt"

// OperatorToken :
type OperatorToken struct {
	SourceFileReference
	operatorType Type
	raw          string `debug:"true"`
	debugString  string
}

func NewOperatorToken(operatorType Type, startPosition SourceFileReference, raw string, debugString string) OperatorToken {
	return OperatorToken{operatorType: operatorType, SourceFileReference: startPosition, raw: raw, debugString: debugString}
}

func (s OperatorToken) Type() Type {
	return s.operatorType
}

func (s OperatorToken) String() string {
	return s.debugString
}

func (s OperatorToken) Raw() string {
	return s.raw
}

func (s OperatorToken) DebugString() string {
	return fmt.Sprintf("[operator %s]", s.debugString)
}

func (s OperatorToken) FetchPositionLength() SourceFileReference {
	return s.SourceFileReference
}
