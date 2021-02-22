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
	literal       Expression
	expression    Expression
	internalIndex int
}

func NewCaseConsequencePatternMatching(internalIndex int, literal Expression, expression Expression) *CaseConsequencePatternMatching {
	return &CaseConsequencePatternMatching{internalIndex: internalIndex, literal: literal, expression: expression}
}

func (c *CaseConsequencePatternMatching) Expression() Expression {
	return c.expression
}

func (c *CaseConsequencePatternMatching) InternalIndex() int {
	return c.internalIndex
}

func (c *CaseConsequencePatternMatching) Literal() Expression {
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
	test        Expression
	cases       []*CaseConsequencePatternMatching
	defaultCase Expression
}

func NewCasePatternMatching(test Expression, cases []*CaseConsequencePatternMatching, defaultCase Expression) (*CasePatternMatching, decshared.DecoratedError) {
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
		return fmt.Sprintf("[dpmcase: %v of %v default: %v]", i.test,
			caseConsequencePatternMatchingArrayToStringEx(i.cases, ";"), i.defaultCase)
	}
	return fmt.Sprintf("[dpmcase: %v of %v]", i.test, caseConsequencePatternMatchingArrayToStringEx(i.cases, ";"))
}

func (i *CasePatternMatching) Test() Expression {
	return i.test
}

func (i *CasePatternMatching) Consequences() []*CaseConsequencePatternMatching {
	return i.cases
}

func (i *CasePatternMatching) DefaultCase() Expression {
	return i.defaultCase
}

func (i *CasePatternMatching) DebugString() string {
	return "[dpmcase]"
}

func (i *CasePatternMatching) FetchPositionLength() token.SourceFileReference {
	return i.test.FetchPositionLength()
}
