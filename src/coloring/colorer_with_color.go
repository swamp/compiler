/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package coloring

import (
	"strings"

	"github.com/swamp/compiler/src/token"
)

type ColorerWithColor struct {
	builder strings.Builder
}

func NewColorerWithColor() *ColorerWithColor {
	return &ColorerWithColor{}
}

func (c *ColorerWithColor) Operator(t token.OperatorToken) {
	c.builder.WriteString(ColorOperator(t))
}

func (c *ColorerWithColor) VariableSymbol(t token.VariableSymbolToken) {
	c.builder.WriteString(ColorVariableSymbol(t))
}

func (c *ColorerWithColor) Parameter(t token.VariableSymbolToken) {
	c.builder.WriteString(ColorVariableSymbol(t))
}

func (c *ColorerWithColor) Definition(t token.VariableSymbolToken) {
	c.builder.WriteString(ColorDefinition(t))
}

func (c *ColorerWithColor) TypeSymbol(t token.TypeSymbolToken) {
	c.builder.WriteString(ColorTypeSymbol(t))
}

func (c *ColorerWithColor) ModuleReference(t token.TypeSymbolToken) {
	c.builder.WriteString(ColorModuleReference(t))
}

func (c *ColorerWithColor) TypeGeneratorName(t token.TypeSymbolToken) {
	c.builder.WriteString(ColorTypeSymbol(t))
}

func (c *ColorerWithColor) PrimitiveType(t token.TypeSymbolToken) {
	c.builder.WriteString(ColorPrimitiveType(t))
}

func (c *ColorerWithColor) AliasNameSymbol(t token.TypeSymbolToken) {
	c.builder.WriteString(ColorTypeSymbol(t))
}

func (c *ColorerWithColor) NumberLiteral(t token.NumberToken) {
	c.builder.WriteString(ColorNumberLiteral(t))
}

func (c *ColorerWithColor) BooleanLiteral(t token.BooleanToken) {
	c.builder.WriteString(colorBooleanLiteral(t))
}

func (c *ColorerWithColor) StringLiteral(t token.StringToken) {
	c.builder.WriteString("\"" + ColorKeywordString(t.Text()) + "\"")
}

func (c *ColorerWithColor) KeywordString(t string)  {
	c.builder.WriteString(ColorKeywordString(t))
}

func (c *ColorerWithColor) NewLine(indentation int)  {
	c.builder.WriteString(writeNewline(indentation))
}

func (c *ColorerWithColor) RecordField(t token.VariableSymbolToken) {
	c.builder.WriteString(ColorRecordField(t))
}

func (c *ColorerWithColor) OneSpace()  {
	c.builder.WriteString(" ")
}

func (c *ColorerWithColor) RightArrow()  {
	c.KeywordString("➞")
}

func (c *ColorerWithColor) LeftPipe()  {
	c.KeywordString("ᐊ")
}

func (c *ColorerWithColor) RightPipe()  {
	c.KeywordString("ᐅ")
}

func (c *ColorerWithColor) OperatorString(operator string)  {
	c.builder.WriteString(ColorOperatorString(operator))
}

func (c *ColorerWithColor) LocalType(operator token.VariableSymbolToken)  {
	c.builder.WriteString(ColorLocalType(operator))
}

func (c *ColorerWithColor) UserInstruction(t string) {
	c.builder.WriteString(ColorKeywordString(t))
}

func (c *ColorerWithColor) String() string {
	return c.builder.String()
}
