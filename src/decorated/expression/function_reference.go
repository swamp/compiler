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

type FunctionReference struct {
	definitionReference     *NamedDefinitionReference
	referencedFunctionValue *FunctionValue
}

func (g *FunctionReference) Type() dtype.Type {
	return g.referencedFunctionValue.Type()
}

func (g *FunctionReference) String() string {
	return fmt.Sprintf("[FunctionRef %v]", g.definitionReference)
}

func (g *FunctionReference) DebugString() string {
	return fmt.Sprintf("[FunctionRef %v]", g.definitionReference)
}

func (g *FunctionReference) HumanReadable() string {
	return "FunctionRef"
}

func (g *FunctionReference) Identifier() ast.ScopedOrNormalVariableIdentifier {
	return g.definitionReference.ident
}

func (g *FunctionReference) NameReference() *NamedDefinitionReference {
	return g.definitionReference
}

func (g *FunctionReference) FunctionValue() *FunctionValue {
	return g.referencedFunctionValue
}

func NewFunctionReference(definitionReference *NamedDefinitionReference,
	referencedFunctionValue *FunctionValue) *FunctionReference {
	if referencedFunctionValue == nil {
		panic("cant be nil")
	}

	ref := &FunctionReference{definitionReference: definitionReference, referencedFunctionValue: referencedFunctionValue}

	referencedFunctionValue.AddReferee(ref)

	return ref
}

func (g *FunctionReference) FetchPositionLength() token.SourceFileReference {
	return g.definitionReference.FetchPositionLength()
}
