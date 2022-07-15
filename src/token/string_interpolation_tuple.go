/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import "fmt"

// StringInterpolationTuple :
type StringInterpolationTuple struct {
	s StringToken
}

func NewStringInterpolationTuple(s StringToken) StringInterpolationTuple {
	return StringInterpolationTuple{s: s}
}

func (s StringInterpolationTuple) Type() Type {
	return StringInterpolationTupleConstant
}

func (s StringInterpolationTuple) String() string {
	return fmt.Sprintf("[str:%s]", s.s)
}

func (s StringInterpolationTuple) StringToken() StringToken {
	return s.s
}

func (s StringInterpolationTuple) FetchPositionLength() SourceFileReference {
	return s.s.FetchPositionLength()
}
