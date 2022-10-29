/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type EmptyExpression struct {
	identifier *VariableIdentifier
}

func NewEmptyExpression(identifier *VariableIdentifier) *EmptyExpression {
	return &EmptyExpression{identifier: identifier}
}

func (i *EmptyExpression) FetchPositionLength() token.SourceFileReference {
	return i.identifier.FetchPositionLength()
}

func (i *EmptyExpression) String() string {
	return fmt.Sprintf("[EmptyExpression]")
}

func (i *EmptyExpression) DebugString() string {
	return "[EmptyExpression]"
}
