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

type Import struct {
	astImport       *ast.Import
	moduleReference *Module
}

func NewImport(astImport *ast.Import, moduleReference *Module) *Import {
	return &Import{astImport: astImport, moduleReference: moduleReference}
}

func (l *Import) String() string {
	return fmt.Sprintf("[import %v %v]", l.astImport, l.moduleReference)
}

func (l *Import) Module() *Module {
	return l.moduleReference
}

func (l *Import) AstImport() *ast.Import {
	return l.astImport
}

func (l *Import) FetchPositionLength() token.SourceFileReference {
	return l.astImport.FetchPositionLength()
}
