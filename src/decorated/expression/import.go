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
	moduleReference *Module
}

func NewImport(astImport *ast.Import, moduleReference *Module) *ImportStatement {
	return &ImportStatement{astImport: astImport, moduleReference: moduleReference}
}

func (l *ImportStatement) String() string {
	return fmt.Sprintf("[import %v %v]", l.astImport, l.moduleReference)
}

func (l *ImportStatement) StatementString() string {
	return fmt.Sprintf("[import %v %v]", l.astImport, l.moduleReference)
}

func (l *ImportStatement) Module() *Module {
	return l.moduleReference
}

func (l *ImportStatement) AstImport() *ast.Import {
	return l.astImport
}

func (l *ImportStatement) FetchPositionLength() token.SourceFileReference {
	return l.astImport.FetchPositionLength()
}
