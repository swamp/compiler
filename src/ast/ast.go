/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"github.com/swamp/compiler/src/token"
)

type Node interface {
	FetchPositionLength() token.Range
	String() string
}

type Statement interface {
	Node
}

type Expression interface {
	Node
	DebugString() string
}

type Literal = Expression
