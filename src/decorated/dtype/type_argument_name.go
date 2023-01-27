/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dtype

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
)

type LocalTypeName struct {
	name *ast.LocalTypeName
}

func NewLocalTypeName(name *ast.LocalTypeName) *LocalTypeName {
	return &LocalTypeName{name: name}
}

func (t *LocalTypeName) String() string {
	return fmt.Sprintf("%v", t.name.Name())
}

func (t *LocalTypeName) Name() string {
	return t.name.Name()
}

func (t *LocalTypeName) LocalType() *ast.LocalTypeName {
	return t.name
}
