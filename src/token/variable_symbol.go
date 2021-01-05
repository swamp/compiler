/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import "fmt"

// VariableSymbolToken :
type VariableSymbolToken struct {
	PositionLength
	raw         string
	Indentation int
}

func NewVariableSymbolToken(raw string, startPosition PositionLength, indentation int) VariableSymbolToken {
	return VariableSymbolToken{raw: raw, PositionLength: startPosition, Indentation: indentation}
}

func (s VariableSymbolToken) Type() Type {
	return VariableSymbol
}

func (s VariableSymbolToken) Name() string {
	return s.raw
}

func (s VariableSymbolToken) Raw() string {
	return s.raw
}

func (s VariableSymbolToken) FetchIndentation() int {
	return s.Indentation
}

func (s VariableSymbolToken) String() string {
	return fmt.Sprintf("$%s", s.raw)
}
