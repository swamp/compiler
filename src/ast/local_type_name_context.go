/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"
	"github.com/swamp/compiler/src/token"
	"strings"
)

type LocalTypeNameDefinitionContextDynamic interface {
	GetOrCreateReferenceFromName(parameter *LocalTypeName) (*LocalTypeNameReference, error)
}

type LocalTypeNameDefinitionContext struct {
	lookup                   map[string]*LocalTypeNameDefinition
	typeParameterIdentifiers []*LocalTypeNameDefinition
	nextType                 Type
	isDynamic                bool
}

func (t *LocalTypeNameDefinitionContext) arrayToString() string {
	var typeParams []string
	if len(t.typeParameterIdentifiers) == 0 {
		return ""
	}
	for _, v := range t.typeParameterIdentifiers {
		typeParams = append(typeParams, v.Identifier().Name())
	}
	return "[" + strings.Join(typeParams, " ") + "]"
}

func (t *LocalTypeNameDefinitionContext) String() string {
	return fmt.Sprintf("[type-param-context %s]", t.arrayToString())
}

func (t *LocalTypeNameDefinitionContext) DebugString() string {
	return fmt.Sprintf("[type-param-context %s]", t.arrayToString())
}

func (t *LocalTypeNameDefinitionContext) GetOrCreateReferenceFromName(parameter *LocalTypeName) (*LocalTypeNameReference, error) {
	if !t.isDynamic {
		return t.ParseReferenceFromName(parameter)
	}
	def, hasLookup := t.lookup[parameter.Identifier().Name()]
	if !hasLookup {
		def = NewLocalTypeNameDefinition(parameter)
		t.typeParameterIdentifiers = append(t.typeParameterIdentifiers, def)
		t.lookup[parameter.Name()] = def
	}

	return NewLocalTypeNameReference(def), nil
}

func (t *LocalTypeNameDefinitionContext) HasTypeParameter(parameter *LocalTypeNameDefinition) bool {
	return t.lookup[parameter.Identifier().Name()] != nil
}

func (t *LocalTypeNameDefinitionContext) IsEmpty() bool {
	return len(t.typeParameterIdentifiers) == 0
}

func (t *LocalTypeNameDefinitionContext) ParseReferenceFromName(parameter *LocalTypeName) (*LocalTypeNameReference, error) {
	definition, foundDefinition := t.lookup[parameter.Identifier().Name()]
	if !foundDefinition {
		return nil, NewUnknownTypeParameterError(parameter, t)
	}

	return NewLocalTypeNameReference(definition), nil
}

func (t *LocalTypeNameDefinitionContext) LocalTypeNames() []*LocalTypeName {
	var names []*LocalTypeName
	for _, x := range t.typeParameterIdentifiers {
		names = append(names, x.ident)
	}

	return names
}

func (t *LocalTypeNameDefinitionContext) FetchPositionLength() token.SourceFileReference {
	return t.nextType.FetchPositionLength()
}

func (t *LocalTypeNameDefinitionContext) Name() string {
	return "LocalTypeNameDefinitionContext"
}

func (t *LocalTypeNameDefinitionContext) SetNextType(p Type) {
	t.nextType = p
}

func NewLocalTypeNameContext(typeParameterNames []*LocalTypeName, nextType Type) *LocalTypeNameDefinitionContext {
	lookup := make(map[string]*LocalTypeNameDefinition)
	var localTypeDefs []*LocalTypeNameDefinition

	for _, typeParameterIdentifier := range typeParameterNames {
		newLocalDef := NewLocalTypeNameDefinition(typeParameterIdentifier)
		localTypeDefs = append(localTypeDefs, newLocalDef)
		lookup[typeParameterIdentifier.Name()] = newLocalDef
	}
	return &LocalTypeNameDefinitionContext{lookup: lookup, nextType: nextType, typeParameterIdentifiers: localTypeDefs, isDynamic: len(localTypeDefs) == 0}
}
