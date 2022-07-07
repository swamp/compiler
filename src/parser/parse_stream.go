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
	addNode(node ast.Node)

	// -----------------------------------------------------------------------------------------------------------------
	// read. Reads the symbol without any spacing
	// -----------------------------------------------------------------------------------------------------------------
	readTypeIdentifier() (*ast.TypeIdentifier, parerr.ParseError)
	readVariableIdentifier() (*ast.VariableIdentifier, parerr.ParseError)
	readKeyword() (token.Keyword, parerr.ParseError)
	readVariableIdentifierAssignOrUpdate(expectedIndentation int) (*ast.VariableIdentifier, bool, int, parerr.ParseError)
	readRightBracket() (token.ParenToken, parerr.ParseError)
	readLeftAngleBracket() (token.ParenToken, parerr.ParseError)
	readRightAngleBracket() (token.ParenToken, parerr.ParseError)
	readRightParen() (token.ParenToken, parerr.ParseError)
	readRightArrayBracket() (token.ParenToken, parerr.ParseError)
	readRightCurly() (token.ParenToken, parerr.ParseError)
	readOf() (token.Keyword, parerr.ParseError)
	readThen() (token.Keyword, parerr.ParseError)
	readElse() (token.Keyword, parerr.ParseError)

	readGuardPipe() (token.GuardToken, parerr.ParseError)

	// -----------------------------------------------------------------------------------------------------------------
	// Spacing
	//
	// Continuation is either a single space or a new line with exactly one indentation.
	//
	// -----------------------------------------------------------------------------------------------------------------
	// maybeOneSpace is usually an optional space after a parenthesis start or end
	maybeOneSpace() (bool, token.IndentationReport, parerr.ParseError)
	eatOneSpace(reason string) (token.IndentationReport, parerr.ParseError)
	eatContinuationReturnIndentation(indentation int) (int, token.IndentationReport, parerr.ParseError)

	// eatCommaSeparatorOrTermination is mostly for items in lists and arrays.
	eatCommaSeparatorOrTermination(indentation int, allowComments tokenize.CommentAllowedType) (bool, token.IndentationReport, parerr.ParseError)

	eatOneSpaceOrDetectEndOfFunctionCallArguments(currentIndentation int) (bool, token.IndentationReport, parerr.ParseError)
	// Block spacing is a helper function for when you specify if single space and new line indent is allowed.
	eatBlockSpacingOneExtraLine(isBlock bool, keywordIndentation int) (token.IndentationReport, parerr.ParseError)

	// eatNewLine is when we require that there should be a new line continuation with exact same indentation.
	eatNewLineContinuationAllowComment(indentation int) (token.IndentationReport, parerr.ParseError)

	// eatNewLineContinuationOrDedent This is used when the blocks that does not allow empty lines. Like the custom type statement.
	eatNewLineContinuationOrDedent(expectedIndentation int) (bool, token.IndentationReport, parerr.ParseError)
	eatOneNewLineContinuationOrDedentAllowComment(expectedIndentation int) (bool, token.IndentationReport, parerr.ParseError)
	// eatOneOrTwoNewLineContinuationOrDedent. This is used when the act of dedenting indicates that the block has ended
	// Mostly used for the `let` block.
	eatOneOrTwoNewLineContinuationOrDedent(indentation int) (bool, token.IndentationReport, parerr.ParseError)

	// -----------------------------------------------------------------------------------------------------------------
	// eat. Similar to read, but discards the result
	// -----------------------------------------------------------------------------------------------------------------

	eatOperatorUpdate() parerr.ParseError
	eatRightArrow() parerr.ParseError
	eatLeftParen() parerr.ParseError
	eatColon() parerr.ParseError
	eatAccessor() parerr.ParseError
	eatAssign() parerr.ParseError

	eatIn() parerr.ParseError

	// -----------------------------------------------------------------------------------------------------------------
	// was. if symbol was detected, it is same as 'read', but returns false without any error if symbol wasn't present.
	// -----------------------------------------------------------------------------------------------------------------
	wasDefaultSymbol() (token.RuneToken, bool)
	wasVariableIdentifier() (*ast.VariableIdentifier, bool)
	wasTypeIdentifier() (*ast.TypeIdentifier, bool)

	// -----------------------------------------------------------------------------------------------------------------
	// maybe. similar to "was*", but only returns a boolean indicating if the symbol was present.
	// -----------------------------------------------------------------------------------------------------------------

	maybeKeywordAlias() (token.Keyword, bool)
	maybeKeywordExposing() (token.Keyword, bool)
	maybeKeywordAs() (token.Keyword, bool)

	maybeAssign() bool
	maybeEllipsis() bool
	maybeAccessor() bool
	maybeRightSquareBracket() (token.ParenToken, bool)
	maybeRightParen() bool
	maybeEmptyParen() bool
	maybeColon() bool
	maybeComma() (token.OperatorToken, bool)
	maybeSpacingAndComma(indentation int) (token.OperatorToken, bool)
	maybePipeLeft() bool
	maybeRightArrow() bool
	maybeOneSpaceAndRightArrow() bool
	maybeLeftParen() (token.ParenToken, bool)
	maybeAsterisk() (token.OperatorToken, bool)
	maybeLeftCurly() (token.ParenToken, bool)
	maybeRightArrayBracket() (token.ParenToken, bool)
	maybeUpToOneLogicalSpaceForContinuation(currentIndentation int) (bool, token.IndentationReport, parerr.ParseError)

	// -----------------------------------------------------------------------------------------------------------------
	// detect. similar to maybe, but doesn't advance the token stream, only reports if the symbol is coming up.
	// -----------------------------------------------------------------------------------------------------------------
	detectNewLine() bool
	detectAssign() bool
	detectNewLineOrSpace() (bool, bool)
	detectOneSpaceAndTermination() bool
	detectTypeIdentifierWithoutScope() bool
	detectNormalOrScopedTypeIdentifier() bool

	// -----------------------------------------------------------------------------------------------------------------
	// parse. parse expressions and terms
	// -----------------------------------------------------------------------------------------------------------------
	parseExpression(precedence Precedence, startIndentation int) (ast.Expression, parerr.ParseError)
	parseExpressionNormal(startIndentation int) (ast.Expression, parerr.ParseError)
	parseExpressionNormalWithComment(startIndentation int, comment token.Comment) (ast.Expression, parerr.ParseError)
	parseTerm(startIndentation int) (ast.Expression, parerr.ParseError)
	parseLiteral(startIndentation int) (ast.Literal, parerr.ParseError)

	// -----------------------------------------------------------------------------------------------------------------
	// debug info
	// -----------------------------------------------------------------------------------------------------------------
	positionLength() token.SourceFileReference
	debugInfo(s string)
	debugInfoRows(s string, rowCount int)
}
