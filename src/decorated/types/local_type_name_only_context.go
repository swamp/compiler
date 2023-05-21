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

type LocalTypeNameOnlyContext struct {
	resolvedArguments             map[string]*LocalTypeName
	definitions                   []*LocalTypeName
	typeThatIsReferencingTheNames dtype.Type
}

func (t *LocalTypeNameOnlyContext) DeclareString() string {
	if len(t.resolvedArguments) == 0 {
		return ""
	}
	return fmt.Sprintf("%v", t.resolvedArguments)
}

func (t *LocalTypeNameOnlyContext) Names() []*dtype.LocalTypeName {
	var x []*dtype.LocalTypeName

	for _, f := range t.definitions {
		x = append(x, f.identifier)
	}

	return x
}

func (t *LocalTypeNameOnlyContext) Definitions() []*LocalTypeName {
	return t.definitions
}

func (t *LocalTypeNameOnlyContext) NamesString() string {
	var x []string

	for _, f := range t.definitions {
		x = append(x, f.identifier.Name())
	}

	return strings.Join(x, ", ")
}

func (t *LocalTypeNameOnlyContext) String() string {
	return fmt.Sprintf("[LocalTypeNameOnlyContext %v => %v]", t.NamesString(), t.typeThatIsReferencingTheNames)
}

func (t *LocalTypeNameOnlyContext) Next() dtype.Type {
	return t.typeThatIsReferencingTheNames
}

func (t *LocalTypeNameOnlyContext) DebugString() string {
	s := ""
	for name, argumentType := range t.resolvedArguments {
		s += fmt.Sprintf("%v = %v", name, argumentType)
	}
	return s
}

func NewLocalTypeNameContext() *LocalTypeNameOnlyContext {
	t := &LocalTypeNameOnlyContext{resolvedArguments: make(map[string]*LocalTypeName),
		typeThatIsReferencingTheNames: nil}

	/*
		for _, name := range names {
			newLocalTypeDef := NewResolvedLocalType(name, NewAnyType())
			t.resolvedArguments[name.Name()] = newLocalTypeDef
			t.definitions = append(t.definitions, newLocalTypeDef)
		}

	*/
	return t
}

func (t *LocalTypeNameOnlyContext) AtomName() string {
	return "NameOnlyContext"
}

func (t *LocalTypeNameOnlyContext) IsEqual(other dtype.Atom) error {
	return nil
}

func (t *LocalTypeNameOnlyContext) HumanReadable() string {
	return fmt.Sprintf("local type name context")
}

func (t *LocalTypeNameOnlyContext) Resolve() (dtype.Atom, error) {
	return t, nil //fmt.Errorf("can not be resolved since it is a type name context %T", t)
}

func (t *LocalTypeNameOnlyContext) ParameterCount() int {
	return 0
}

func (t *LocalTypeNameOnlyContext) FetchPositionLength() token.SourceFileReference {
	return t.Next().FetchPositionLength()
}

func (t *LocalTypeNameOnlyContext) WasReferenced() bool {
	return true
}

func (t *LocalTypeNameOnlyContext) SetType(d dtype.Type) {
	t.typeThatIsReferencingTheNames = d
}

func (t *LocalTypeNameOnlyContext) HasDefinitions() bool {
	return len(t.definitions) > 0
}

func (t *LocalTypeNameOnlyContext) AddDef(identifier *dtype.LocalTypeName) *LocalTypeName {
	nameDef := NewLocalTypeName(identifier)
	t.resolvedArguments[identifier.Name()] = nameDef
	t.definitions = append(t.definitions, nameDef)

	return nameDef
}

func (t *LocalTypeNameOnlyContext) ReferenceNameOnly(identifier *ast.LocalTypeNameReference) (*LocalTypeNameReference, error) {
	found := t.resolvedArguments[identifier.Name()]
	if found == nil {
		return nil, fmt.Errorf("could not find %v", identifier)
	}

	return NewLocalTypeNameReference(identifier, found), nil
}
