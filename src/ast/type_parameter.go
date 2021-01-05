/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import "fmt"

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
	return fmt.Sprintf("[type-param %v]", t.ident)
}

func (t *TypeParameter) Name() string {
	return t.ident.Name()
}
