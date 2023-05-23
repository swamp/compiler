/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"bytes"
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type CaseConsequenceParameterForCustomType struct {
	name          *ast.VariableIdentifier
	parameterType dtype.Type `debug:"true"`
	references    []*CaseConsequenceParameterReference
}

func (c *CaseConsequenceParameterForCustomType) String() string {
	return fmt.Sprintf("[dcaseparm %v type:%v]", c.name, c.parameterType)
}

func (c *CaseConsequenceParameterForCustomType) Identifier() *ast.VariableIdentifier {
	return c.name
}

func (c *CaseConsequenceParameterForCustomType) Type() dtype.Type {
	return c.parameterType
}

func (c *CaseConsequenceParameterForCustomType) FetchPositionLength() token.SourceFileReference {
	return c.name.FetchPositionLength()
}

func (c *CaseConsequenceParameterForCustomType) HumanReadable() string {
	return "custom type variant parameter"
}

func (c *CaseConsequenceParameterForCustomType) References() []*CaseConsequenceParameterReference {
	return c.references
}

func (c *CaseConsequenceParameterForCustomType) AddReferee(ref *CaseConsequenceParameterReference) {
	c.references = append(c.references, ref)
}

func NewCaseConsequenceParameterForCustomType(name *ast.VariableIdentifier,
	parameterType dtype.Type) *CaseConsequenceParameterForCustomType {
	return &CaseConsequenceParameterForCustomType{name: name, parameterType: parameterType}
}

type CaseConsequenceForCustomType struct {
	variantName    *dectype.CustomTypeVariantReference
	parameters     []*CaseConsequenceParameterForCustomType `debug:"true"`
	expression     Expression
	internalIndex  int
	astConsequence *ast.CaseConsequenceForCustomType
}

func NewCaseConsequenceForCustomType(internalIndex int, variantName *dectype.CustomTypeVariantReference,
	parameters []*CaseConsequenceParameterForCustomType,
	expression Expression, astConsequence *ast.CaseConsequenceForCustomType) *CaseConsequenceForCustomType {
	return &CaseConsequenceForCustomType{
		internalIndex: internalIndex, variantName: variantName, parameters: parameters,
		expression: expression, astConsequence: astConsequence,
	}
}

func (c *CaseConsequenceForCustomType) AstConsequence() *ast.CaseConsequenceForCustomType {
	return c.astConsequence
}

func (c *CaseConsequenceForCustomType) Expression() Expression {
	return c.expression
}

func (c *CaseConsequenceForCustomType) InternalIndex() int {
	return c.internalIndex
}

func (c *CaseConsequenceForCustomType) VariantReference() *dectype.CustomTypeVariantReference {
	return c.variantName
}

func (c *CaseConsequenceForCustomType) Parameters() []*CaseConsequenceParameterForCustomType {
	return c.parameters
}

func (c *CaseConsequenceForCustomType) String() string {
	return fmt.Sprintf("[dcasecons %v (%v) => %v]", c.variantName, c.parameters, c.expression)
}

func caseConsequenceArrayToStringEx(expressions []*CaseConsequenceForCustomType, ch string) string {
	var out bytes.Buffer

	for index, expression := range expressions {
		if index > 0 {
			out.WriteString(ch)
		}
		out.WriteString(expression.String())
	}
	return out.String()
}

type CaseCustomType struct {
	test           Expression                      `debug:"true"`
	cases          []*CaseConsequenceForCustomType `debug:"true"`
	defaultCase    Expression
	caseExpression *ast.CaseForCustomType
}

func NewCaseCustomType(caseExpression *ast.CaseForCustomType, test Expression, cases []*CaseConsequenceForCustomType,
	defaultCase Expression) (*CaseCustomType, decshared.DecoratedError) {
	return &CaseCustomType{caseExpression: caseExpression, test: test, cases: cases, defaultCase: defaultCase}, nil
}

func (i *CaseCustomType) AstCaseCustomType() *ast.CaseForCustomType {
	return i.caseExpression
}

func (i *CaseCustomType) Type() dtype.Type {
	if len(i.cases) == 0 {
		return i.defaultCase.Type()
	}
	firstCase := i.cases[0]
	return firstCase.Expression().Type()
}

func (i *CaseCustomType) String() string {
	if i.defaultCase != nil {
		return fmt.Sprintf("[dcase: %v of %v default: %v]", i.test,
			caseConsequenceArrayToStringEx(i.cases, ";"), i.defaultCase)
	}
	return fmt.Sprintf("[dcase: %v of %v]", i.test, caseConsequenceArrayToStringEx(i.cases, ";"))
}

func (i *CaseCustomType) Test() Expression {
	return i.test
}

func (i *CaseCustomType) Consequences() []*CaseConsequenceForCustomType {
	return i.cases
}

func (i *CaseCustomType) DefaultCase() Expression {
	return i.defaultCase
}

func (i *CaseCustomType) DebugString() string {
	return "[dcase]"
}

func (i *CaseCustomType) FetchPositionLength() token.SourceFileReference {
	return i.caseExpression.FetchPositionLength()
}
