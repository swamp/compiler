/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

// TypeId :
type TypeId struct {
	Range
	raw         string
	Indentation int
}

func NewTypeId(raw string, startPosition Range, indentation int) TypeId {
	return TypeId{raw: raw, Range: startPosition, Indentation: indentation}
}

func (s TypeId) Type() Type {
	return TypeIdSymbol
}

func (s TypeId) Name() string {
	return s.raw
}

func (s TypeId) Raw() string {
	return s.raw
}

func (s TypeId) FetchIndentation() int {
	return s.Indentation
}

func (s TypeId) String() string {
	return "@"
}
