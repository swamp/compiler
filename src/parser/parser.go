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

var precedences = map[token.Type]Precedence{
	token.OperatorEqual:          EQUALS,
	token.OperatorNotEqual:       EQUALS,
	token.OperatorLess:           LESSGREATER,
	token.OperatorLessOrEqual:    LESSGREATER,
	token.OperatorGreater:        LESSGREATER,
	token.OperatorGreaterOrEqual: LESSGREATER,
	token.OperatorPlus:           SUM,
	token.OperatorMinus:          SUM,
	token.OperatorAppend:         SUM,
	token.OperatorCons:           SUM,
	token.OperatorDivide:         PRODUCT,
	token.OperatorMultiply:       PRODUCT,
	token.OperatorAnd:            ANDOR,
	token.OperatorOr:             ANDOR,
	token.OperatorBitwiseOr:      ANDOR,
	token.OperatorBitwiseAnd:     ANDOR,
	token.OperatorBitwiseXor:     ANDOR,
	token.OperatorBitwiseNot:     PREFIX,
	token.OperatorUpdate:         UPDATE,
	token.OperatorAssign:         ASSIGN,
	token.OperatorPipeRight:      PIPE,
	token.OperatorPipeLeft:       PIPE,
}

type Parser struct {
	stream *ParseStreamImpl
}

func NewParser(tokenizer *tokenize.Tokenizer, enforceStyle bool) *Parser {
	p := &Parser{}
	p.stream = NewParseStreamImpl(p, tokenizer, enforceStyle)

	return p
}

func (p *Parser) Parse() (*ast.Program, parerr.ParseError) {
	var expressions []ast.Expression

	for !p.stream.tokenizer.MaybeEOF() {
		report, skipLinesErr := p.stream.tokenizer.SkipWhitespaceAllowCommentsToNextIndentation()
		if skipLinesErr != nil {
			return nil, skipLinesErr
		}

		if report.SpacesUntilMaybeNewline > 0 && report.IndentationSpaces > 0 {
			return nil, parerr.NewExtraSpacing(p.stream.positionLength())
		}

		expression, expressionErr := p.parseExpressionStatement(report.Comments)
		if expressionErr != nil {
			return nil, expressionErr
		}

		expressions = append(expressions, expression)

		linesToPad := ast.ExpectedLinePaddingAfter(expression)

		report, mustHaveLineAfterStatementErr := p.stream.eatNewLinesAfterStatement(linesToPad)
		if mustHaveLineAfterStatementErr != nil {
			return nil, parerr.NewExpectedTwoLinesAfterStatement(mustHaveLineAfterStatementErr)
		}
	}

	program := ast.NewProgram(expressions)

	return program, nil
}

func (p *Parser) ParseExpression() (*ast.Program, parerr.ParseError) {
	expressions := []ast.Expression{}

	expression, expressionErr := p.parseExpressionNormal(0)
	if expressionErr != nil {
		return nil, expressionErr
	}
	if expression == nil {
		panic("should not be nil")
	}
	expressions = append(expressions, expression)
	program := ast.NewProgram(expressions)
	return program, nil
}

func (p *Parser) peekUpcomingPrecedence(indentation int) (Precedence, bool) {
	var currentPrecedence Precedence
	saveInfo := p.stream.tokenizer.Tell()
	if _, _, checkingForwardErr := p.stream.maybeUpToOneLogicalSpaceForContinuation(indentation); checkingForwardErr != nil {
		return currentPrecedence, false
	}

	operatorToken, operatorTokenErr := p.stream.readOperatorToken()
	p.stream.tokenizer.Seek(saveInfo)
	if operatorTokenErr != nil {
		return currentPrecedence, false
	}

	currentPrecedence = p.stream.getPrecedenceFromToken(operatorToken)

	return currentPrecedence, true
}

func tokenIsLeaf(tok token.Token) bool {
	t := tok.Type()
	return t == token.VariableSymbol || t == token.TypeSymbol || t == token.NumberInteger || t == token.NumberFixed || t == token.BooleanType || t == token.StringConstant || t == token.ResourceNameSymbol || t == token.TypeIdSymbol || t == token.OperatorUnaryMinus || t == token.OperatorUnaryNot
}

func parenIsLeft(parenToken token.ParenToken) bool {
	t := parenToken.Type()
	return t == token.LeftParen || t == token.LeftBracket || t == token.LeftCurlyBrace || t == token.LeftArrayBracket
}

func (p *Parser) peekIsCall() bool {
	saveInfo := p.stream.tokenizer.Tell()

	wasContinuation, _, spaceErr := p.stream.eatContinuationOrNoSpaceInternal()
	if spaceErr != nil {
		p.stream.tokenizer.Seek(saveInfo)
		return false
	}

	if !wasContinuation {
		p.stream.tokenizer.Seek(saveInfo)
		return false
	}

	anyToken, anyTokenErr := p.stream.readTermToken()
	p.stream.tokenizer.Seek(saveInfo)
	if anyTokenErr != nil {
		return false
	}
	parenToken, wasParen := anyToken.(token.ParenToken)
	isLeftParen := false
	if wasParen {
		isLeftParen = parenIsLeft(parenToken)
	}
	if isLeftParen {
		return true
	}
	wasLeaf := tokenIsLeaf(anyToken)
	return wasLeaf
}

func (p *Parser) internalParseExpression(filterPrecedence Precedence, startIndentation int) (ast.Expression, parerr.ParseError) {
	t, tErr := p.stream.readTermToken()
	if tErr != nil {
		return nil, tErr
	}

	switch tok := t.(type) {
	case token.MultiLineCommentToken:
		{
			return parseMultiLineComment(p.stream, tok)
		}
	case token.SingleLineCommentToken:
		{
			return parseSingleLineComment(p.stream, tok)
		}
	}

	term, termErr := p.parseTermUsingToken(t, startIndentation)
	if termErr != nil {
		switch termErr.(type) {
		case parerr.NotATermError:
			term = nil
		default:
			return nil, termErr
		}
	}

	leftExp := term
	leftExpErr := termErr
	if leftExp == nil {
		// It is not a term, lets try a prefix
		leftExp, leftExpErr = p.parsePrefix(t, startIndentation)
		if leftExpErr != nil {
			return nil, leftExpErr
		}
		if leftExp == nil {
			panic("not allowed")
		}
	}

	_, isTypeIdentifier := term.(*ast.TypeIdentifier)
	_, isVariableIdentifier := term.(*ast.VariableIdentifier)

	if isTypeIdentifier || isVariableIdentifier {
		isCall := p.peekIsCall()
		if isCall {
			_, _, spaceErr := p.stream.maybeOneSpace()
			if spaceErr != nil {
				return nil, spaceErr
			}
			leftExp, leftExpErr = parseFunctionCall(p.stream, startIndentation, leftExp)
		} else {
			typeIdentifier, isTypeIdentifier := term.(*ast.TypeIdentifier)
			if isTypeIdentifier {
				leftExp = ast.NewConstructorCall(typeIdentifier, nil)
			}
		}
	}

	for {
		peekPrecedence, hasUpcomingOperator := p.peekUpcomingPrecedence(startIndentation)
		if !hasUpcomingOperator {
			break
		}
		if filterPrecedence >= peekPrecedence {
			break
		}
		if leftExp == nil {
			panic("not allowed to parse infix with nil")
		}
		leftExp, leftExpErr = p.parseInfix(leftExp, startIndentation)
		if leftExpErr != nil {
			return nil, leftExpErr
		}
	}
	return leftExp, leftExpErr
}

func (p *Parser) parseExpression(precedence Precedence, startIndentation int) (ast.Expression, parerr.ParseError) {
	e, eErr := p.internalParseExpression(precedence, startIndentation)
	return e, eErr
}

func (p *Parser) parseExpressionNormal(startIndentation int) (ast.Expression, parerr.ParseError) {
	return p.parseExpression(LOWEST, startIndentation)
}
