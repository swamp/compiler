/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"bytes"
	"fmt"

	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type CaseConsequencePatternMatching struct {
	literal    DecoratedExpression
	expression    DecoratedExpression
	internalIndex int
}

func NewCaseConsequencePatternMatching(internalIndex int, literal DecoratedExpression, expression DecoratedExpression) *CaseConsequencePatternMatching {
	return &CaseConsequencePatternMatching{internalIndex: internalIndex, literal: literal, expression: expression}
}

func (c *CaseConsequencePatternMatching) Expression() DecoratedExpression {
	return c.expression
}

func (c *CaseConsequencePatternMatching) InternalIndex() int {
	return c.internalIndex
}

func (c *CaseConsequencePatternMatching) Literal() DecoratedExpression {
	return c.literal
}

func (c *CaseConsequencePatternMatching) String() string {
	return fmt.Sprintf("[dpmcasecons %v => %v]", c.literal, c.expression)
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
	test        DecoratedExpression
	cases       []*CaseConsequencePatternMatching
	defaultCase DecoratedExpression
}

func NewCasePatternMatching(test DecoratedExpression, cases []*CaseConsequencePatternMatching, defaultCase DecoratedExpression) (*CasePatternMatching, decshared.DecoratedError) {
	return &CasePatternMatching{test: test, cases: cases, defaultCase: defaultCase}, nil
}

func (i *CasePatternMatching) Type() dtype.Type {
	if len(i.cases) == 0 {
		return i.defaultCase.Type()
	}
	firstCase := i.cases[0]
	return firstCase.Expression().Type()
}

func (i *CasePatternMatching) String() string {
	if i.defaultCase != nil {
		return fmt.Sprintf("[dpmcase: %v of %v default: %v]", i.test, caseConsequencePatternMatchingArrayToStringEx(i.cases, ";"), i.defaultCase)
	}
	return fmt.Sprintf("[dpmcase: %v of %v]", i.test, caseConsequencePatternMatchingArrayToStringEx(i.cases, ";"))
}

func (i *CasePatternMatching) Test() DecoratedExpression {
	return i.test
}

func (i *CasePatternMatching) Consequences() []*CaseConsequencePatternMatching {
	return i.cases
}

func (i *CasePatternMatching) DefaultCase() DecoratedExpression {
	return i.defaultCase
}

func (i *CasePatternMatching) DebugString() string {
	return fmt.Sprintf("[dpmcase]")
}

func (i *CasePatternMatching) FetchPositionAndLength() token.PositionLength {
	return i.test.FetchPositionAndLength()
}
