/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package coloring

import (
	"github.com/swamp/compiler/src/ast"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type DecoratedColorer interface {
	Operator(t token.OperatorToken)
	CustomType(t *dectype.CustomTypeAtom)
	CustomTypeName(t *dectype.CustomTypeReference)
	NewLine(indentation int)
	OperatorString(s string)
	OneSpace()
	CustomTypeVariant(t *dectype.CustomTypeVariantAtom)
	RecordTypeField(t *dectype.RecordField)
	AliasName(t *dectype.Alias)
	KeywordString(s string)
	UnmanagedName(s *ast.UnmanagedType)
	PrimitiveTypeName(identifier *ast.TypeIdentifier)
	RightArrow()
	LocalTypeName(localType *dectype.LocalTypeDefinition)
}
