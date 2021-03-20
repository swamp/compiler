package modref

import (
	"github.com/swamp/compiler/src/token"
)

type ModuleReferencer interface {
	FetchPositionLength() token.SourceFileReference
	String() string
}
