/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
)

type ModuleDeclarations struct {
	types map[string]dtype.Type

	ownedByModule *Module
}

func NewModuleDeclarations(ownedByModule *Module) *ModuleDeclarations {
	if ownedByModule == nil {
		panic("sorry, all localDefinitions must be owned by a module")
	}
	return &ModuleDeclarations{
		ownedByModule: ownedByModule, types: make(map[string]dtype.Type),
	}
}

func (d *ModuleDeclarations) AddDeclaration(identifier *ast.VariableIdentifier, declaredType dtype.Type) error {
	declaration, wasFound := d.types[identifier.Name()]
	if wasFound {
		return fmt.Errorf("already have type %v %v", identifier, declaration)
	}
	d.types[identifier.Name()] = declaredType

	return nil
}

func (t *ModuleDeclarations) String() string {
	s := ""

	for _, declaredType := range t.types {
		s += fmt.Sprintf(" %v \n", declaredType)
	}

	return s
}
