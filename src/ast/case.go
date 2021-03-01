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

type CaseConsequenceCustomType struct {
	variantName *TypeIdentifier
	arguments   []*VariableIdentifier
	expression  Expression
}

func NewCaseConsequenceForCustomType(variantName *TypeIdentifier, arguments []*VariableIdentifier, expression Expression) *CaseConsequenceCustomType {
	return &CaseConsequenceCustomType{variantName: variantName, arguments: arguments, expression: expression}
}

func (c *CaseConsequenceCustomType) Identifier() *TypeIdentifier {
	return c.variantName
}

func (c *CaseConsequenceCustomType) Arguments() []*VariableIdentifier {
	return c.arguments
}

func (c *CaseConsequenceCustomType) Expression() Expression {
	return c.expression
}

func (c *CaseConsequenceCustomType) String() string {
	return fmt.Sprintf("[casecons %v (%v) => %v]", c.variantName, c.arguments, c.expression)
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
	test        Expression
	cases       []*CaseConsequenceCustomType
	keywordCase token.Keyword
	keywordOf   token.Keyword
	inclusive   token.SourceFileReference
}

func NewCaseForCustomType(keywordCase token.Keyword, keywordOf token.Keyword, test Expression, cases []*CaseConsequenceCustomType) *CaseCustomType {
	inclusive := token.MakeInclusiveSourceFileReference(keywordCase.FetchPositionLength(), cases[len(cases)-1].Expression().FetchPositionLength())
	return &CaseCustomType{keywordCase: keywordCase, keywordOf: keywordOf, test: test, cases: cases, inclusive: inclusive}
}

func (i *CaseCustomType) String() string {
	return fmt.Sprintf("[case: %v of %v]", i.test, caseConsequenceArrayToStringEx(i.cases, ";"))
}

func (i *CaseCustomType) Test() Expression {
	return i.test
}

func (i *CaseCustomType) KeywordCase() token.Keyword {
	return i.keywordCase
}

func (i *CaseCustomType) KeywordOf() token.Keyword {
	return i.keywordOf
}

func (i *CaseCustomType) FetchPositionLength() token.SourceFileReference {
	return i.inclusive
}

func (i *CaseCustomType) Consequences() []*CaseConsequenceCustomType {
	return i.cases
}

func (i *CaseCustomType) DebugString() string {
	return "[case]"
}
