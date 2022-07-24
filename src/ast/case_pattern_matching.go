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
	comment    token.Comment
}

func NewCaseConsequenceForPatternMatching(index int, literal Literal, expression Expression, comment token.Comment) *CaseConsequencePatternMatching {
	return &CaseConsequencePatternMatching{index: index, literal: literal, expression: expression, comment: comment}
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

func (c *CaseConsequencePatternMatching) Comment() token.Comment {
	return c.comment
}

func (c *CaseConsequencePatternMatching) String() string {
	var literalString string
	if c.literal == nil {
		literalString = "'_'"
	} else {
		literalString = c.literal.String()
	}
	return fmt.Sprintf("[CaseConsPm %v => %v]", literalString, c.expression)
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

type CaseForPatternMatching struct {
	test        Expression
	cases       []*CaseConsequencePatternMatching
	keywordCase token.Keyword
	keywordOf   token.Keyword
	inclusive   token.SourceFileReference
}

func NewCaseForPatternMatching(keywordCase token.Keyword, keywordOf token.Keyword, test Expression, cases []*CaseConsequencePatternMatching) *CaseForPatternMatching {
	inclusive := token.MakeInclusiveSourceFileReference(keywordCase.FetchPositionLength(), cases[len(cases)-1].Expression().FetchPositionLength())

	return &CaseForPatternMatching{keywordCase: keywordCase, keywordOf: keywordOf, test: test, cases: cases, inclusive: inclusive}
}

func (i *CaseForPatternMatching) String() string {
	return fmt.Sprintf("[CasePm %v of %v]", i.test, caseConsequencePatternMatchingArrayToStringEx(i.cases, ";"))
}

func (i *CaseForPatternMatching) Test() Expression {
	return i.test
}

func (i *CaseForPatternMatching) KeywordCase() token.Keyword {
	return i.keywordCase
}

func (i *CaseForPatternMatching) KeywordOf() token.Keyword {
	return i.keywordOf
}

func (i *CaseForPatternMatching) FetchPositionLength() token.SourceFileReference {
	return i.inclusive
}

func (i *CaseForPatternMatching) Consequences() []*CaseConsequencePatternMatching {
	return i.cases
}

func (i *CaseForPatternMatching) DebugString() string {
	return "[CasePm]"
}
