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

const interpolationRegExp = `\$\{.*?\}`

func isInterpolation(s string) (bool, error) {
	didMatch, matchErr := regexp.MatchString(interpolationRegExp, s)
	if matchErr != nil {
		return false, matchErr
	}
	return didMatch, nil
}

func parseInterpolationString(s string) [][]int {
	re, reErr := regexp.Compile(interpolationRegExp)
	if reErr != nil {
		panic(reErr)
	}

	locations := re.FindAllStringIndex(s, -1)

	return locations
}

func replaceInterpolationString(s string) string {
	ranges := parseInterpolationString(s)

	result := ""
	lastPos := 0

	for index, item := range ranges {
		start := item[0]
		end := item[1]

		if index > 0 {
			result += " ++ "
		}

		result += fmt.Sprintf("\"%s\"", s[lastPos:start])

		inside := s[start+2 : end-1]

		result += fmt.Sprintf(" ++ Debug.toString(%v)", inside)

		lastPos = end
	}

	remaining := s[lastPos:]

	if len(remaining) > 0 {
		if len(ranges) > 0 {
			result += " ++ "
		}

		result += fmt.Sprintf("\"%s\"", remaining)
	}

	return result
}

func replaceInterpolationStringToExpression(s string) (ast.Expression, parerr.ParseError) {
	replaced := replaceInterpolationString(s)
	reader := strings.NewReader(replaced)
	runeReader, _ := runestream.NewRuneReader(reader, "internal interpolation string")

	const exactWhitespace = true
	tokenizer, tokenizerErr := tokenize.NewTokenizer(runeReader, exactWhitespace)
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

func parseStringLiteral(p ParseStream, stringToken token.StringToken) (ast.Expression, parerr.ParseError) {
	lit := ast.NewStringConstant(stringToken, stringToken.Text())

	wasInterpolation, interpolationErr := isInterpolation(stringToken.Text())
	if interpolationErr != nil {
		return nil, parerr.NewInternalError(p.positionLength(), interpolationErr)
	}

	if wasInterpolation {
		return replaceInterpolationStringToExpression(stringToken.Text())
	}

	return lit, nil
}
