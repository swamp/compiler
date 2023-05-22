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

type Constant struct {
	astConstant       *ast.ConstantDefinition
	identifier        ast.ScopedOrNormalVariableIdentifier
	expression        Expression `debug:"true"`
	references        []*ConstantReference
	localCommentBlock *ast.MultilineComment
	inclusive         token.SourceFileReference
}

func NewConstant(identifier ast.ScopedOrNormalVariableIdentifier, astConstant *ast.ConstantDefinition, expression Expression, localCommentBlock *ast.MultilineComment) *Constant {
	inclusive := token.MakeInclusiveSourceFileReference(identifier.FetchPositionLength(), expression.FetchPositionLength())
	return &Constant{astConstant: astConstant, identifier: identifier, expression: expression, localCommentBlock: localCommentBlock, inclusive: inclusive}
}

func (c *Constant) String() string {
	return fmt.Sprintf("[Constant %v]", c.expression)
}

func (c *Constant) CommentBlock() *ast.MultilineComment {
	return c.localCommentBlock
}

func (c *Constant) AstConstant() *ast.ConstantDefinition {
	return c.astConstant
}

func (a *Constant) AddReferee(ref *ConstantReference) {
	a.references = append(a.references, ref)
}

func (a *Constant) References() []*ConstantReference {
	return a.references
}

func (c *Constant) Expression() Expression {
	return c.expression
}

func (c *Constant) FetchPositionLength() token.SourceFileReference {
	return c.inclusive
}

func (c *Constant) HumanReadable() string {
	return "Constant"
}

func (n *Constant) StatementString() string {
	return fmt.Sprintf("constant customTypeAtom %v = %v", n.identifier, n.expression)
}

func (c *Constant) Type() dtype.Type {
	return c.expression.Type()
}
