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
}

func NewImport(astImport *ast.Import, moduleReference *ModuleReference, alias *ModuleReference) *ImportStatement {
	return &ImportStatement{astImport: astImport, moduleReference: moduleReference, alias: alias}
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

func (l *ImportStatement) Alias() *ModuleReference {
	return l.alias
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
