/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type CustomTypeStatement struct {
	customTypeOrWrapped Type
	name                *TypeIdentifier
	precedingComments   token.CommentBlock
}

func NewCustomTypeStatement(name *TypeIdentifier, customType Type,
	precedingComments token.CommentBlock) *CustomTypeStatement {
	return &CustomTypeStatement{name: name, customTypeOrWrapped: customType, precedingComments: precedingComments}
}

func (i *CustomTypeStatement) CommentBlock() token.CommentBlock {
	return i.precedingComments
}

func (i *CustomTypeStatement) String() string {
	return fmt.Sprintf("[custom-type-statement %v]", i.customTypeOrWrapped)
}

func (i *CustomTypeStatement) TypeIdentifier() *TypeIdentifier {
	return i.name
}

func (i *CustomTypeStatement) Type() Type {
	return i.customTypeOrWrapped
}

func (i *CustomTypeStatement) Name() string {
	return i.customTypeOrWrapped.Name()
}

func (i *CustomTypeStatement) PositionLength() token.PositionLength {
	return i.name.symbolToken.FetchPositionLength()
}

func (i *CustomTypeStatement) DebugString() string {
	return fmt.Sprintf("[custom-type-statement %v]", i.customTypeOrWrapped)
}
