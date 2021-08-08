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
	CustomTypeVariant(t *dectype.CustomTypeVariant)
	InvokerType(t *dectype.InvokerType)
	RecordTypeField(t *dectype.RecordField)
	AliasName(t *dectype.Alias)
	KeywordString(s string)
	UnmanagedName(s *ast.UnmanagedType)
	PrimitiveTypeName(identifier *ast.TypeIdentifier)
	RightArrow()
	LocalTypeName(localType *dectype.LocalType)
}
