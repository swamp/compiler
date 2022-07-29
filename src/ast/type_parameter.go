/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type TypeParameter struct {
	ident *VariableIdentifier
}

func (t *TypeParameter) Identifier() *VariableIdentifier {
	return t.ident
}

func NewTypeParameter(ident *VariableIdentifier) *TypeParameter {
	return &TypeParameter{ident: ident}
}

func (t *TypeParameter) String() string {
	return fmt.Sprintf("[TypeParam %v]", t.ident)
}

func (t *TypeParameter) Name() string {
	return t.ident.Name()
}

func (t *TypeParameter) FetchPositionLength() token.SourceFileReference {
	return t.ident.FetchPositionLength()
}
