/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dtype

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
)

type TypeArgumentName struct {
	name *ast.VariableIdentifier
}

func NewTypeArgumentName(name *ast.VariableIdentifier) *TypeArgumentName {
	return &TypeArgumentName{name: name}
}

func (t *TypeArgumentName) String() string {
	return fmt.Sprintf("%v", t.name.Name())
}

func (t *TypeArgumentName) Name() string {
	return t.name.Name()
}

func (t *TypeArgumentName) VariableIdentifier() *ast.VariableIdentifier {
	return t.name
}
