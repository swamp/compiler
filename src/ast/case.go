/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"bytes"
	"fmt"
	"github.com/swamp/compiler/src/token"
)

type CaseConsequence struct {
	variantName *TypeIdentifier
	arguments   []*VariableIdentifier
	expression  Expression
}

func NewCaseConsequence(variantName *TypeIdentifier, arguments []*VariableIdentifier, expression Expression) *CaseConsequence {
	return &CaseConsequence{variantName: variantName, arguments: arguments, expression: expression}
}

func (c *CaseConsequence) Identifier() *TypeIdentifier {
	return c.variantName
}

func (c *CaseConsequence) Arguments() []*VariableIdentifier {
	return c.arguments
}

func (c *CaseConsequence) Expression() Expression {
	return c.expression
}

func (c *CaseConsequence) String() string {
	return fmt.Sprintf("[casecons %v (%v) => %v]", c.variantName, c.arguments, c.expression)
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
	test    Expression
	cases   []*CaseConsequence
	keyword token.Keyword
}

func NewCase(keyword token.Keyword, test Expression, cases []*CaseConsequence) *Case {
	return &Case{keyword: keyword, test: test, cases: cases}
}

func (i *Case) String() string {
	return fmt.Sprintf("[case: %v of %v]", i.test, caseConsequenceArrayToStringEx(i.cases, ";"))
}

func (i *Case) Test() Expression {
	return i.test
}
func (i *Case) Keyword() token.Keyword {
	return i.keyword
}

func (i *Case) PositionLength() token.PositionLength {
	return i.keyword.PositionLength
}

func (i *Case) Consequences() []*CaseConsequence {
	return i.cases
}

func (i *Case) DebugString() string {
	return fmt.Sprintf("[case]")
}
