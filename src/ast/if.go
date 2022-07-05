/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type IfExpression struct {
	condition   Expression
	consequence Expression
	alternative Expression
	inclusive   token.SourceFileReference
	keywordIf   token.Keyword
	keywordThen token.Keyword
	keywordElse token.Keyword
}

func NewIfExpression(keywordIf token.Keyword, keywordThen token.Keyword, keywordElse token.Keyword, condition Expression, consequence Expression, alternative Expression) *IfExpression {
	inclusive := token.MakeInclusiveSourceFileReference(condition.FetchPositionLength(), alternative.FetchPositionLength())
	return &IfExpression{inclusive: inclusive, keywordElse: keywordElse, keywordIf: keywordIf, keywordThen: keywordThen, condition: condition, consequence: consequence, alternative: alternative}
}

func (i *IfExpression) Condition() Expression {
	return i.condition
}

func (i *IfExpression) Consequence() Expression {
	return i.consequence
}

func (i *IfExpression) KeywordIf() token.Keyword {
	return i.keywordIf
}

func (i *IfExpression) KeywordThen() token.Keyword {
	return i.keywordThen
}

func (i *IfExpression) KeywordElse() token.Keyword {
	return i.keywordElse
}

func (i *IfExpression) Alternative() Expression {
	return i.alternative
}

func (i *IfExpression) FetchPositionLength() token.SourceFileReference {
	return i.inclusive
}

func (i *IfExpression) String() string {
	return fmt.Sprintf("[if: %v then %v else %v]", i.condition, i.consequence, i.alternative)
}

func (i *IfExpression) DebugString() string {
	return "[if]"
}
