/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

// TypeId :
type TypeId struct {
	SourceFileReference
	raw         string
	Indentation int
}

func NewTypeId(raw string, startPosition SourceFileReference, indentation int) TypeId {
	return TypeId{raw: raw, SourceFileReference: startPosition, Indentation: indentation}
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

func (s TypeId) FetchPositionLength() SourceFileReference {
	return s.SourceFileReference
}
