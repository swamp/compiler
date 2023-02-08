/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type CustomTypeNamedDefinition struct {
	customType Type
}

func NewCustomTypeNamedDefinition(customType Type) *CustomTypeNamedDefinition {
	return &CustomTypeNamedDefinition{customType: customType}
}

func (i *CustomTypeNamedDefinition) CustomType() Type {
	return i.customType
}

func (i *CustomTypeNamedDefinition) FetchPositionLength() token.SourceFileReference {
	return i.customType.FetchPositionLength()
}

func (i *CustomTypeNamedDefinition) String() string {
	return fmt.Sprintf("[CustomTypeStatement %v]", i.customType)
}

func (i *CustomTypeNamedDefinition) DebugString() string {
	return "[customtypedefinition]"
}
