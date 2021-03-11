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
	"github.com/swamp/compiler/src/token"
)

type CaseConsequenceForPatternMatching struct {
	literal        Expression
	expression     Expression
	astConsequence *ast.CaseConsequencePatternMatching
	internalIndex  int
}

func NewCaseConsequencePatternMatching(astConsequence *ast.CaseConsequencePatternMatching, internalIndex int, literal Expression, expression Expression) *CaseConsequenceForPatternMatching {
	return &CaseConsequenceForPatternMatching{internalIndex: internalIndex, astConsequence: astConsequence, literal: literal, expression: expression}
}

func (c *CaseConsequenceForPatternMatching) Expression() Expression {
	return c.expression
}

func (c *CaseConsequenceForPatternMatching) InternalIndex() int {
	return c.internalIndex
}

func (c *CaseConsequenceForPatternMatching) Literal() Expression {
	return c.literal
}

func (c *CaseConsequenceForPatternMatching) String() string {
	return fmt.Sprintf("[dpmcasecons %v => %v]", c.literal, c.expression)
}

func caseConsequencePatternMatchingArrayToStringEx(expressions []*CaseConsequenceForPatternMatching, ch string) string {
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
	cases       []*CaseConsequenceForPatternMatching
	defaultCase Expression
	astCase     *ast.CaseForPatternMatching
}

func NewCaseForPatternMatching(astCase *ast.CaseForPatternMatching, test Expression, cases []*CaseConsequenceForPatternMatching, defaultCase Expression) (*CaseForPatternMatching, decshared.DecoratedError) {
	return &CaseForPatternMatching{astCase: astCase, test: test, cases: cases, defaultCase: defaultCase}, nil
}

func (i *CaseForPatternMatching) Type() dtype.Type {
	if len(i.cases) == 0 {
		return i.defaultCase.Type()
	}
	firstCase := i.cases[0]
	return firstCase.Expression().Type()
}

func (i *CaseForPatternMatching) AstCasePatternMatching() *ast.CaseForPatternMatching {
	return i.astCase
}

func (i *CaseForPatternMatching) String() string {
	if i.defaultCase != nil {
		return fmt.Sprintf("[dpmcase: %v of %v default: %v]", i.test,
			caseConsequencePatternMatchingArrayToStringEx(i.cases, ";"), i.defaultCase)
	}
	return fmt.Sprintf("[dpmcase: %v of %v]", i.test, caseConsequencePatternMatchingArrayToStringEx(i.cases, ";"))
}

func (i *CaseForPatternMatching) Test() Expression {
	return i.test
}

func (i *CaseForPatternMatching) Consequences() []*CaseConsequenceForPatternMatching {
	return i.cases
}

func (i *CaseForPatternMatching) DefaultCase() Expression {
	return i.defaultCase
}

func (i *CaseForPatternMatching) DebugString() string {
	return "[dpmcase]"
}

func (i *CaseForPatternMatching) FetchPositionLength() token.SourceFileReference {
	return i.astCase.FetchPositionLength()
}
