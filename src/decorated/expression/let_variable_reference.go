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
	ident      ast.ScopedOrNormalVariableIdentifier
	assignment *LetVariable
}

func (g *LetVariableReference) Type() dtype.Type {
	return g.assignment.Type()
}

func (g *LetVariableReference) String() string {
	return fmt.Sprintf("[LetVarRef %v]", g.assignment.Name())
}

func (g *LetVariableReference) HumanReadable() string {
	return "Reference to Let variable"
}

func (g *LetVariableReference) LetVariable() *LetVariable {
	return g.assignment
}

func NewLetVariableReference(ident ast.ScopedOrNormalVariableIdentifier, assignment *LetVariable) *LetVariableReference {
	if assignment == nil {
		panic("cant be nil")
	}

	ref := &LetVariableReference{ident: ident, assignment: assignment}

	assignment.AddReferee(ref)

	return ref
}

func (g *LetVariableReference) FetchPositionLength() token.SourceFileReference {
	return g.ident.FetchPositionLength()
}
