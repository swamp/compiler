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

type GetVariableOrReferenceFunction struct {
	ident                *ast.VariableIdentifier
	referencedExpression *NamedDecoratedExpression
}

func (g *GetVariableOrReferenceFunction) Type() dtype.Type {
	return g.referencedExpression.Type()
}

func (g *GetVariableOrReferenceFunction) String() string {
	return fmt.Sprintf("[getvar %v %v]", g.ident, g.referencedExpression.Expression().Type())
}

func (g *GetVariableOrReferenceFunction) Identifier() *ast.VariableIdentifier {
	return g.ident
}

func (g *GetVariableOrReferenceFunction) NamedExpression() *NamedDecoratedExpression {
	return g.referencedExpression
}

func NewGetVariable(ident *ast.VariableIdentifier, referencedExpression *NamedDecoratedExpression) *GetVariableOrReferenceFunction {
	if referencedExpression == nil {
		panic("cant be nil")
	}
	return &GetVariableOrReferenceFunction{ident: ident, referencedExpression: referencedExpression}
}

func (g *GetVariableOrReferenceFunction) FetchPositionAndLength() token.PositionLength {
	return g.ident.PositionLength()
}
