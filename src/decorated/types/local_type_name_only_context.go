/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"
	"strings"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type LocalTypeNameOnlyContext struct {
	localTypeNamesMap             map[string]*LocalTypeName
	localTypeNames                []*LocalTypeName
	typeThatIsReferencingTheNames dtype.Type `debug:"true"`
}

func (t *LocalTypeNameOnlyContext) DeclareString() string {
	if len(t.localTypeNamesMap) == 0 {
		return ""
	}
	return fmt.Sprintf("%v", t.localTypeNamesMap)
}

func (t *LocalTypeNameOnlyContext) Names() []*dtype.LocalTypeName {
	var x []*dtype.LocalTypeName

	for _, f := range t.localTypeNames {
		x = append(x, f.LocalTypeName())
	}

	return x
}

func (t *LocalTypeNameOnlyContext) Definitions() []*LocalTypeName {
	return t.localTypeNames
}

func (t *LocalTypeNameOnlyContext) NamesString() string {
	var x []string

	for _, f := range t.localTypeNames {
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
	for name, argumentType := range t.localTypeNamesMap {
		s += fmt.Sprintf("%v = %v", name, argumentType)
	}
	return s
}

func NewLocalTypeNameContext(names []*ast.LocalTypeName) *LocalTypeNameOnlyContext {
	t := &LocalTypeNameOnlyContext{
		localTypeNamesMap:             make(map[string]*LocalTypeName),
		typeThatIsReferencingTheNames: nil,
	}

	for _, name := range names {
		if !name.FetchPositionLength().Verify() {
			panic(fmt.Errorf("wrong position length"))
		}
		dLocalTypeName := dtype.NewLocalTypeName(name)
		localTypeName := NewLocalTypeName(dLocalTypeName)
		t.localTypeNamesMap[name.Name()] = localTypeName
		t.localTypeNames = append(t.localTypeNames, localTypeName)
	}

	return t
}

func (t *LocalTypeNameOnlyContext) LookupNameReference(some *ast.LocalTypeNameReference) (*LocalTypeNameReference,
	error) {
	found, wasFound := t.localTypeNamesMap[some.Name()]
	if !wasFound {
		return nil, fmt.Errorf("could not find %v", some)
	}
	ref := NewLocalTypeNameReference(some, found)
	return ref, nil
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

func (t *LocalTypeNameOnlyContext) FetchPositionLength() token.SourceFileReference {
	return t.Next().FetchPositionLength()
}

func (t *LocalTypeNameOnlyContext) WasReferenced() bool {
	return true
}

func (t *LocalTypeNameOnlyContext) SetType(d dtype.Type) {
	if !d.FetchPositionLength().Verify() {
		panic(fmt.Errorf("suspicious sub type %T", d))
	}
	t.typeThatIsReferencingTheNames = d
}

func (t *LocalTypeNameOnlyContext) HasNames() bool {
	return len(t.localTypeNames) > 0
}
