/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import (
	"fmt"
)

// Keyword :
type BooleanToken struct {
	Range
	value bool
	raw   string
}

func NewBooleanToken(raw string, v bool, startPosition Range) BooleanToken {
	return BooleanToken{raw: raw, value: v, Range: startPosition}
}

func (s BooleanToken) Type() Type {
	return BooleanType
}

func (s BooleanToken) Value() bool {
	return s.value
}

func (s BooleanToken) Raw() string {
	return s.raw
}

func (s BooleanToken) String() string {
	return fmt.Sprintf("{bool: %v}", s.value)
}
