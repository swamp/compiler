/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"
	"strings"
)

type TypeParameterIdentifierContext struct {
	lookup                   map[string]*TypeParameter
	typeParameterIdentifiers []*TypeParameter
}

func (t *TypeParameterIdentifierContext) AllTypeParameters() []*TypeParameter {
	return t.typeParameterIdentifiers
}

func (t *TypeParameterIdentifierContext) arrayToString() string {
	var typeParams []string
	if len(t.typeParameterIdentifiers) == 0 {
		return ""
	}
	for _, v := range t.typeParameterIdentifiers {
		typeParams = append(typeParams, v.Identifier().String())
	}
	return "[" + strings.Join(typeParams, " ") + "]"
}

func (t *TypeParameterIdentifierContext) String() string {
	return fmt.Sprintf("[type-param-context %s]", t.arrayToString())
}

func (t *TypeParameterIdentifierContext) HasTypeParameter(parameter *TypeParameter) bool {
	return t.lookup[parameter.Identifier().Name()] != nil
}

func NewTypeParameterIdentifierContext(typeParameterIdentifiers []*TypeParameter) *TypeParameterIdentifierContext {
	lookup := make(map[string]*TypeParameter)
	for _, typeParameterIdentifier := range typeParameterIdentifiers {
		lookup[typeParameterIdentifier.Identifier().Name()] = typeParameterIdentifier
	}
	return &TypeParameterIdentifierContext{lookup: lookup, typeParameterIdentifiers: typeParameterIdentifiers}
}
