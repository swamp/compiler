/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"regexp"
	"strings"

	"github.com/swamp/compiler/src/ast"
	parerr "github.com/swamp/compiler/src/parser/errors"
	"github.com/swamp/compiler/src/runestream"
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/compiler/src/tokenize"
)

const interpolationRegExp = `\{.*?\}`

func parseInterpolationString(s string) [][]int {
	re, reErr := regexp.Compile(interpolationRegExp)
	if reErr != nil {
		panic(reErr)
	}

	locations := re.FindAllStringIndex(s, -1)

	return locations
}

func replaceInterpolationString(stringToken token.StringToken) ([]ast.Expression, parerr.ParseError) {
	ranges := parseInterpolationString(stringToken.Text())

	var expressions []ast.Expression
	lastPos := 0

	for _, item := range ranges {
		start := item[0]
		end := item[1]

		stringPart := stringToken.Text()[lastPos:start]
		if len(stringPart) > 0 {
			sourceFileReference := token.SourceFileReference{
				Range: token.MakeRange(
					stringToken.Range.Start(),
					stringToken.Range.End(),
				),
				Document: stringToken.Document,
			}
			stringToken := token.NewStringToken(stringPart, stringPart, sourceFileReference)
			expression := ast.NewStringLiteral(stringToken)
			expressions = append(expressions, expression)
		}

		expressionString := stringToken.Text()[start+1 : end-1]
		if len(expressionString) > 0 {
			if strings.HasSuffix(expressionString, "=") {
				expressionString = expressionString[:len(expressionString)-1]
			}
			expression, expressionErr := stringToExpression(expressionString, stringToken.FetchPositionLength())
			if expressionErr != nil {
				return nil, expressionErr
			}
			expressions = append(expressions, expression)
		}

		lastPos = end
	}

	remainingString := stringToken.Text()[lastPos:]

	if len(remainingString) > 0 {
		sourceFileReference := token.SourceFileReference{
			Range: token.MakeRange(
				stringToken.Range.Start(),
				stringToken.Range.End(),
			),
			Document: stringToken.Document,
		}
		stringToken := token.NewStringToken(remainingString, remainingString, sourceFileReference)
		expression := ast.NewStringLiteral(stringToken)
		expressions = append(expressions, expression)
	}

	return expressions, nil
}

/*
func replaceInterpolationStringToString(s string) string {
	return replaceInterpolationString(s, " ++ ", "Debug.toString(%v)")
}

func replaceInterpolationStringToTuple(s string) string {
	return "( " + replaceInterpolationString(s, ", ", "%v") + " )"
}
*/

func stringToExpression(replaced string, sourceFileReference token.SourceFileReference) (ast.Expression, parerr.ParseError) {
	reader := strings.NewReader(replaced)
	localPath, localErr := sourceFileReference.Document.Uri.ToLocalFilePath()
	if localErr != nil {
		panic(localErr)
	}
	runeReader, _ := runestream.NewRuneReader(reader, localPath)

	const exactWhitespace = true
	tokenizer, tokenizerErr := tokenize.NewTokenizerInternalWithStartPosition(runeReader, sourceFileReference.Range.Start(), exactWhitespace)
	if tokenizerErr != nil {
		return nil, tokenizerErr
	}
	parser := NewParser(tokenizer, exactWhitespace)
	expr, exprErr := parser.parseExpressionNormal(0)
	if exprErr != nil {
		return nil, exprErr
	}

	return expr, nil
}

func parseInterpolationStringToTupleExpression(p ParseStream, stringToken token.StringToken) (ast.Expression, parerr.ParseError) {
	expressions, interpolateErr := replaceInterpolationString(stringToken)
	if interpolateErr != nil {
		return nil, interpolateErr
	}

	startParen := token.NewParenToken("(", token.LeftParen, stringToken.FetchPositionLength(), "(")
	endParen := token.NewParenToken(")", token.RightParen, stringToken.FetchPositionLength(), ")")

	tupleLiteral := ast.NewTupleLiteral(startParen, endParen, expressions)

	return ast.NewStringInterpolation(stringToken, tupleLiteral, expressions), nil
}

func makeItString(expression ast.Expression, stringToken token.StringToken) ast.Expression {
	_, wasStringLiteral := expression.(*ast.StringLiteral)
	if wasStringLiteral {
		return expression
	}

	debugModuleName := token.NewTypeSymbolToken("Debug", stringToken.FetchPositionLength(), 0)
	debugModuleNamePart := ast.NewModuleNamePart(ast.NewTypeIdentifier(debugModuleName))
	debugModuleRef := ast.NewModuleReference([]*ast.ModuleNamePart{debugModuleNamePart})
	toStringVar := ast.NewVariableIdentifier(token.NewVariableSymbolToken("toString", stringToken.FetchPositionLength(), 0))
	debugToString := ast.NewQualifiedVariableIdentifierScoped(debugModuleRef, toStringVar)
	return ast.NewFunctionCall(debugToString, []ast.Expression{expression})
}

func parseInterpolationStringToStringExpression(p ParseStream, stringToken token.StringToken) (ast.Expression, parerr.ParseError) {
	expressions, interpolateErr := replaceInterpolationString(stringToken)
	if interpolateErr != nil {
		return nil, interpolateErr
	}

	var lastExpression ast.Expression
	for _, expression := range expressions {
		convertedExpression := makeItString(expression, stringToken)

		if lastExpression != nil {
			appendOperatorToken := token.NewOperatorToken(token.OperatorAppend, stringToken.FetchPositionLength(), "++", "++")
			convertedExpression = ast.NewBinaryOperator(appendOperatorToken, appendOperatorToken, lastExpression, convertedExpression)
		}
		lastExpression = convertedExpression
	}

	return ast.NewStringInterpolation(stringToken, lastExpression, expressions), nil
}

func parseStringLiteral(p ParseStream, stringToken token.StringToken) (ast.Expression, parerr.ParseError) {
	lit := ast.NewStringLiteral(stringToken)

	return lit, nil
}
