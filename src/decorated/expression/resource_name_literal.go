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

type ResourceNameLiteral struct {
	resourceName           *ast.ResourceNameLiteral
	globalFixedType dtype.Type
}

func NewResourceNameLiteral(resourceName *ast.ResourceNameLiteral, globalFixedType dtype.Type) *ResourceNameLiteral {
	return &ResourceNameLiteral{resourceName: resourceName, globalFixedType: globalFixedType}
}

func (i *ResourceNameLiteral) Type() dtype.Type {
	return i.globalFixedType
}

func (i *ResourceNameLiteral) Value() string {
	return i.resourceName.Value()
}

func (i *ResourceNameLiteral) String() string {
	return fmt.Sprintf("[resource name %v]", i.resourceName.Value())
}

func (i *ResourceNameLiteral) FetchPositionAndLength() token.PositionLength {
	return i.resourceName.Token.FetchPositionLength()
}
