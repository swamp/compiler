/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import "fmt"

// VariableSymbolToken :
type TypeSymbolToken struct {
	PositionLength
	raw         string
	Indentation int
}

func NewTypeSymbolToken(raw string, startPosition PositionLength, indentation int) TypeSymbolToken {
	return TypeSymbolToken{raw: raw, PositionLength: startPosition, Indentation: indentation}
}

func (s TypeSymbolToken) Type() Type {
	return TypeSymbol
}

func (s TypeSymbolToken) Name() string {
	return s.raw
}

func (s TypeSymbolToken) Raw() string {
	return s.raw
}

func (s TypeSymbolToken) FetchIndentation() int {
	return s.Indentation
}

func (s TypeSymbolToken) String() string {
	return fmt.Sprintf("$%s", s.raw)
}
