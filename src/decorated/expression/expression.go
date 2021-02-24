/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type Node interface {
	FetchPositionLength() token.Range
	String() string
}

type DecoratedExpression interface {
	Node
	Type() dtype.Type
}

type DecoratedExpressionNode struct {
	decoratedType dtype.Type
}

func (d DecoratedExpressionNode) Type() dtype.Type {
	return d.decoratedType
}
