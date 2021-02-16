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

type CaseConsequencePatternMatching struct {
	literal    Literal
	expression Expression
	index      int
}

func NewCaseConsequenceForPatternMatching(index int, literal Literal, expression Expression) *CaseConsequencePatternMatching {
	return &CaseConsequencePatternMatching{index: index, literal: literal, expression: expression}
}

func (c *CaseConsequencePatternMatching) Literal() Literal {
	return c.literal
}

func (c *CaseConsequencePatternMatching) Index() int {
	return c.index
}

func (c *CaseConsequencePatternMatching) Expression() Expression {
	return c.expression
}

func (c *CaseConsequencePatternMatching) String() string {
	return fmt.Sprintf("[caseconspm %v => %v]", c.literal, c.expression)
}

func caseConsequencePatternMatchingArrayToStringEx(expressions []*CaseConsequencePatternMatching, ch string) string {
	var out bytes.Buffer

	for index, expression := range expressions {
		if index > 0 {
			out.WriteString(ch)
		}
		out.WriteString(expression.String())
	}
	return out.String()
}

type CasePatternMatching struct {
	test    Expression
	cases   []*CaseConsequencePatternMatching
	keyword token.Keyword
}

func NewCaseForPatternMatching(keyword token.Keyword, test Expression, cases []*CaseConsequencePatternMatching) *CasePatternMatching {
	return &CasePatternMatching{keyword: keyword, test: test, cases: cases}
}

func (i *CasePatternMatching) String() string {
	return fmt.Sprintf("[casepm: %v of %v]", i.test, caseConsequencePatternMatchingArrayToStringEx(i.cases, ";"))
}

func (i *CasePatternMatching) Test() Expression {
	return i.test
}

func (i *CasePatternMatching) Keyword() token.Keyword {
	return i.keyword
}

func (i *CasePatternMatching) PositionLength() token.PositionLength {
	return i.keyword.PositionLength
}

func (i *CasePatternMatching) Consequences() []*CaseConsequencePatternMatching {
	return i.cases
}

func (i *CasePatternMatching) DebugString() string {
	return "[casepm]"
}
