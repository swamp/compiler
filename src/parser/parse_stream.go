/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"github.com/swamp/compiler/src/ast"
	parerr "github.com/swamp/compiler/src/parser/errors"
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/compiler/src/tokenize"
)

type ParseStream interface {
	addWarning(description string, length token.PositionLength)

	// -----------------------------------------------------------------------------------------------------------------
	// read. Reads the symbol without any spacing
	// -----------------------------------------------------------------------------------------------------------------
	readTypeIdentifier() (*ast.TypeIdentifier, parerr.ParseError)
	readVariableIdentifier() (*ast.VariableIdentifier, parerr.ParseError)
	readVariableIdentifierAssignOrUpdate() (*ast.VariableIdentifier, bool, int, parerr.ParseError)

	// -----------------------------------------------------------------------------------------------------------------
	// eat. Similar to read, but discards the result
	// -----------------------------------------------------------------------------------------------------------------

	// Continuation is either a single space or a new line with exactly one indentation.
	eatOneSpace(reason string) (token.IndentationReport, parerr.ParseError)
	eatOneSpaceOrIndent(reason string) (token.IndentationReport, parerr.ParseError)
	eatContinuationReturnIndentation(indentation int) (int, token.IndentationReport, parerr.ParseError)
	eatArgumentSpaceOrDetectEndOfArguments(currentIndentation int) (bool, token.IndentationReport, parerr.ParseError)

	eatCommaSeparatorOrTermination(indentation int, allowComments tokenize.CommentAllowedType) (bool, token.IndentationReport, parerr.ParseError)

	eatNewLineContinuation(indentation int) (token.IndentationReport, parerr.ParseError)
	eatNewLineContinuationAllowComment(indentation int) (token.IndentationReport, parerr.ParseError)
	eatOneOrTwoNewLineContinuationOrDedent(indentation int) (bool, token.IndentationReport, parerr.ParseError)

	// Block spacing is a helper function for when you both allow single space and new line continuations.
	eatBlockSpacing(isBlock bool, keywordIndentation int) (token.IndentationReport, parerr.ParseError)
	eatBlockSpacingOneExtraLine(isBlock bool, keywordIndentation int) (token.IndentationReport, parerr.ParseError)

	eatOperatorUpdate() parerr.ParseError
	eatRightArrow() parerr.ParseError
	eatLeftParen() parerr.ParseError
	eatRightParen() parerr.ParseError
	eatRightCurly() parerr.ParseError
	eatRightBracket() parerr.ParseError
	eatColon() parerr.ParseError
	eatAccessor() parerr.ParseError
	eatAssign() parerr.ParseError

	eatOf() parerr.ParseError
	eatThen() parerr.ParseError
	eatElse() parerr.ParseError
	eatIn() parerr.ParseError

	// -----------------------------------------------------------------------------------------------------------------
	// was. it symbol was detected, it is same as 'eat', but returns false without any error if symbol wasn't present.
	// -----------------------------------------------------------------------------------------------------------------
	wasDefaultSymbol() (token.RuneToken, bool)
	wasVariableIdentifier() (*ast.VariableIdentifier, bool)
	wasTypeIdentifier() (*ast.TypeIdentifier, bool)

	// -----------------------------------------------------------------------------------------------------------------
	// maybe. similar to "was*", but only returns a boolean indicating if the symbol was present.
	// -----------------------------------------------------------------------------------------------------------------
	maybeOneSpace() (bool, token.IndentationReport, parerr.ParseError)
	maybeNewLineContinuation(expectedIndentation int) (bool, token.IndentationReport, parerr.ParseError)
	maybeNewLineContinuationAllowComment(expectedIndentation int) (bool, token.IndentationReport, parerr.ParseError)
	maybeNewLineContinuationWithExtraEmptyLine(expectedIndentation int) (bool, token.IndentationReport, parerr.ParseError)

	maybeKeywordAlias() bool
	maybeKeywordExposing() bool
	maybeKeywordAs() bool

	maybeAssign() bool
	maybeEllipsis() bool
	maybeAccessor() bool
	maybeRightBracket() bool
	maybeRightParen() bool
	maybeEmptyParen() bool
	maybeColon() bool
	maybePipeLeft() bool
	maybeRightArrow() bool
	maybeOneSpaceAndRightArrow() bool
	maybeLeftParen() bool
	maybeLeftCurly() bool

	// -----------------------------------------------------------------------------------------------------------------
	// detect. similar to maybe, but doesn't advance the token stream, only reports if the symbol is coming up.
	// -----------------------------------------------------------------------------------------------------------------
	detectNewLine() bool
	detectAssign() bool
	detectNewLineOrSpace() (bool, bool)
	detectOneSpaceAndTermination() bool
	detectTypeIdentifier() bool

	// -----------------------------------------------------------------------------------------------------------------
	// parse. parse expressions and terms
	// -----------------------------------------------------------------------------------------------------------------
	parseExpression(precedence Precedence, startIndentation int) (ast.Expression, parerr.ParseError)
	parseExpressionNormal(startIndentation int) (ast.Expression, parerr.ParseError)
	parseTerm(startIndentation int) (ast.Expression, parerr.ParseError)
	parseLiteral(startIndentation int) (ast.Literal, parerr.ParseError)

	// -----------------------------------------------------------------------------------------------------------------
	// debug info
	// -----------------------------------------------------------------------------------------------------------------
	positionLength() token.PositionLength
	debugInfo(s string)
	debugInfoRows(s string, rowCount int)
}
