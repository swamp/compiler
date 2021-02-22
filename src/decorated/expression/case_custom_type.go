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

type CaseConsequenceCustomType struct {
	variantName   *ast.TypeIdentifier
	parameters    []*CaseConsequenceParameter
	expression    DecoratedExpression
	internalIndex int
}

func NewCaseConsequenceCustomType(internalIndex int, variantName *ast.TypeIdentifier, parameters []*CaseConsequenceParameter,
	expression DecoratedExpression) *CaseConsequenceCustomType {
	return &CaseConsequenceCustomType{
		internalIndex: internalIndex, variantName: variantName, parameters: parameters,
		expression: expression,
	}
}

func (c *CaseConsequenceCustomType) Expression() DecoratedExpression {
	return c.expression
}

func (c *CaseConsequenceCustomType) InternalIndex() int {
	return c.internalIndex
}

func (c *CaseConsequenceCustomType) Identifier() *ast.TypeIdentifier {
	return c.variantName
}

func (c *CaseConsequenceCustomType) Parameters() []*CaseConsequenceParameter {
	return c.parameters
}

func (c *CaseConsequenceCustomType) String() string {
	return fmt.Sprintf("[dcasecons %v (%v) => %v]", c.variantName, c.parameters, c.expression)
}

func caseConsequenceArrayToStringEx(expressions []*CaseConsequenceCustomType, ch string) string {
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
	test        DecoratedExpression
	cases       []*CaseConsequenceCustomType
	defaultCase DecoratedExpression
}

func NewCaseCustomType(test DecoratedExpression, cases []*CaseConsequenceCustomType, defaultCase DecoratedExpression) (*CaseCustomType, decshared.DecoratedError) {
	return &CaseCustomType{test: test, cases: cases, defaultCase: defaultCase}, nil
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

func (i *CaseCustomType) Test() DecoratedExpression {
	return i.test
}

func (i *CaseCustomType) Consequences() []*CaseConsequenceCustomType {
	return i.cases
}

func (i *CaseCustomType) DefaultCase() DecoratedExpression {
	return i.defaultCase
}

func (i *CaseCustomType) DebugString() string {
	return "[dcase]"
}

func (i *CaseCustomType) FetchPositionLength() token.Range {
	return i.test.FetchPositionLength()
}
