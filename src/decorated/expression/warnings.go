package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type UnusedWarning struct {
	definition ModuleDef
}

func NewUnusedWarning(definition ModuleDef) *UnusedWarning {
	return &UnusedWarning{definition: definition}
}

func (e *UnusedWarning) Warning() string {
	return fmt.Sprintf("unused definition '%v'", e.definition.Identifier().Name())
}

func (e *UnusedWarning) FetchPositionLength() token.SourceFileReference {
	return e.definition.Identifier().FetchPositionLength()
}

type UnusedTypeWarning struct {
	unusedType dtype.Type
}

func NewUnusedTypeWarning(unusedType dtype.Type) *UnusedTypeWarning {
	return &UnusedTypeWarning{unusedType: unusedType}
}

func (e *UnusedTypeWarning) Warning() string {
	return fmt.Sprintf("unused type '%v'", e.unusedType.HumanReadable())
}

func (e *UnusedTypeWarning) FetchPositionLength() token.SourceFileReference {
	return e.unusedType.FetchPositionLength()
}

type UnusedImportWarning struct {
	definition  *ImportedModule
	description string
}

func NewUnusedImportWarning(definition *ImportedModule, description string) *UnusedImportWarning {
	if definition == nil {
		//		panic("must have definition")
	}
	return &UnusedImportWarning{definition: definition, description: description}
}

func (e *UnusedImportWarning) Warning() string {
	if e.definition == nil {
		return fmt.Sprintf("unused import (%v)", e.description)
	}
	return fmt.Sprintf("unused import %v (%v)", e.definition.ModuleName().ModuleName(), e.description)
}

func (e *UnusedImportWarning) FetchPositionLength() token.SourceFileReference {
	if e.definition == nil {
		return token.SourceFileReference{}
	}
	return e.definition.ModuleName().FetchPositionLength()
}

type UnusedImportStatementWarning struct {
	definition *ImportStatement
}

func NewUnusedImportStatementWarning(definition *ImportStatement) *UnusedImportStatementWarning {
	if definition.astImport == nil {
		panic("what is this")
	}
	return &UnusedImportStatementWarning{definition: definition}
}

func (e *UnusedImportStatementWarning) Warning() string {
	return fmt.Sprintf("unused import '%v'", e.definition.astImport.ModuleName().ModuleName())
}

func (e *UnusedImportStatementWarning) FetchPositionLength() token.SourceFileReference {
	return e.definition.FetchPositionLength()
}
