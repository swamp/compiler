/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import (
	"fmt"
)

type SourceFileURI struct {
	name string
}

func MakeSourceFileURI(name string) *SourceFileURI {
	return &SourceFileURI{name: name}
}

func (s *SourceFileURI) String() string {
	return s.name
}

func (s *SourceFileURI) ReferenceString() string {
	return fmt.Sprintf("%v:", s.name)
}

func (s *SourceFileURI) ReferenceWithPositionString(pos Position) string {
	return fmt.Sprintf("%v:%d:%d:", s.name, pos.line+1, pos.column+1)
}
