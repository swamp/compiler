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
	errors []parerr.ParseError
}

func NewParser(tokenizer *tokenize.Tokenizer, enforceStyle bool) *Parser {
	p := &Parser{}
	p.stream = NewParseStreamImpl(p, tokenizer, enforceStyle)

	return p
}

func (p *Parser) Nodes() []ast.Node {
	return p.stream.nodes
}

func writeAllNodes(nodes []ast.Node) {
	lastLine := -1

	for _, node := range nodes {
		rangeFound := node.FetchPositionLength().Range
		if rangeFound.Start().Line() != lastLine {
			fmt.Printf("\n")
			lastLine = rangeFound.End().Line()
		}
		fmt.Printf("node pos %v ", node.FetchPositionLength())
	}
}

func (p *Parser) Parse() (*ast.SourceFile, parerr.ParseError) {
	var statements []ast.Expression

	linesToPad := -2

	for !p.stream.tokenizer.MaybeEOF() {
		report, mustHaveLineAfterStatementErr := p.stream.eatNewLinesAfterStatement(linesToPad)
		if mustHaveLineAfterStatementErr != nil {
			return nil, parerr.NewExpectedTwoLinesAfterStatement(mustHaveLineAfterStatementErr)
		}

		if p.stream.tokenizer.MaybeEOF() {
			break
		}

		if (report.SpacesUntilMaybeNewline > 0 || linesToPad == -1) && report.IndentationSpaces > 0 {
			return nil, parerr.NewExtraSpacing(p.stream.sourceFileReference())
		}
		var previousComment *ast.MultilineComment
		for _, comment := range report.Comments.Comments {
			astMultilineComment := ast.NewMultilineComment(token.NewMultiLineCommentToken(comment.RawString, comment.CommentString, comment.ForDocumentation, comment.SourceFileReference))
			statements = append(statements, astMultilineComment)
			previousComment = astMultilineComment
		}

		expression, expressionErr := p.parseExpressionStatement(previousComment)
		if expressionErr != nil {
			return nil, expressionErr
		}

		statements = append(statements, expression)

		linesToPad = ast.ExpectedLinePaddingAfter(expression)
	}

	program := ast.NewSourceFile(statements)

	program.SetNodes(p.stream.nodes)

	return program, nil
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

func (p *Parser) AddError(parseError parerr.ParseError) {
	p.errors = append(p.errors, parseError)
}

func (p *Parser) parseExpressionNormal(startIndentation int) (ast.Expression, parerr.ParseError) {
	expr, err := p.parseExpression(LOWEST, startIndentation)
	if err != nil {
		return nil, err
	}

	return expr, nil
}
