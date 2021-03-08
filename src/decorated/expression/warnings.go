package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type UnusedWarning struct {
	definition *ModuleDefinition
}

func NewUnusedWarning(definition *ModuleDefinition) *UnusedWarning {
	return &UnusedWarning{definition: definition}
}

func (e *UnusedWarning) Warning() string {
	return fmt.Sprintf("unused definition %v", e.definition.Identifier().Name())
}

func (e *UnusedWarning) FetchPositionLength() token.SourceFileReference {
	return e.definition.Identifier().FetchPositionLength()
}
