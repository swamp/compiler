/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import (
	"fmt"
)

// Keyword :
type Keyword struct {
	Range
	raw       string
	tokenType Type
}

func NewKeyword(raw string, tokenType Type, startPosition Range) Keyword {
	return Keyword{raw: raw, tokenType: tokenType, Range: startPosition}
}

func (s Keyword) Type() Type {
	return s.tokenType
}

func (s Keyword) Raw() string {
	return s.raw
}

func (s Keyword) String() string {
	if s.tokenType == TypeDef {
		return "TYPE"
	}
	return fmt.Sprintf("{%v} (%s)", s.tokenType, s.raw)
}
