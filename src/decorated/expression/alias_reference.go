package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type AliasReference struct {
	referencedType *dectype.AliasReference
}

func NewAliasReference(referencedType *dectype.AliasReference) *AliasReference {
	return &AliasReference{referencedType: referencedType}
}

func (c *AliasReference) FetchPositionLength() token.SourceFileReference {
	return c.referencedType.FetchPositionLength()
}

func (c *AliasReference) String() string {
	return fmt.Sprintf("[AliasRef %v]", c.referencedType)
}

func (c *AliasReference) Type() dtype.Type {
	return c.referencedType.Type()
}
