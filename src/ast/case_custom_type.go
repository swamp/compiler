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

type CaseConsequenceForCustomType struct {
	variantName *TypeIdentifier
	arguments   []*VariableIdentifier
	expression  Expression
	comment     token.Comment
}

func NewCaseConsequenceForCustomType(variantName *TypeIdentifier, arguments []*VariableIdentifier, expression Expression, comment token.Comment) *CaseConsequenceForCustomType {
	return &CaseConsequenceForCustomType{variantName: variantName, arguments: arguments, expression: expression, comment: comment}
}

func (c *CaseConsequenceForCustomType) Identifier() *TypeIdentifier {
	return c.variantName
}

func (c *CaseConsequenceForCustomType) Arguments() []*VariableIdentifier {
	return c.arguments
}

func (c *CaseConsequenceForCustomType) Expression() Expression {
	return c.expression
}

func (c *CaseConsequenceForCustomType) Comment() token.Comment {
	return c.comment
}

func (c *CaseConsequenceForCustomType) String() string {
	return fmt.Sprintf("[casecons %v (%v) => %v]", c.variantName, c.arguments, c.expression)
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

type CaseForCustomType struct {
	test        Expression
	cases       []*CaseConsequenceForCustomType
	keywordCase token.Keyword
	keywordOf   token.Keyword
	inclusive   token.SourceFileReference
}

func NewCaseForCustomType(keywordCase token.Keyword, keywordOf token.Keyword, test Expression, cases []*CaseConsequenceForCustomType) *CaseForCustomType {
	inclusive := token.MakeInclusiveSourceFileReference(keywordCase.FetchPositionLength(), cases[len(cases)-1].Expression().FetchPositionLength())
	return &CaseForCustomType{keywordCase: keywordCase, keywordOf: keywordOf, test: test, cases: cases, inclusive: inclusive}
}

func (i *CaseForCustomType) String() string {
	return fmt.Sprintf("[case: %v of %v]", i.test, caseConsequenceArrayToStringEx(i.cases, ";"))
}

func (i *CaseForCustomType) Test() Expression {
	return i.test
}

func (i *CaseForCustomType) KeywordCase() token.Keyword {
	return i.keywordCase
}

func (i *CaseForCustomType) KeywordOf() token.Keyword {
	return i.keywordOf
}

func (i *CaseForCustomType) FetchPositionLength() token.SourceFileReference {
	return i.inclusive
}

func (i *CaseForCustomType) Consequences() []*CaseConsequenceForCustomType {
	return i.cases
}

func (i *CaseForCustomType) DebugString() string {
	return "[case]"
}
