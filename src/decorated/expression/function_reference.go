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
	ident                   *ast.VariableIdentifier
	referencedFunctionValue *FunctionValue
}

func (g *FunctionReference) Type() dtype.Type {
	return g.referencedFunctionValue.Type()
}

func (g *FunctionReference) String() string {
	return fmt.Sprintf("[functionref %v %v]", g.ident, g.referencedFunctionValue)
}

func (g *FunctionReference) DebugString() string {
	return fmt.Sprintf("[functionref %v %v]", g.ident, g.referencedFunctionValue)
}

func (g *FunctionReference) HumanReadable() string {
	return "function reference"
}

func (g *FunctionReference) Identifier() *ast.VariableIdentifier {
	return g.ident
}

func (g *FunctionReference) FunctionValue() *FunctionValue {
	return g.referencedFunctionValue
}

func NewFunctionReference(ident *ast.VariableIdentifier,
	referencedFunctionValue *FunctionValue) *FunctionReference {
	if referencedFunctionValue == nil {
		panic("cant be nil")
	}

	ref := &FunctionReference{ident: ident, referencedFunctionValue: referencedFunctionValue}

	referencedFunctionValue.AddReferee(ref)

	return ref
}

func (g *FunctionReference) FetchPositionLength() token.SourceFileReference {
	return g.ident.FetchPositionLength()
}
