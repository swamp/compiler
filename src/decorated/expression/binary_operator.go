/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type BinaryOperator struct {
	ExpressionNode `debug:"true"`
	left           Expression `debug:"true"`
	right          Expression `debug:"true"`
}

func (b *BinaryOperator) Left() Expression {
	return b.left
}

func (b *BinaryOperator) Right() Expression {
	return b.right
}

func (b *BinaryOperator) FetchPositionLength() token.SourceFileReference {
	inclusive := token.MakeInclusiveSourceFileReference(b.left.FetchPositionLength(), b.right.FetchPositionLength())
	return inclusive
}

func (b *BinaryOperator) String() string {
	return fmt.Sprintf("binary operator %v %v", b.left, b.right)
}
