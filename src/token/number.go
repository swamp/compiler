/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import (
	"fmt"
)

// NumberToken :
type NumberToken struct {
	PositionLength
	number int32
	raw    string
	isFixed bool
}

func NewNumberToken(raw string, v int32, isFixed bool, startPosition PositionLength) NumberToken {
	return NumberToken{raw: raw, number: v, isFixed: isFixed, PositionLength: startPosition}
}

func (s NumberToken) Type() Type {
	if s.isFixed {
		return NumberFixed
	}

	return NumberInteger
}

func (s NumberToken) Value() int32 {
	return s.number
}

func (s NumberToken) Raw() string {
	return s.raw
}

func (s NumberToken) String() string {
	return fmt.Sprintf("Number:%v", s.number)
}
