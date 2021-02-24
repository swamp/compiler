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
}

func NewIfExpression(condition Expression, consequence Expression, alternative Expression) *IfExpression {
	return &IfExpression{condition: condition, consequence: consequence, alternative: alternative}
}

func (i *IfExpression) Condition() Expression {
	return i.condition
}

func (i *IfExpression) Consequence() Expression {
	return i.consequence
}

func (i *IfExpression) Alternative() Expression {
	return i.alternative
}

func (i *IfExpression) FetchPositionLength() token.SourceFileReference {
	return i.consequence.FetchPositionLength()
}

func (i *IfExpression) String() string {
	return fmt.Sprintf("[if: %v then %v else %v]", i.condition, i.consequence, i.alternative)
}

func (i *IfExpression) DebugString() string {
	return fmt.Sprintf("[if]")
}
