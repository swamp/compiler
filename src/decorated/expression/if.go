/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type If struct {
	condition   Expression
	consequence Expression
	alternative Expression
	astIf       *ast.IfExpression
}

func (l *If) Type() dtype.Type {
	return l.consequence.Type()
}

func NewIf(astIf *ast.IfExpression, condition Expression, consequence Expression, alternative Expression) *If {
	return &If{astIf: astIf, condition: condition, consequence: consequence, alternative: alternative}
}

func (l *If) String() string {
	return fmt.Sprintf("[if %v then %v else %v]", l.condition, l.consequence, l.alternative)
}

func (l *If) AstIf() *ast.IfExpression {
	return l.astIf
}

func (l *If) Condition() Expression {
	return l.condition
}

func (l *If) Consequence() Expression {
	return l.consequence
}

func (l *If) Alternative() Expression {
	return l.alternative
}

func (l *If) FetchPositionLength() token.SourceFileReference {
	return l.astIf.FetchPositionLength()
}
