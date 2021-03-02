/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type StringInterpolation struct {
	str        *ast.StringInterpolation
	expression Expression
}

func NewStringInterpolation(str *ast.StringInterpolation, expression Expression) *StringInterpolation {
	return &StringInterpolation{str: str, expression: expression}
}

func (i *StringInterpolation) Type() dtype.Type {
	return i.expression.Type()
}

func (i *StringInterpolation) Expression() Expression {
	return i.expression
}

func (i *StringInterpolation) AstStringInterpolation() *ast.StringInterpolation {
	return i.str
}

func (i *StringInterpolation) String() string {
	return fmt.Sprintf("[strinterpolation %v]", i.str)
}

func (i *StringInterpolation) FetchPositionLength() token.SourceFileReference {
	return i.str.FetchPositionLength()
}
