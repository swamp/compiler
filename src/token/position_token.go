/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import "fmt"

// PositionToken :
type PositionToken struct {
	position    Position
	indentation int
}

func NewPositionToken(position Position, indentation int) PositionToken {
	return PositionToken{position: position, indentation: indentation}
}

func (p PositionToken) FetchPositionToken() PositionToken {
	return p
}

func (p PositionToken) Position() Position {
	return p.position
}

func (p PositionToken) Indentation() int {
	return p.indentation
}

func (p PositionToken) String() string {
	return fmt.Sprintf("[%d:%d  (%d)] ", p.position.Line(), p.position.Column(), p.indentation)
}
