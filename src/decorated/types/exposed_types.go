/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type ExposedTypes struct {
	identifierToType map[string]dtype.Type
	types            []dtype.Type
}

func NewExposedTypes() *ExposedTypes {
	return &ExposedTypes{identifierToType: make(map[string]dtype.Type)}
}

func (e *ExposedTypes) AddType(t dtype.Type) {
	e.identifierToType[t.DecoratedName()] = t
	e.types = append(e.types, t)
}

func (e *ExposedTypes) AddTypes(allTypes []dtype.Type) {
	for _, t := range allTypes {
		e.AddType(t)
	}
}

func (e *ExposedTypes) AddTypesWithModulePrefix(allTypes []dtype.Type, prefix PackageRelativeModuleName) {
	for _, t := range allTypes {
		fakeVariable := ast.NewVariableIdentifier(token.NewVariableSymbolToken(t.DecoratedName(), nil, token.PositionLength{}, 0))
		fullyQualifiedName := prefix.JoinLocalName(fakeVariable)
		e.identifierToType[fullyQualifiedName] = t
		e.types = append(e.types, t)
	}
}

func (e *ExposedTypes) FindTypeFromSignature(complete string) dtype.Type {
	return e.identifierToType[complete]
}

func (e *ExposedTypes) AllExposedTypes() []dtype.Type {
	return e.types
}

func (e *ExposedTypes) ShortString() string {
	s := ""
	for _, exposedType := range e.types {
		s += fmt.Sprintf("%v : %v\n", exposedType.DecoratedName(), exposedType.DecoratedName())
	}

	return s
}

func (e *ExposedTypes) DebugString() string {
	s := ""
	for _, exposedType := range e.types {
		s += fmt.Sprintf("%p %v : %v\n", exposedType, exposedType.DecoratedName(), exposedType.DecoratedName())
	}

	return s
}
