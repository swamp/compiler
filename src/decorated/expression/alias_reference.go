package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type AliasReference struct {
	castToType *dectype.AliasReference
}

func NewAliasReference(castToType *dectype.AliasReference) *AliasReference {
	return &AliasReference{castToType: castToType}
}

func (c *AliasReference) FetchPositionLength() token.SourceFileReference {
	return c.castToType.FetchPositionLength()
}

func (c *AliasReference) String() string {
	return fmt.Sprintf("cast %v %v", c.castToType, c.castToType)
}

func (c *AliasReference) Type() dtype.Type {
	return c.castToType.Type()
}
