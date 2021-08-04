/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type ConstantDefinition struct {
	identifier *VariableIdentifier
	expression Expression
	comment    *MultilineComment
}

func NewConstantDefinition(identifier *VariableIdentifier, expression Expression, comment *MultilineComment) *ConstantDefinition {
	return &ConstantDefinition{identifier: identifier, expression: expression, comment: comment}
}

func (i *ConstantDefinition) Identifier() *VariableIdentifier {
	return i.identifier
}

func (i *ConstantDefinition) Expression() Expression {
	return i.expression
}

func (i *ConstantDefinition) Comment() *MultilineComment {
	return i.comment
}

func (i *ConstantDefinition) FetchPositionLength() token.SourceFileReference {
	return i.identifier.FetchPositionLength()
}

func (i *ConstantDefinition) String() string {
	return fmt.Sprintf("[constant: %v = %v]", i.identifier, i.expression)
}

func (i *ConstantDefinition) DebugString() string {
	return "[constant]"
}
