/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

// SpaceToken :
type SpaceToken struct {
	Range
	r rune
}

func NewSpaceToken(position Range, r rune) SpaceToken {
	return SpaceToken{Range: position, r: r}
}

func (s SpaceToken) String() string {
	return "space"
}

func (s SpaceToken) Type() Type {
	return Space
}
