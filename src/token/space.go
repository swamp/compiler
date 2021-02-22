/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

// SpaceToken :
type SpaceToken struct {
	SourceFileReference
	r rune
}

func NewSpaceToken(position SourceFileReference, r rune) SpaceToken {
	return SpaceToken{SourceFileReference: position, r: r}
}

func (s SpaceToken) String() string {
	return "space"
}

func (s SpaceToken) Type() Type {
	return Space
}

func (s SpaceToken) FetchPositionLength() SourceFileReference {
	return s.SourceFileReference
}
