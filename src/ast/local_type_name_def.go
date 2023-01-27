/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type LocalTypeNameDefinition struct {
	ident      *LocalTypeName
	references []*LocalTypeNameReference
}

func (t *LocalTypeNameDefinition) Identifier() *LocalTypeName {
	return t.ident
}

func NewLocalTypeNameDefinition(name *LocalTypeName) *LocalTypeNameDefinition {
	return &LocalTypeNameDefinition{ident: name}
}

func (t *LocalTypeNameDefinition) String() string {
	return fmt.Sprintf("[LocalTypeNameDef %v]", t.ident)
}

func (t *LocalTypeNameDefinition) Name() string {
	return t.ident.Name()
}

func (t *LocalTypeNameDefinition) FetchPositionLength() token.SourceFileReference {
	return t.ident.FetchPositionLength()
}

func (t *LocalTypeNameDefinition) AddReference(reference *LocalTypeNameReference) {
	t.references = append(t.references, reference)
}
