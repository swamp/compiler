/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	dectype "github.com/swamp/compiler/src/decorated/types"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type ResourceNameLiteral struct {
	resourceName    *ast.ResourceNameLiteral `debug:"true"`
	globalFixedType *dectype.PrimitiveTypeReference
}

func NewResourceNameLiteral(resourceName *ast.ResourceNameLiteral, globalFixedType *dectype.PrimitiveTypeReference) *ResourceNameLiteral {
	return &ResourceNameLiteral{resourceName: resourceName, globalFixedType: globalFixedType}
}

func (i *ResourceNameLiteral) Type() dtype.Type {
	return i.globalFixedType
}

func (i *ResourceNameLiteral) Value() string {
	return i.resourceName.Value()
}

func (i *ResourceNameLiteral) String() string {
	return fmt.Sprintf("[ResourceName %v]", i.resourceName.Value())
}

func (i *ResourceNameLiteral) HumanReadable() string {
	return "Resource Name"
}

func (i *ResourceNameLiteral) FetchPositionLength() token.SourceFileReference {
	return i.resourceName.Token.FetchPositionLength()
}
