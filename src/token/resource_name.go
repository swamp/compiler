/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import "fmt"

// ResourceName :
type ResourceName struct {
	Range
	raw         string
	Indentation int
}

func NewResourceName(raw string, startPosition Range, indentation int) ResourceName {
	return ResourceName{raw: raw, Range: startPosition, Indentation: indentation}
}

func (s ResourceName) Type() Type {
	return ResourceNameSymbol
}

func (s ResourceName) Name() string {
	return s.raw
}

func (s ResourceName) Raw() string {
	return s.raw
}

func (s ResourceName) FetchIndentation() int {
	return s.Indentation
}

func (s ResourceName) String() string {
	return fmt.Sprintf("@%s", s.raw)
}
