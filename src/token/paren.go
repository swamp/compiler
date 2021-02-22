/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

// ParenToken :
type ParenToken struct {
	Range
	operatorType Type
	raw          string
	debugString  string
}

func NewParenToken(raw string, operatorType Type, startPosition Range, debugString string) ParenToken {
	return ParenToken{operatorType: operatorType, Range: startPosition, raw: raw, debugString: debugString}
}

func (s ParenToken) Type() Type {
	return s.operatorType
}

func (s ParenToken) String() string {
	return s.debugString
}

func (s ParenToken) Raw() string {
	return s.raw
}
