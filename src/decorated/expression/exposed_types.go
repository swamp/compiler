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

type ImportedType struct {
	referencedType dtype.Type
	createdBy      *ImportedModule
	name           string
	wasReferenced  bool
}

func (i *ImportedType) String() string {
	return i.referencedType.String()
}

func (i *ImportedType) MarkAsReferenced() {
	i.wasReferenced = true
	if i.createdBy != nil {
		i.createdBy.MarkAsReferenced()
	}
}

func (i *ImportedType) WasReferenced() bool {
	return i.wasReferenced
}

func (i *ImportedType) CreatedByModuleImport() *ImportedModule {
	return i.createdBy
}

func (i *ImportedType) ReferencedType() dtype.Type {
	return i.referencedType
}

type ExposedTypes struct {
	identifierToType map[string]*ImportedType
	allInOrderTypes  []*ImportedType
}

func NewExposedTypes(module *Module) *ExposedTypes {
	return &ExposedTypes{identifierToType: make(map[string]*ImportedType)}
}

func (e *ExposedTypes) internalAddType(name string, t dtype.Type, importedModule *ImportedModule) {
	importedType := &ImportedType{name: name, referencedType: t, createdBy: importedModule}
	e.identifierToType[name] = importedType
	e.allInOrderTypes = append(e.allInOrderTypes, importedType)
}

func (e *ExposedTypes) AddTypes(allTypes []*ImportedType, importedModule *ImportedModule) {
	for _, t := range allTypes {
		e.internalAddType(t.name, t.referencedType, importedModule)
	}
}

func (e *ExposedTypes) AddTypesFromModule(allTypes []NamedType, module *Module) {
	for _, t := range allTypes {
		e.internalAddType(t.name, t.realType, nil)
	}
}

func (e *ExposedTypes) AddTypesWithModulePrefix(allTypes map[string]dtype.Type, prefix dectype.PackageRelativeModuleName) {
	for name, t := range allTypes {
		fakeVariable := ast.NewVariableIdentifier(token.NewVariableSymbolToken(name, token.SourceFileReference{}, 0))
		fullyQualifiedName := prefix.JoinLocalName(fakeVariable)
		e.identifierToType[fullyQualifiedName] = &ImportedType{referencedType: t}
	}
}

func (t *ExposedTypes) AddBuiltInTypes(name *ast.TypeIdentifier, referencedType dtype.Type, localComments []ast.LocalComment) TypeError {
	t.internalAddType(name.Name(), referencedType, nil)

	return nil
}

func (e *ExposedTypes) FindType(complete *ast.TypeIdentifier) *ImportedType {
	return e.identifierToType[complete.Name()]
}

func (e *ExposedTypes) FindBuiltInType(s string) dtype.Type {
	return e.identifierToType[s].referencedType
}

func (e *ExposedTypes) AllInOrderTypes() []*ImportedType {
	return e.allInOrderTypes
}

func (e *ExposedTypes) DebugString() string {
	s := ""
	for _, exposedType := range e.identifierToType {
		s += fmt.Sprintf("%p %v : %v\n", exposedType, exposedType.String(), exposedType.String())
	}

	return s
}
