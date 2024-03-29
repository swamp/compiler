/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	parerr "github.com/swamp/compiler/src/parser/errors"
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/compiler/src/tokenize"
)

var precedences = map[token.Type]Precedence{
	token.OperatorEqual:             EQUALS,
	token.OperatorNotEqual:          EQUALS,
	token.OperatorLess:              LESSGREATER,
	token.OperatorLessOrEqual:       LESSGREATER,
	token.OperatorGreater:           LESSGREATER,
	token.OperatorGreaterOrEqual:    LESSGREATER,
	token.OperatorPlus:              SUM,
	token.OperatorMinus:             SUM,
	token.OperatorAppend:            SUM,
	token.OperatorCons:              SUM,
	token.Colon:                     ASSIGN,
	token.OperatorDivide:            PRODUCT,
	token.OperatorMultiply:          PRODUCT,
	token.OperatorRemainder:         PRODUCT,
	token.OperatorAnd:               ANDOR,
	token.OperatorOr:                ANDOR,
	token.OperatorBitwiseOr:         ANDOR,
	token.OperatorBitwiseAnd:        ANDOR,
	token.OperatorBitwiseXor:        ANDOR,
	token.OperatorBitwiseShiftLeft:  ANDOR,
	token.OperatorBitwiseShiftRight: ANDOR,
	token.OperatorBitwiseNot:        PREFIX,
	token.OperatorUpdate:            UPDATE,
	token.OperatorAssign:            ASSIGN,
	token.OperatorPipeRight:         PIPE,
	token.OperatorPipeLeft:          PIPE,
}

type Parser struct {
	stream          *ParseStreamImpl
	previousComment token.Comment
}

func NewParser(tokenizer *tokenize.Tokenizer, enforceStyle bool) *Parser {
	p := &Parser{}
	p.stream = NewParseStreamImpl(p, tokenizer, enforceStyle)

	return p
}

func (p *Parser) Nodes() []ast.Node {
	return p.stream.nodes
}

func (p *Parser) Errors() parerr.ParseError {
	return p.stream.Errors()
}

func (p *Parser) Parse() (*ast.SourceFile, parerr.ParseError) {
	var statements []ast.Expression

	linesToPadMin := -2
	linesToPadMax := -2

	var errors parerr.ParseError

	for !p.stream.tokenizer.MaybeEOF() {
		report, mustHaveLineAfterStatementErr := p.stream.eatNewLinesAfterStatement(linesToPadMin, linesToPadMax)
		if mustHaveLineAfterStatementErr != nil {
			p.stream.debugInfo(fmt.Sprintf("problem was here %d %d (%d)", linesToPadMin, linesToPadMax,
				report.NewLineCount))
			//errors = parerr.AppendError(errors, parerr.NewExpectedTwoLinesAfterStatement(mustHaveLineAfterStatementErr))
		}

		if p.stream.tokenizer.MaybeEOF() {
			break
		}

		if (report.SpacesUntilMaybeNewline > 0 || linesToPadMin == -1) && report.IndentationSpaces > 0 {
			errors = parerr.AppendError(errors, parerr.NewExtraSpacing(p.stream.sourceFileReference()))
		}

		astMultilineComments := ast.CommentBlockToAst(report.Comments)
		if len(astMultilineComments) > 1 {
			for _, comment := range astMultilineComments[:len(astMultilineComments)-1] {
				statements = append(statements, comment)
			}
		}

		var lastComment *ast.MultilineComment
		if len(astMultilineComments) > 0 {
			lastComment = astMultilineComments[len(astMultilineComments)-1]
		}

		expression, expressionErr := p.parseExpressionStatement(lastComment)
		if expressionErr != nil {
			if IsCompileErr(expressionErr) {
				return nil, expressionErr
			}
			errors = parerr.AppendError(errors, expressionErr)
		}

		if IsCompileError(expressionErr) {
			return nil, errors
		}

		statements = append(statements, expression)

		linesToPadMin, linesToPadMax = ast.ExpectedLinePaddingAfter(expression)
	}

	program := ast.NewSourceFile(statements)

	program.SetNodes(p.stream.nodes)

	return program, errors
}

func (p *Parser) ParseExpression() (*ast.SourceFile, parerr.ParseError) {
	var expressions []ast.Expression

	expression, expressionErr := p.parseExpressionNormal(0)
	if expressionErr != nil {
		return nil, expressionErr
	}
	if expression == nil {
		panic("should not be nil")
	}
	expressions = append(expressions, expression)
	program := ast.NewSourceFile(expressions)
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
	return t == token.VariableSymbol || t == token.TypeSymbol || t == token.NumberInteger || t == token.NumberFixed || t == token.BooleanType || t == token.StringConstant || t == token.StringInterpolationTupleConstant || t == token.StringInterpolationStringConstant || t == token.ResourceNameSymbol || t == token.TypeIdSymbol || t == token.OperatorUnaryMinus || t == token.OperatorUnaryNot
}

func parenIsLeft(parenToken token.ParenToken) bool {
	t := parenToken.Type()
	return t == token.LeftParen || t == token.LeftSquareBracket || t == token.LeftCurlyBrace || t == token.LeftArrayBracket
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

func (p *Parser) internalParseExpression(filterPrecedence Precedence, startIndentation int) (ast.Expression,
	parerr.ParseError) {
	t, tErr := p.stream.readTermToken()
	if tErr != nil {
		return nil, tErr
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
	// log.Printf("term token %T %v", term, term)

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

	p.stream.nodes = append(p.stream.nodes, leftExp)
	notScopedTypeIdent, isTypeIdentifier := term.(*ast.TypeIdentifier)
	scopedTypeIdent, isScopedTypeIdentifier := term.(*ast.TypeIdentifierScoped)
	_, isVariableIdentifier := term.(*ast.VariableIdentifier)
	_, isScopedVariableIdentifier := term.(*ast.VariableIdentifierScoped)

	if isTypeIdentifier || isScopedTypeIdentifier || isVariableIdentifier || isScopedVariableIdentifier {
		isCall := p.peekIsCall()
		if isCall {
			_, _, spaceErr := p.stream.maybeOneSpace()
			if spaceErr != nil {
				return nil, spaceErr
			}
			leftExp, leftExpErr = parseFunctionCall(p.stream, startIndentation, leftExp)
		} else {
			if isTypeIdentifier || isScopedTypeIdentifier {
				// This happens when there is a constructor call which takes no arguments
				// Usually a custom type variant constructor
				var someRef ast.TypeReferenceScopedOrNormal
				if isScopedTypeIdentifier {
					someRef = ast.NewScopedTypeReference(scopedTypeIdent, nil)
				} else {
					someRef = ast.NewTypeReference(notScopedTypeIdent, nil)
				}
				leftExp = ast.NewConstructorCall(someRef, nil)
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
	expr, err := p.parseExpression(LOWEST, startIndentation)
	if err != nil {
		return nil, err
	}

	return expr, nil
}

func (p *Parser) parseExpressionNormalWithComment(startIndentation int, comment token.Comment) (ast.Expression,
	parerr.ParseError) {
	p.previousComment = comment
	expr, err := p.parseExpressionNormal(startIndentation)
	p.previousComment = nil
	if err != nil {
		return nil, err
	}

	return expr, nil
}
