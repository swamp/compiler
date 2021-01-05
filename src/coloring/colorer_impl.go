/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package coloring

import (
	"strings"

	"github.com/swamp/compiler/src/token"
)

type ColorerWithoutColor struct {
	builder strings.Builder
}

func NewColorerWithoutColor() *ColorerWithoutColor {
	return &ColorerWithoutColor{}
}

func (c *ColorerWithoutColor) Operator(t token.OperatorToken) {
	c.builder.WriteString(t.Raw())
}

func (c *ColorerWithoutColor) VariableSymbol(t token.VariableSymbolToken) {
	c.builder.WriteString(t.Raw())
}

func (c *ColorerWithoutColor) RecordField(t token.VariableSymbolToken) {
	c.builder.WriteString(t.Raw())
}

func (c *ColorerWithoutColor) Parameter(t token.VariableSymbolToken) {
	c.builder.WriteString(t.Raw())
}

func (c *ColorerWithoutColor) Definition(t token.VariableSymbolToken) {
	c.builder.WriteString(t.Raw())
}

func (c *ColorerWithoutColor) TypeSymbol(t token.TypeSymbolToken) {
	c.builder.WriteString(t.Raw())
}

func (c *ColorerWithoutColor) ModuleReference(t token.TypeSymbolToken) {
	c.builder.WriteString(t.Raw())
}

func (c *ColorerWithoutColor) TypeGeneratorName(t token.TypeSymbolToken) {
	c.builder.WriteString(t.Raw())
}

func (c *ColorerWithoutColor) PrimitiveType(t token.TypeSymbolToken) {
	c.builder.WriteString(t.Raw())
}

func (c *ColorerWithoutColor) AliasNameSymbol(t token.TypeSymbolToken) {
	c.builder.WriteString(t.Raw())
}

func (c *ColorerWithoutColor) NumberLiteral(t token.NumberToken) {
	c.builder.WriteString(t.Raw())
}

func (c *ColorerWithoutColor) BooleanLiteral(t token.BooleanToken) {
	c.builder.WriteString(t.Raw())
}

func (c *ColorerWithoutColor) StringLiteral(t token.StringToken) {
	c.builder.WriteString("\"" + t.Text() + "\"")
}

func (c *ColorerWithoutColor) LocalType(t token.VariableSymbolToken) {
	c.builder.WriteString(t.Raw())
}

func (c *ColorerWithoutColor) UserInstruction(t string) {
	c.builder.WriteString(t)
}

func (c *ColorerWithoutColor) KeywordString(t string)  {
	c.builder.WriteString(t)
}

func writeIndentation(indentation int) string {
	return strings.Repeat("    ", indentation)
}

func writeNewline(indentation int) string {
	return "\n" + writeIndentation(indentation)
}

func (c *ColorerWithoutColor) NewLine(indentation int)  {
	c.builder.WriteString(writeNewline(indentation))
}

func (c *ColorerWithoutColor) OneSpace()  {
	c.builder.WriteString(" ")
}


func (c *ColorerWithoutColor) RightArrow()  {
	c.KeywordString("->")
}

func (c *ColorerWithoutColor) LeftPipe()  {
	c.KeywordString("<|")
}

func (c *ColorerWithoutColor) RightPipe()  {
	c.KeywordString("|>")
}


func (c *ColorerWithoutColor) OperatorString(operator string)  {
	c.builder.WriteString(operator)
}

func (c *ColorerWithoutColor) String() string {
	return c.builder.String()
}
