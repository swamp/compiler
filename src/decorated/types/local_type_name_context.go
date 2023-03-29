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
	"strings"
)

type LocalTypeNameContext struct {
	resolvedArguments             map[string]*LocalTypeNameDefinition
	definitions                   []*LocalTypeNameDefinition
	typeThatIsReferencingTheNames dtype.Type
}

func (t *LocalTypeNameContext) DeclareString() string {
	if len(t.resolvedArguments) == 0 {
		return ""
	}
	return fmt.Sprintf("%v", t.resolvedArguments)
}

func (t *LocalTypeNameContext) Names() []*dtype.LocalTypeName {
	var x []*dtype.LocalTypeName

	for _, f := range t.definitions {
		x = append(x, f.identifier)
	}

	return x
}

func (t *LocalTypeNameContext) Definitions() []*LocalTypeNameDefinition {
	return t.definitions
}

func (t *LocalTypeNameContext) NamesString() string {
	var x []string

	for _, f := range t.definitions {
		x = append(x, f.identifier.Name())
	}

	return strings.Join(x, ", ")
}

func (t *LocalTypeNameContext) String() string {
	return fmt.Sprintf("[LocalTypeNameContext %v = %v]", t.NamesString(), t.typeThatIsReferencingTheNames)
}

func (t *LocalTypeNameContext) Next() dtype.Type {
	return t.typeThatIsReferencingTheNames
}

func (t *LocalTypeNameContext) DebugString() string {
	s := ""
	for name, argumentType := range t.resolvedArguments {
		s += fmt.Sprintf("%v = %v", name, argumentType)
	}
	return s
}

func NewLocalTypeNameContext() *LocalTypeNameContext {
	t := &LocalTypeNameContext{resolvedArguments: make(map[string]*LocalTypeNameDefinition),
		typeThatIsReferencingTheNames: nil}

	/*
		for _, name := range names {
			newLocalTypeDef := NewLocalTypeDefinition(name, NewAnyType())
			t.resolvedArguments[name.Name()] = newLocalTypeDef
			t.definitions = append(t.definitions, newLocalTypeDef)
		}

	*/
	return t
}

func (t *LocalTypeNameContext) AtomName() string {
	return "NameOnlyContext"
}

func (t *LocalTypeNameContext) IsEqual(other dtype.Atom) error {
	return nil
}

func (t *LocalTypeNameContext) HumanReadable() string {
	return fmt.Sprintf("local type name context")
}

func (t *LocalTypeNameContext) Resolve() (dtype.Atom, error) {
	return t, nil //fmt.Errorf("can not be resolved since it is a type name context %T", t)
}

func (t *LocalTypeNameContext) ParameterCount() int {
	return 0
}

func (t *LocalTypeNameContext) FetchPositionLength() token.SourceFileReference {
	return t.Next().FetchPositionLength()
}

func (t *LocalTypeNameContext) WasReferenced() bool {
	return true
}

func (t *LocalTypeNameContext) SetType(d dtype.Type) {
	t.typeThatIsReferencingTheNames = d
}

func (t *LocalTypeNameContext) HasDefinitions() bool {
	return len(t.definitions) > 0
}

func (t *LocalTypeNameContext) AddDef(identifier *dtype.LocalTypeName) *LocalTypeNameDefinition {
	nameDef := NewLocalTypeNameDefinition(identifier)
	t.resolvedArguments[identifier.Name()] = nameDef
	t.definitions = append(t.definitions, nameDef)

	return nameDef
}

func (t *LocalTypeNameContext) ReferenceNameOnly(identifier *ast.LocalTypeNameReference) (*LocalTypeNameReference, error) {
	found := t.resolvedArguments[identifier.Name()]
	if found == nil {
		return nil, fmt.Errorf("could not find %v", identifier)
	}

	return NewLocalTypeNameReference(identifier, found), nil
}
