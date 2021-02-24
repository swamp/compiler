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

type FunctionParameterReference struct {
	ident               *ast.VariableIdentifier
	referencedParameter *FunctionParameterDefinition
}

func (g *FunctionParameterReference) Type() dtype.Type {
	return g.referencedParameter.Type()
}

func (g *FunctionParameterReference) String() string {
	return fmt.Sprintf("[functionparamref %v %v]", g.ident, g.referencedParameter)
}

func (g *FunctionParameterReference) Identifier() *ast.VariableIdentifier {
	return g.ident
}

func (g *FunctionParameterReference) ParameterRef() *FunctionParameterDefinition {
	return g.referencedParameter
}

func NewFunctionParameterReference(ident *ast.VariableIdentifier,
	referencedParameter *FunctionParameterDefinition) *FunctionParameterReference {
	if referencedParameter == nil {
		panic("cant be nil")
	}

	return &FunctionParameterReference{ident: ident, referencedParameter: referencedParameter}
}

func (g *FunctionParameterReference) FetchPositionLength() token.SourceFileReference {
	return g.ident.FetchPositionLength()
}
