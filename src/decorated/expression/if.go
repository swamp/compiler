/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type If struct {
	condition   DecoratedExpression
	consequence DecoratedExpression
	alternative DecoratedExpression
}

func (l *If) Type() dtype.Type {
	return l.consequence.Type()
}

func NewIf(condition DecoratedExpression, consequence DecoratedExpression, alternative DecoratedExpression) *If {
	return &If{condition: condition, consequence: consequence, alternative: alternative}
}

func (l *If) String() string {
	return fmt.Sprintf("[if %v then %v else %v]", l.condition, l.consequence, l.alternative)
}

func (l *If) Condition() DecoratedExpression {
	return l.condition
}

func (l *If) Consequence() DecoratedExpression {
	return l.consequence
}

func (l *If) Alternative() DecoratedExpression {
	return l.alternative
}

func (l *If) FetchPositionLength() token.SourceFileReference {
	return l.condition.FetchPositionLength()
}
