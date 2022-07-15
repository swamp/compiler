package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type AliasReference struct {
	referencedType      *dectype.AliasReference
	definitionReference *dectype.NamedDefinitionTypeReference
}

func NewAliasReference(definitionReference *dectype.NamedDefinitionTypeReference, referencedType *dectype.AliasReference) *AliasReference {
	return &AliasReference{definitionReference: definitionReference, referencedType: referencedType}
}

func (c *AliasReference) FetchPositionLength() token.SourceFileReference {
	return c.referencedType.FetchPositionLength()
}

func (c *AliasReference) String() string {
	return fmt.Sprintf("[AliasRef %v]", c.referencedType)
}

func (c *AliasReference) Type() dtype.Type {
	return c.referencedType
}

func (c *AliasReference) TypeAliasReference() *dectype.AliasReference {
	return c.referencedType
}

func (c *AliasReference) NameReference() *dectype.NamedDefinitionTypeReference {
	return c.definitionReference
}

func (c *AliasReference) HumanReadable() string {
	return c.referencedType.HumanReadable()
}
