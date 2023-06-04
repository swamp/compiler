/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"
	"log"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type ResolvedLocalTypeContext struct {
	resolvedArguments                map[string]*ResolvedLocalType
	definitions                      []*ResolvedLocalType               `debug:"true"`
	contextRefThatWantsResolvedTypes *LocalTypeNameOnlyContextReference `debug:"true"`
}

func (t *ResolvedLocalTypeContext) DeclareString() string {
	if len(t.resolvedArguments) == 0 {
		return ""
	}
	return fmt.Sprintf("%v", t.resolvedArguments)
}

func (t *ResolvedLocalTypeContext) String() string {
	return fmt.Sprintf("[ResolvedContext %v => %v]", t.DebugString(), t.contextRefThatWantsResolvedTypes)
}

func (t *ResolvedLocalTypeContext) DebugString() string {
	s := ""
	for _, argumentType := range t.definitions {
		if s != "" {
			s += ", "
		}
		s += fmt.Sprintf("%v", argumentType)

	}
	return s
}

func (t *ResolvedLocalTypeContext) Resolve() (dtype.Atom, error) {
	resolvedType, resolveErr := Collapse(t.contextRefThatWantsResolvedTypes.nameContext.Next(), t)
	if resolveErr != nil {
		return nil, resolveErr
	}

	return resolvedType.Resolve()
}

func (t *ResolvedLocalTypeContext) HumanReadable() string {
	return "resolved context"
}

func (t *ResolvedLocalTypeContext) Next() dtype.Type {
	return t.contextRefThatWantsResolvedTypes
}

func (t *ResolvedLocalTypeContext) WasReferenced() bool {
	return false
}

func (t *ResolvedLocalTypeContext) FetchPositionLength() token.SourceFileReference {
	return t.contextRefThatWantsResolvedTypes.FetchPositionLength()
}

func NewResolvedLocalTypeContext(contextRefThatWantsResolvedTypes *LocalTypeNameOnlyContextReference,
	types []dtype.Type) (
	*ResolvedLocalTypeContext, error,
) {
	t := &ResolvedLocalTypeContext{
		resolvedArguments:                make(map[string]*ResolvedLocalType),
		contextRefThatWantsResolvedTypes: contextRefThatWantsResolvedTypes,
	}

	if len(contextRefThatWantsResolvedTypes.nameContext.Definitions()) != len(types) {
		return nil, fmt.Errorf("must have same number of types as names")
	}

	for index, resolvedType := range types {
		foundName := contextRefThatWantsResolvedTypes.nameContext.Definitions()[index]
		log.Printf("resolved '%s' <- %T %v", foundName.identifier.Name(), resolvedType, resolvedType)
		newLocalTypeDef := NewResolvedLocalType(foundName, resolvedType)
		t.resolvedArguments[foundName.identifier.Name()] = newLocalTypeDef
		t.definitions = append(t.definitions, newLocalTypeDef)
	}

	return t, nil
}

func (t *ResolvedLocalTypeContext) Definitions() []*ResolvedLocalType {
	return t.definitions
}

/*
	func (t *ResolvedLocalTypeContext) SetTypes(types []dtype.Type) error {
		if len(types) != len(t.localTypeNames) {
			return fmt.Errorf("wrong number of localTypeNames")
		}

		for index, typeToSet := range types {
			def := t.localTypeNames[index]
			_, wasPrimitive := typeToSet.(*PrimitiveAtom)
			if wasPrimitive {
				panic(fmt.Errorf("should not have primitive atoms, I want type references"))
			}
			if err := def.SetDefinition(typeToSet); err != nil {
				return err
			}
		}

		return nil
	}
*/
func (t *ResolvedLocalTypeContext) Debug() {
	for _, def := range t.definitions {
		log.Printf("%v +> %T", def.Identifier(), def.ReferencedType())
	}
}

/*
func (t *ResolvedLocalTypeContext) AddExpectedDef(name *dtype.LocalTypeName) {
	nameDef := NewLocalTypeName(name)
	def := NewResolvedLocalType(nameDef)
	t.localTypeNamesMap[name.Name()] = def
	t.localTypeNames = append(t.localTypeNames, def)
}

func (t *ResolvedLocalTypeContext) AddExpectedDefs(names []*dtype.LocalTypeName) {
	for _, name := range names {
		t.AddExpectedDef(name)
	}
}



*/

func (t *ResolvedLocalTypeContext) HasDefinitions() bool {
	return len(t.definitions) > 0
}

/*
func (t *ResolvedLocalTypeContext) ResolveTypeRef(defRef *ResolvedLocalTypeReference) (*ResolvedLocalTypeReference, error) {
	definition, found := t.localTypeNamesMap[defRef.Identifier().Name()]
	if !found {
		return nil, fmt.Errorf("could not find %v", defRef.Identifier().Name())
	}

	return NewResolvedLocalTypeReference(definition, definition), nil
}
*/

func (t *ResolvedLocalTypeContext) LookupTypeAstRef(astReference *ast.LocalTypeNameReference) (
	*ResolvedLocalTypeReference, decshared.DecoratedError,
) {
	definition, found := t.resolvedArguments[astReference.Name()]
	if !found {
		return nil, NewCouldNotFindLocalTypeName(astReference, fmt.Errorf("could not find %v", astReference.Name()))
	}

	log.Printf("looking up '%v' and got %T %v", astReference.Name(), definition.referencedType, definition)
	return NewResolvedLocalTypeReference(definition.debugLocalTypeName, definition), nil
}

func (t *ResolvedLocalTypeContext) LookupTypeRef(decReference *LocalTypeNameReference) (
	*ResolvedLocalTypeReference, decshared.DecoratedError,
) {
	return t.LookupTypeAstRef(decReference.identifier)
}

func (t *ResolvedLocalTypeContext) LookupTypeName(astReference *ast.LocalTypeName) (dtype.Type, error) {
	if astReference == nil {
		panic(fmt.Errorf("ast localtypename is nil"))
	}
	definition, found := t.resolvedArguments[astReference.Name()]
	if !found {
		return nil, fmt.Errorf("could not find %v", astReference.Name())
	}
	return definition, nil
}

func (t *ResolvedLocalTypeContext) Verify() error {
	for _, x := range t.resolvedArguments {
		if !x.WasReferenced() {
			return fmt.Errorf(
				"ResolvedLocalTypeContext:Verify. Argument name %v has not been resolved", t.resolvedArguments,
			)
		}
	}

	return nil
}
