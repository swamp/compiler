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
	FetchPositionLength() token.SourceFileReference
	String() string
}

type Expression interface {
	Node
	Type() dtype.Type
}

type ExpressionNode struct {
	decoratedType dtype.Type
}

func (d ExpressionNode) Type() dtype.Type {
	return d.decoratedType
}
