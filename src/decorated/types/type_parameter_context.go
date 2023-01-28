/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	"log"
)

type TypeParameterContext struct {
	resolvedArguments map[string]*LocalTypeDefinition
	definitions       []*LocalTypeDefinition
}

func (t *TypeParameterContext) DeclareString() string {
	if len(t.resolvedArguments) == 0 {
		return ""
	}
	return fmt.Sprintf("%v", t.resolvedArguments)
}

func (t *TypeParameterContext) String() string {
	return fmt.Sprintf("[typeparamcontext %v]", t.resolvedArguments)
}

func (t *TypeParameterContext) DebugString() string {
	s := ""
	for name, argumentType := range t.resolvedArguments {
		s += fmt.Sprintf("%v = %v", name, argumentType)
	}
	return s
}

func NewTypeParameterContext() *TypeParameterContext {
	t := &TypeParameterContext{resolvedArguments: make(map[string]*LocalTypeDefinition)}

	/*
		for _, name := range names {
			newLocalTypeDef := NewLocalTypeNameDefinition(name, NewAnyType())
			t.resolvedArguments[name.Name()] = newLocalTypeDef
			t.definitions = append(t.definitions, newLocalTypeDef)
		}

	*/
	return t
}

func (t *TypeParameterContext) SetTypes(types []dtype.Type) error {
	if len(types) != len(t.definitions) {
		return fmt.Errorf("wrong number of definitions")
	}

	for index, typeToSet := range types {
		def := t.definitions[index]
		if err := def.SetDefinition(typeToSet); err != nil {
			return err
		}
	}

	return nil
}

func (t *TypeParameterContext) IsDefined() bool {
	for _, def := range t.definitions {
		if !def.hasBeenDefined {
			return false
		}
	}

	return true
}

func (t *TypeParameterContext) AddExpectedDef(name *dtype.LocalTypeName) {
	def := NewLocalTypeDefinition(name, NewAnyType())
	t.resolvedArguments[name.Name()] = def
	t.definitions = append(t.definitions, def)
}

func (t *TypeParameterContext) AddExpectedDefs(names []*dtype.LocalTypeName) {
	for _, name := range names {
		t.AddExpectedDef(name)
	}
}

func (t *TypeParameterContext) SetType(defRef *LocalTypeNameReference, definedType dtype.Type) (*LocalTypeDefinitionReference, error) {
	definition, found := t.resolvedArguments[defRef.Identifier().Name()]
	if !found {
		return nil, fmt.Errorf("could not find %v %v", defRef.Identifier().Name(), t)
	}

	localNameRef, wasLocal := definedType.(*LocalTypeNameReference)
	if wasLocal {
		panic(fmt.Errorf("not allowed to set a type to a name reference, that won't help us %v", localNameRef))
	}
	if definition.hasBeenDefined {
		if err := CompatibleTypes(definition.referencedType, definedType); err != nil {
			return nil, fmt.Errorf(" %v was already set %w", defRef.Identifier().Name(), err)
		}
	} else {
		log.Printf("set %v to %v", defRef.Identifier().Name(), definedType)
		definition.referencedType = definedType
		definition.hasBeenDefined = true
	}

	return NewLocalTypeDefinitionReference(defRef.Identifier(), definition), nil
}

func (t *TypeParameterContext) HasDefinitions() bool {
	return len(t.definitions) > 0
}

func (t *TypeParameterContext) LookupType(definition *LocalTypeDefinitionReference) *LocalTypeDefinition {
	return t.resolvedArguments[definition.identifier.Name()]
}

/*
func (t *TypeParameterContext) ResolveTypeRef(defRef *LocalTypeDefinitionReference) (*LocalTypeDefinitionReference, error) {
	definition, found := t.resolvedArguments[defRef.Identifier().Name()]
	if !found {
		return nil, fmt.Errorf("could not find %v", defRef.Identifier().Name())
	}

	return NewLocalTypeDefinitionReference(definition, definition), nil
}
*/

func (t *TypeParameterContext) LookupTypeAstRef(astReference *ast.LocalTypeNameReference) (*LocalTypeDefinitionReference, decshared.DecoratedError) {
	definition, found := t.resolvedArguments[astReference.Name()]
	if !found {
		return nil, NewCouldNotFindLocalTypeName(astReference, fmt.Errorf("could not find %v", astReference.Name()))
	}
	return NewLocalTypeDefinitionReference(astReference, definition), nil
}

func (t *TypeParameterContext) LookupTypeRef(decReference *LocalTypeNameReference) (*LocalTypeDefinitionReference, decshared.DecoratedError) {
	return t.LookupTypeAstRef(decReference.identifier)
}

func (t *TypeParameterContext) LookupTypeName(astReference *ast.LocalTypeName) (dtype.Type, error) {
	definition, found := t.resolvedArguments[astReference.Name()]
	if !found {
		return nil, fmt.Errorf("could not find %v", astReference.Name())
	}
	return definition, nil
}

func (t *TypeParameterContext) Verify() error {
	for _, x := range t.resolvedArguments {
		if !x.WasReferenced() {
			return fmt.Errorf("TypeParameterContext:Verify. Argument name %v has not been resolved", t.resolvedArguments)
		}
	}

	return nil
}
