/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type DefinitionAssignment struct {
	identifier *VariableIdentifier
	expression Expression
}

func NewDefinitionAssignment(identifier *VariableIdentifier, expression Expression) *DefinitionAssignment {
	return &DefinitionAssignment{identifier: identifier, expression: expression}
}

func (i *DefinitionAssignment) Identifier() *VariableIdentifier {
	return i.identifier
}

func (i *DefinitionAssignment) Expression() Expression {
	return i.expression
}

func (i *DefinitionAssignment) PositionLength() token.PositionLength {
	return i.identifier.PositionLength()
}

func (i *DefinitionAssignment) String() string {
	return fmt.Sprintf("[definition: %v = %v]", i.identifier, i.expression)
}

func (i *DefinitionAssignment) DebugString() string {
	return fmt.Sprintf("[assignmentdef]")
}
