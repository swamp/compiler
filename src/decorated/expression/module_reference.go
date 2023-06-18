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

type ModuleReference struct {
	module     *Module
	identifier *ast.ModuleReference
	inclusive  token.SourceFileReference
}

func NewModuleReference(identifier *ast.ModuleReference, module *Module) *ModuleReference {
	inclusive := token.MakeInclusiveSourceFileReference(identifier.First().FetchPositionLength(),
		identifier.Last().FetchPositionLength())
	if !inclusive.Verify() {
		//panic(fmt.Errorf("wrong here"))
	}
	ref := &ModuleReference{module: module, identifier: identifier, inclusive: inclusive}

	module.AddReference(ref)

	return ref
}

func (m *ModuleReference) String() string {
	return fmt.Sprintf("moduleref %v", m.module)
}

func (m *ModuleReference) Module() *Module {
	return m.module
}

func (m *ModuleReference) AstModuleReference() *ast.ModuleReference {
	return m.identifier
}

func (m *ModuleReference) HumanReadable() string {
	return "Module Reference"
}

func (m *ModuleReference) FetchPositionLength() token.SourceFileReference {
	return m.inclusive
}
