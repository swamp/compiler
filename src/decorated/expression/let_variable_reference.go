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

type LetVariableReference struct {
	ident               *ast.VariableIdentifier
	decoratedExpression DecoratedExpression
}

func (g *LetVariableReference) Type() dtype.Type {
	return g.decoratedExpression.Type()
}

func (g *LetVariableReference) String() string {
	return fmt.Sprintf("[letvarref %v %v]", g.ident, g.decoratedExpression)
}

func (g *LetVariableReference) Identifier() *ast.VariableIdentifier {
	return g.ident
}

func (g *LetVariableReference) Expression() DecoratedExpression {
	return g.decoratedExpression
}

func NewLetVariableReference(ident *ast.VariableIdentifier,
	decoratedExpression DecoratedExpression) *LetVariableReference {
	if decoratedExpression == nil {
		panic("cant be nil")
	}

	return &LetVariableReference{ident: ident, decoratedExpression: decoratedExpression}
}

func (g *LetVariableReference) FetchPositionLength() token.SourceFileReference {
	return g.ident.FetchPositionLength()
}
