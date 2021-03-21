/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type ExposedTypes struct {
	identifierToType map[string]dtype.Type
}

func NewExposedTypes(module *Module) *ExposedTypes {
	return &ExposedTypes{identifierToType: make(map[string]dtype.Type)}
}

func (e *ExposedTypes) internalAddType(name string, t dtype.Type) {
	e.identifierToType[name] = t
}

func (e *ExposedTypes) AddTypes(allTypes map[string]dtype.Type) {
	for name, t := range allTypes {
		e.internalAddType(name, t)
	}
}

func (e *ExposedTypes) AddTypesWithModulePrefix(allTypes map[string]dtype.Type, prefix dectype.PackageRelativeModuleName) {
	for name, t := range allTypes {
		fakeVariable := ast.NewVariableIdentifier(token.NewVariableSymbolToken(name, token.SourceFileReference{}, 0))
		fullyQualifiedName := prefix.JoinLocalName(fakeVariable)
		e.identifierToType[fullyQualifiedName] = t
	}
}

func (t *ExposedTypes) AddBuiltInTypes(name *ast.TypeIdentifier, referencedType dtype.Type, localComments []ast.LocalComment) TypeError {
	t.internalAddType(name.Name(), referencedType)

	return nil
}

func (e *ExposedTypes) FindType(complete *ast.TypeIdentifier) dtype.Type {
	return e.identifierToType[complete.Name()]
}

func (e *ExposedTypes) FindBuiltInType(s string) dtype.Type {
	return e.identifierToType[s]
}

func (e *ExposedTypes) AllTypes() map[string]dtype.Type {
	return e.identifierToType
}

func (e *ExposedTypes) DebugString() string {
	s := ""
	for _, exposedType := range e.identifierToType {
		s += fmt.Sprintf("%p %v : %v\n", exposedType, exposedType.String(), exposedType.String())
	}

	return s
}
