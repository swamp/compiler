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
	SourceFileReference
	value bool `debug:"true"`
	raw   string
}

func NewBooleanToken(raw string, v bool, sourceFileReference SourceFileReference) BooleanToken {
	return BooleanToken{raw: raw, value: v, SourceFileReference: sourceFileReference}
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

func (s BooleanToken) FetchPositionLength() SourceFileReference {
	return s.SourceFileReference
}
