/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package coloring

import (
	"github.com/swamp/compiler/src/token"
)

type Colorer interface {
	Operator(t token.OperatorToken)
	VariableSymbol(t token.VariableSymbolToken)
	Definition(t token.VariableSymbolToken)
	LocalType(t token.VariableSymbolToken)
	Parameter(t token.VariableSymbolToken)
	RecordField(t token.VariableSymbolToken)
	TypeSymbol(t token.TypeSymbolToken)
	TypeGeneratorName(t token.TypeSymbolToken)
	ModuleReference(t token.TypeSymbolToken)
	PrimitiveType(t token.TypeSymbolToken)
	AliasNameSymbol(t token.TypeSymbolToken)
	NumberLiteral(t token.NumberToken)
	BooleanLiteral(t token.BooleanToken)
	KeywordString(t string)
	NewLine(indentation int)
	OperatorString(t string)
	StringLiteral(s token.StringToken)
	OneSpace()
	RightArrow()
	RightPipe()
	LeftPipe()
	String() string
	UserInstruction(t string)
}
