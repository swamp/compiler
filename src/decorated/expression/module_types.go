/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type ModuleTypes struct {
	identifierToType map[string]dtype.Type
	allNamedTypes    []NamedType
	sourceModule     *Module
}

func NewModuleTypes(sourceModule *Module) *ModuleTypes {
	t := &ModuleTypes{sourceModule: sourceModule, identifierToType: make(map[string]dtype.Type)}
	return t
}

func (t *ModuleTypes) AllInOrderTypes() []NamedType {
	return t.allNamedTypes
}

func (t *ModuleTypes) SourceModule() *Module {
	return t.sourceModule
}

// -----------------------------------------------------
//
//	Find
//
// -----------------------------------------------------
func (t *ModuleTypes) FindType(identifier *ast.TypeIdentifier) dtype.Type {
	found := t.identifierToType[identifier.Name()]
	return found
}

func (t *ModuleTypes) FindBuiltInType(name string) dtype.Type {
	found := t.identifierToType[name]
	return found
}

// -----------------------------------------------------
//                    Declare
// -----------------------------------------------------

func (t *ModuleTypes) InternalAddType(typeIdentifier *ast.TypeIdentifier, realType dtype.Type) error {
	existingType := t.FindType(typeIdentifier)
	if existingType != nil {
		return fmt.Errorf("sorry, '%v' already declared", existingType)
	}
	t.internalAdd(typeIdentifier, realType)
	return nil
}

func (t *ModuleTypes) InternalAddPrimitive(typeIdentifier *ast.TypeIdentifier, atom *dectype.PrimitiveAtom) error {
	return t.InternalAddType(typeIdentifier, atom)
}

type TypeError interface {
	Error() string
	FetchPositionLength() token.SourceFileReference
}

/*


 */

type UnknownType struct {
	sourceFileReference token.SourceFileReference
	errString           string
}

func (c *UnknownType) Error() string {
	return c.errString
}

func (c *UnknownType) FetchPositionLength() token.SourceFileReference {
	return c.sourceFileReference
}

func (t *ModuleTypes) AddTypeAlias(alias *dectype.Alias) TypeError {
	if err := t.InternalAddType(alias.TypeIdentifier(), alias); err != nil {
		return &UnknownType{sourceFileReference: alias.AstAlias().FetchPositionLength(), errString: err.Error()}
	}

	return nil
}

func (t *ModuleTypes) addCustomTypeVariantConstructors(customType *dectype.CustomTypeAtom) {
	for _, variant := range customType.Variants() {
		t.InternalAddType(variant.Name(), variant)
	}
}

func (t *ModuleTypes) AddCustomType(customType *dectype.CustomTypeAtom) TypeError {
	t.InternalAddType(customType.TypeIdentifier(), customType)
	t.addCustomTypeVariantConstructors(customType)
	return nil
}

func (t *ModuleTypes) AddCustomTypeWrappedInNameOnlyContext(customTypeWrappedInContext *dectype.LocalTypeNameOnlyContext) TypeError {
	customType, _ := customTypeWrappedInContext.Next().(*dectype.CustomTypeAtom)
	t.InternalAddType(customType.TypeIdentifier(), customTypeWrappedInContext)
	t.addCustomTypeVariantConstructors(customType)
	return nil
}

// -----------------------------------------------------
//                    Other
// -----------------------------------------------------

func (t *ModuleTypes) DebugOutput() {
	fmt.Println(t.DebugString())
}

func TraverseToString(t dtype.Type) string {
	s := t.String()
	concreteType, isConcreteType := t.(dtype.Type)
	if isConcreteType {
		next := concreteType.Next()
		if next != nil {
			s += " => " + TraverseToString(next)
		}
	}
	return s
}

func (t *ModuleTypes) DebugString() string {
	s := ""
	for _, v := range t.allNamedTypes {
		s += fmt.Sprintf("%v : %v\n", v.name, v.realType)
	}
	return s
}

func (t *ModuleTypes) String() string {
	return t.DebugString()
}

func (t *ModuleTypes) internalAdd(identifier *ast.TypeIdentifier, realType dtype.Type) {
	t.internalAddWithString(identifier.Name(), realType)
}

type NamedType struct {
	name     string
	realType dtype.Type
}

func (n NamedType) Name() string {
	return n.name
}

func (n NamedType) RealType() dtype.Type {
	return n.realType
}

func (t *ModuleTypes) internalAddWithString(name string, realType dtype.Type) {
	hasType := t.identifierToType[name]
	if hasType != nil {
		panic("already have name " + name)
	}
	t.identifierToType[name] = realType
	t.allNamedTypes = append(t.allNamedTypes, NamedType{name: name, realType: realType})
}

func (t *ModuleTypes) CopyTypes(realTypes []NamedType) decshared.DecoratedError {
	for _, copyType := range realTypes {
		symbol := token.NewTypeSymbolToken(copyType.name, t.sourceModule.FetchPositionLength(), 0)
		fakeIdentifier := ast.NewTypeIdentifier(symbol)
		copyErr := t.CopyType(fakeIdentifier, copyType.realType)
		if copyErr != nil {
			return copyErr
		}
	}
	return nil
}

func (t *ModuleTypes) CopyType(nameOfType *ast.TypeIdentifier, realType dtype.Type) decshared.DecoratedError {
	existingType := t.FindType(nameOfType)
	if existingType != nil {
		return NewInternalError(fmt.Errorf("copy: sorry, '%v' already declared", existingType))
	}

	t.internalAddWithString(nameOfType.Name(), realType)
	return nil
}
