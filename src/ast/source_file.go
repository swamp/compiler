/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"bytes"
)

type SourceFile struct {
	statements []Expression
	nodes      []Node
}

func (p *SourceFile) Statements() []Expression {
	return p.statements
}

func NewSourceFile(expressions []Expression) *SourceFile {
	for _, expr := range expressions {
		if expr == nil {
			panic("nil expression")
		}
	}
	return &SourceFile{statements: expressions}
}

func (p *SourceFile) String() string {
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

func (p *SourceFile) DebugString() string {
	return expressionArrayToDebugStringEx(p.statements, "\n", 0)
}

func (p *SourceFile) Nodes() []Node {
	return p.nodes
}

func (p *SourceFile) SetNodes(nodes []Node) {
	p.nodes = nodes
}
