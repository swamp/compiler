/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"github.com/swamp/compiler/src/token"
)

type UnaryOperator struct {
	ExpressionNode
	left Expression
}

func (u *UnaryOperator) String() string {
	return "unary"
}

func (u *UnaryOperator) FetchPositionLength() token.SourceFileReference {
	return u.left.FetchPositionLength()
}

func (u *UnaryOperator) Left() Expression {
	return u.left
}
