/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package modref

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/token"
)

type ModuleReferencer interface {
	FetchPositionLength() token.SourceFileReference
	AstModuleReference() *ast.ModuleReference
	String() string
}
