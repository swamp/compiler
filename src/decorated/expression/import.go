/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/token"
)

type ImportStatement struct {
	astImport       *ast.Import
	moduleReference *ModuleReference
	alias           *ModuleReference
	references      []*ImportStatementReference
	isReferenced    bool
	exposeAll       bool
}

// source *decorated.Module, sourceMountedModuleName dectype.PackageRelativeModuleName, exposeAll bool

func NewImport(astImport *ast.Import, moduleReference *ModuleReference, alias *ModuleReference,
	exposeAll bool) *ImportStatement {
	if astImport.FetchPositionLength().Document == nil {
		panic("astImport is wrong")
	}
	return &ImportStatement{astImport: astImport, moduleReference: moduleReference, alias: alias, exposeAll: exposeAll}
}

func (l *ImportStatement) String() string {
	return fmt.Sprintf("[import %v %v]", l.astImport, l.moduleReference)
}

func (l *ImportStatement) StatementString() string {
	return fmt.Sprintf("[import %v %v]", l.astImport, l.moduleReference)
}

func (l *ImportStatement) Module() *Module {
	return l.moduleReference.module
}

func (l *ImportStatement) ExposeAll() bool {
	return l.exposeAll
}

func (l *ImportStatement) Alias() *ModuleReference {
	return l.alias
}

func (l *ImportStatement) ImportAsName() *ModuleReference {
	if l.alias != nil {
		return l.alias
	}

	return l.moduleReference
}

func (l *ImportStatement) ModuleReference() *ModuleReference {
	return l.moduleReference
}

func (l *ImportStatement) AstImport() *ast.Import {
	return l.astImport
}

func (l *ImportStatement) FetchPositionLength() token.SourceFileReference {
	return l.astImport.FetchPositionLength()
}

func (l *ImportStatement) MarkAsReferenced() {
	l.isReferenced = true
}

func (l *ImportStatement) AddReference(statementReference *ImportStatementReference) {
	l.references = append(l.references, statementReference)
}

func (l *ImportStatement) WasReferenced() bool {
	return len(l.references) > 0 || l.isReferenced
}
