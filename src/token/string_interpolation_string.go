/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import "fmt"

// StringInterpolationString :
type StringInterpolationString struct {
	s StringToken
}

func NewStringInterpolationString(s StringToken) StringInterpolationString {
	return StringInterpolationString{s: s}
}

func (s StringInterpolationString) Type() Type {
	return StringInterpolationStringConstant
}

func (s StringInterpolationString) String() string {
	return fmt.Sprintf("[str:%s]", s.s)
}

func (s StringInterpolationString) StringToken() StringToken {
	return s.s
}

func (s StringInterpolationString) FetchPositionLength() SourceFileReference {
	return s.s.FetchPositionLength()
}
