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

type ConstantReference struct {
	definitionReference *NamedDefinitionReference
	referencedConstant  *Constant
}

func (g *ConstantReference) Type() dtype.Type {
	return g.referencedConstant.Type()
}

func (g *ConstantReference) String() string {
	return fmt.Sprintf("[ConstantRef %v]", g.definitionReference)
}

func (g *ConstantReference) DebugString() string {
	return fmt.Sprintf("[ConstantRef %v]", g.definitionReference)
}

func (g *ConstantReference) HumanReadable() string {
	return "constant reference"
}

func (g *ConstantReference) Identifier() ast.ScopedOrNormalVariableIdentifier {
	return g.definitionReference.ident
}

func (g *ConstantReference) NameReference() *NamedDefinitionReference {
	return g.definitionReference
}

func (g *ConstantReference) Constant() *Constant {
	return g.referencedConstant
}

func NewConstantReference(definitionReference *NamedDefinitionReference,
	referencedConstant *Constant) *ConstantReference {
	if referencedConstant == nil {
		panic("cant be nil")
	}

	ref := &ConstantReference{definitionReference: definitionReference, referencedConstant: referencedConstant}

	referencedConstant.AddReferee(ref)

	return ref
}

func (g *ConstantReference) FetchPositionLength() token.SourceFileReference {
	return g.definitionReference.FetchPositionLength()
}
