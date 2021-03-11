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

type CaseConsequenceParameterReference struct {
	ident               ast.ScopedOrNormalVariableIdentifier
	referencedParameter *CaseConsequenceParameterForCustomType
}

func (g *CaseConsequenceParameterReference) Type() dtype.Type {
	return g.referencedParameter.Type()
}

func (g *CaseConsequenceParameterReference) String() string {
	return fmt.Sprintf("[functionparamref %v %v]", g.ident, g.referencedParameter)
}

func (g *CaseConsequenceParameterReference) HumanReadable() string {
	return fmt.Sprintf("custom type variant Parameter Reference")
}

func (g *CaseConsequenceParameterReference) Identifier() ast.ScopedOrNormalVariableIdentifier {
	return g.ident
}

func (g *CaseConsequenceParameterReference) ParameterRef() *CaseConsequenceParameterForCustomType {
	return g.referencedParameter
}

func NewCaseConsequenceParameterReference(ident ast.ScopedOrNormalVariableIdentifier,
	referencedParameter *CaseConsequenceParameterForCustomType) *CaseConsequenceParameterReference {
	if referencedParameter == nil {
		panic("cant be nil")
	}

	ref := &CaseConsequenceParameterReference{ident: ident, referencedParameter: referencedParameter}

	referencedParameter.AddReferee(ref)

	return ref
}

func (g *CaseConsequenceParameterReference) FetchPositionLength() token.SourceFileReference {
	return g.ident.FetchPositionLength()
}
