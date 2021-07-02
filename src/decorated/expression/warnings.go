package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type UnusedWarning struct {
	definition ModuleDef
}

func NewUnusedWarning(definition ModuleDef) *UnusedWarning {
	return &UnusedWarning{definition: definition}
}

func (e *UnusedWarning) Warning() string {
	return fmt.Sprintf("unused definition %v", e.definition.Identifier().Name())
}

func (e *UnusedWarning) FetchPositionLength() token.SourceFileReference {
	return e.definition.Identifier().FetchPositionLength()
}

type UnusedImportWarning struct {
	definition  *ImportedModule
	description string
}

func NewUnusedImportWarning(definition *ImportedModule, description string) *UnusedImportWarning {
	return &UnusedImportWarning{definition: definition, description: description}
}

func (e *UnusedImportWarning) Warning() string {
	return fmt.Sprintf("unused import %v (%v)", e.definition.ModuleName().ModuleName(), e.description)
}

func (e *UnusedImportWarning) FetchPositionLength() token.SourceFileReference {
	return e.definition.ModuleName().FetchPositionLength()
}
