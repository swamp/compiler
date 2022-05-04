/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"fmt"
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

func replaceInterpolationString(s string, delimiter string, encapsulate string) string {
	ranges := parseInterpolationString(s)

	result := ""
	lastPos := 0

	for _, item := range ranges {
		start := item[0]
		end := item[1]

		first := s[lastPos:start]
		if len(first) > 0 {
			if len(result) > 0 {
				result += delimiter
			}

			result += fmt.Sprintf("\"%s\"", first)
		}

		expressionString := s[start+1 : end-1]
		if len(expressionString) > 0 {
			if len(result) > 0 {
				result += delimiter
			}
			if strings.HasSuffix(expressionString, "=") {
				expressionString = expressionString[:len(expressionString)-1]
				result += fmt.Sprintf("\"%v=\"", expressionString)
				result += delimiter
			}
			result += fmt.Sprintf(encapsulate, expressionString)
		}

		lastPos = end
	}

	remaining := s[lastPos:]

	if len(remaining) > 0 {
		if len(ranges) > 0 {
			result += delimiter
		}

		result += fmt.Sprintf("\"%s\"", remaining)
	}

	return result
}

func replaceInterpolationStringToString(s string) string {
	return replaceInterpolationString(s, " ++ ", "Debug.toString(%v)")
}

func replaceInterpolationStringToTuple(s string) string {
	return "( " + replaceInterpolationString(s, ", ", "%v") + " )"
}

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
	replaced := replaceInterpolationStringToTuple(stringToken.Text())
	expr, exprErr := stringToExpression(replaced, stringToken.FetchPositionLength())
	if exprErr != nil {
		return nil, exprErr
	}

	return expr, nil
}

func parseInterpolationStringToStringExpression(p ParseStream, stringToken token.StringToken) (ast.Expression, parerr.ParseError) {
	replaced := replaceInterpolationStringToString(stringToken.Text())
	expr, exprErr := stringToExpression(replaced, stringToken.FetchPositionLength())
	if exprErr != nil {
		return nil, exprErr
	}

	return expr, nil
}

func parseStringLiteral(p ParseStream, stringToken token.StringToken) (ast.Expression, parerr.ParseError) {
	lit := ast.NewStringConstant(stringToken, stringToken.Text())

	return lit, nil
}
