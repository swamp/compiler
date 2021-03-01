/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type BooleanLiteral struct {
	boolean           *ast.BooleanLiteral
	globalBooleanType dtype.Type
}

func NewBooleanLiteral(boolean *ast.BooleanLiteral, globalBooleanType dtype.Type) *BooleanLiteral {
	return &BooleanLiteral{boolean: boolean, globalBooleanType: globalBooleanType}
}

func (i *BooleanLiteral) Type() dtype.Type {
	return i.globalBooleanType
}

func (i *BooleanLiteral) Value() bool {
	return i.boolean.Value()
}

func (i *BooleanLiteral) String() string {
	return fmt.Sprintf("[bool %v]", i.boolean.Value())
}

func (i *BooleanLiteral) FetchPositionLength() token.SourceFileReference {
	return i.boolean.FetchPositionLength()
}
