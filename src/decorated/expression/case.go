/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"bytes"
	"fmt"

	"github.com/swamp/compiler/src/ast"
	decshared "github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type CaseConsequenceParameter struct {
	name          *ast.VariableIdentifier
	parameterType dtype.Type
}

func (c *CaseConsequenceParameter) String() string {
	return fmt.Sprintf("[dcaseparm %v type:%v]", c.name, c.parameterType)
}

func (c *CaseConsequenceParameter) Identifier() *ast.VariableIdentifier {
	return c.name
}

func (c *CaseConsequenceParameter) Type() dtype.Type {
	return c.parameterType
}

func NewCaseConsequenceParameter(name *ast.VariableIdentifier, parameterType dtype.Type) *CaseConsequenceParameter {
	return &CaseConsequenceParameter{name: name, parameterType: parameterType}
}

type CaseConsequence struct {
	variantName   *ast.TypeIdentifier
	parameters    []*CaseConsequenceParameter
	expression    DecoratedExpression
	internalIndex int
}

func NewCaseConsequence(internalIndex int, variantName *ast.TypeIdentifier, parameters []*CaseConsequenceParameter, expression DecoratedExpression) *CaseConsequence {
	return &CaseConsequence{internalIndex: internalIndex, variantName: variantName, parameters: parameters, expression: expression}
}

func (c *CaseConsequence) Expression() DecoratedExpression {
	return c.expression
}

func (c *CaseConsequence) InternalIndex() int {
	return c.internalIndex
}

func (c *CaseConsequence) Identifier() *ast.TypeIdentifier {
	return c.variantName
}

func (c *CaseConsequence) Parameters() []*CaseConsequenceParameter {
	return c.parameters
}

func (c *CaseConsequence) String() string {
	return fmt.Sprintf("[dcasecons %v (%v) => %v]", c.variantName, c.parameters, c.expression)
}

func caseConsequenceArrayToStringEx(expressions []*CaseConsequence, ch string) string {
	var out bytes.Buffer

	for index, expression := range expressions {
		if index > 0 {
			out.WriteString(ch)
		}
		out.WriteString(expression.String())
	}
	return out.String()
}

type Case struct {
	test        DecoratedExpression
	cases       []*CaseConsequence
	defaultCase DecoratedExpression
}

func NewCase(test DecoratedExpression, cases []*CaseConsequence, defaultCase DecoratedExpression) (*Case, decshared.DecoratedError) {
	return &Case{test: test, cases: cases, defaultCase: defaultCase}, nil
}

func (i *Case) Type() dtype.Type {
	if len(i.cases) == 0 {
		return i.defaultCase.Type()
	}
	firstCase := i.cases[0]
	return firstCase.Expression().Type()
}

func (i *Case) String() string {
	if i.defaultCase != nil {
		return fmt.Sprintf("[dcase: %v of %v default: %v]", i.test, caseConsequenceArrayToStringEx(i.cases, ";"), i.defaultCase)
	}
	return fmt.Sprintf("[dcase: %v of %v]", i.test, caseConsequenceArrayToStringEx(i.cases, ";"))
}

func (i *Case) Test() DecoratedExpression {
	return i.test
}

func (i *Case) Consequences() []*CaseConsequence {
	return i.cases
}

func (i *Case) DefaultCase() DecoratedExpression {
	return i.defaultCase
}

func (i *Case) DebugString() string {
	return fmt.Sprintf("[dcase]")
}

func (i *Case) FetchPositionAndLength() token.PositionLength {
	return i.test.FetchPositionAndLength()
}
