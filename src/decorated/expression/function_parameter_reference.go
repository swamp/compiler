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
	ident               ast.ScopedOrNormalVariableIdentifier
	referencedParameter *FunctionParameterDefinition `debug:"true"`
}

func (g *FunctionParameterReference) Type() dtype.Type {
	return g.referencedParameter.Type()
}

func (g *FunctionParameterReference) String() string {
	return fmt.Sprintf("[ParamRef %v]", g.referencedParameter)
}

func (g *FunctionParameterReference) HumanReadable() string {
	return "FunctionParamRef"
}

func (g *FunctionParameterReference) Identifier() ast.ScopedOrNormalVariableIdentifier {
	return g.ident
}

func (g *FunctionParameterReference) ParameterRef() *FunctionParameterDefinition {
	return g.referencedParameter
}

func NewFunctionParameterReference(
	ident ast.ScopedOrNormalVariableIdentifier,
	referencedParameter *FunctionParameterDefinition,
) *FunctionParameterReference {
	if referencedParameter == nil {
		panic("cant be nil")
	}

	ref := &FunctionParameterReference{ident: ident, referencedParameter: referencedParameter}

	referencedParameter.AddReferee(ref)

	return ref
}

func (g *FunctionParameterReference) FetchPositionLength() token.SourceFileReference {
	return g.ident.FetchPositionLength()
}
