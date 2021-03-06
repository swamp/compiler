/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type CustomTypeVariantReference struct {
	named             *NamedDefinitionTypeReference
	customTypeVariant *dectype.CustomTypeVariant
}

func (g *CustomTypeVariantReference) Type() dtype.Type {
	return g.customTypeVariant
}

func (g *CustomTypeVariantReference) String() string {
	return fmt.Sprintf("[customtypevariantref %v %v]", g.named, g.customTypeVariant)
}

func (g *CustomTypeVariantReference) HumanReadable() string {
	return "Variant Reference"
}

func (g *CustomTypeVariantReference) CustomTypeVariant() *dectype.CustomTypeVariant {
	return g.customTypeVariant
}

func (g *CustomTypeVariantReference) AstIdentifier() ast.TypeIdentifierNormalOrScoped {
	return g.named.ident
}

func NewCustomTypeVariantReference(named *NamedDefinitionTypeReference, customTypeVariant *dectype.CustomTypeVariant) *CustomTypeVariantReference {
	ref := &CustomTypeVariantReference{named: named, customTypeVariant: customTypeVariant}

	return ref
}

func (g *CustomTypeVariantReference) FetchPositionLength() token.SourceFileReference {
	return g.named.FetchPositionLength()
}
