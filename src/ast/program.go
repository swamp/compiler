/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"bytes"
)

type Program struct {
	statements []Expression
}

func (p *Program) Statements() []Expression {
	return p.statements
}

func NewProgram(expressions []Expression) *Program {
	for _, expr := range expressions {
		if expr == nil {
			panic("nil expression")
		}
	}
	return &Program{statements: expressions}
}

func (p *Program) String() string {
	return expressionStatementsToString(p.statements)
}

func expressionArrayToDebugStringEx(expressions []Expression, ch string, level int) string {
	var out bytes.Buffer

	for index, expression := range expressions {
		if index > 0 {
			out.WriteString(ch)
		}
		out.WriteString(expression.DebugString())

		if level < 1 {
			infix, isInfix := expression.(*BinaryOperator)
			if isInfix {
				sub := []Expression{infix.Left(), infix.Right()}
				expressionArrayToDebugStringEx(sub, "\n", level+1)
			}
		}
	}
	return out.String()
}

func (p *Program) DebugString() string {
	return expressionArrayToDebugStringEx(p.statements, "\n", 0)
}
