/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import "fmt"

// SpaceToken :
type SpaceToken struct {
	PositionLength
	r rune
}

func NewSpaceToken(position PositionLength, r rune) SpaceToken {
	return SpaceToken{PositionLength: position, r: r}
}

func (s SpaceToken) String() string {
	return fmt.Sprintf("space")
}

func (s SpaceToken) Type() Type {
	return Space
}
