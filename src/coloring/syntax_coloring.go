/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package coloring

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/swamp/compiler/src/runestream"
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/compiler/src/tokenize"
)

func colorSpace(t token.SpaceToken) string {
	return " "
}

func ColorOperator(t token.OperatorToken) string {
	return color.HiGreenString(t.Raw())
}

func ColorOperatorString(t string) string {
	return color.HiGreenString(t)
}

func ColorVariableSymbol(t token.VariableSymbolToken) string {
	return color.HiMagentaString(t.Raw())
}

func ColorResourceNameSymbol(t token.ResourceName) string {
	return color.HiMagentaString(t.Raw())
}

func ColorTypeIdSymbol(t token.TypeId) string {
	return color.WhiteString(t.Raw())
}

func ColorCharacterToken(t token.CharacterToken) string {
	return color.GreenString(t.Raw())
}

func ColorDefinition(t token.VariableSymbolToken) string {
	return color.HiBlueString(t.Raw())
}

func ColorLocalType(t token.VariableSymbolToken) string {
	return color.GreenString(t.Raw())
}

func ColorTypeSymbol(t token.TypeSymbolToken) string {
	return color.HiMagentaString(t.Raw())
}

func ColorModuleReference(t token.TypeSymbolToken) string {
	return color.MagentaString(t.Raw())
}

func ColorTypeGeneratorName(t token.TypeSymbolToken) string {
	return color.HiCyanString(t.Raw())
}

func ColorPrimitiveType(t token.TypeSymbolToken) string {
	return color.HiBlueString(t.Raw())
}

func ColorAliasNameSymbol(t token.TypeSymbolToken) string {
	return color.CyanString(t.Raw())
}

func ColorNumberLiteral(t token.NumberToken) string {
	return color.GreenString(t.Raw())
}

func ColorRecordField(t token.VariableSymbolToken) string {
	return color.HiYellowString(t.Raw())
}

func colorStringLiteral(t token.StringToken) string {
	return color.HiBlueString(t.Raw())
}

func colorBooleanLiteral(t token.BooleanToken) string {
	return color.BlueString(t.Raw())
}

func colorParen(t token.ParenToken) string {
	return color.HiRedString(t.Raw())
}

func colorSingleLineComment(t token.SingleLineCommentToken) string {
	return color.HiBlackString(t.Raw())
}

func colorMultiLineComment(t token.MultiLineCommentToken) string {
	return color.HiBlackString(t.Raw())
}

func colorSpecialKeyword(t token.Keyword) string {
	return t.Raw()
}

func ColorKeywordString(t string) string {
	return color.HiCyanString(t)
}

func colorKeyword(t token.Keyword) string {
	return color.HiBlackString(t.Raw())
}

func colorExternalFunction(t token.ExternalFunctionToken) string {
	return color.HiBlackString(t.Raw())
}

func colorAsm(t token.AsmToken) string {
	return color.HiBlackString(t.Raw())
}

func colorLambda(t token.LambdaToken) string {
	return color.HiGreenString("\\")
}

func colorNewLine(t token.LineDelimiterToken) string {
	return "\n"
}

func colorToken(t token.Token) string {
	switch v := t.(type) {
	case token.SpaceToken:
		return colorSpace(v)
	case token.OperatorToken:
		return ColorOperator(v)
	case token.BooleanToken:
		return colorBooleanLiteral(v)
	case token.VariableSymbolToken:
		return ColorVariableSymbol(v)
	case token.ResourceName:
		return ColorResourceNameSymbol(v)
	case token.TypeId:
		return ColorTypeIdSymbol(v)
	case token.CharacterToken:
		return ColorCharacterToken(v)
	case token.TypeSymbolToken:
		return ColorTypeSymbol(v)
	case token.NumberToken:
		return ColorNumberLiteral(v)
	case token.StringToken:
		return colorStringLiteral(v)
	case token.ParenToken:
		return colorParen(v)
	case token.SingleLineCommentToken:
		return colorSingleLineComment(v)
	case token.MultiLineCommentToken:
		return colorMultiLineComment(v)
	case token.Keyword:
		return colorKeyword(v)
	case token.LineDelimiterToken:
		return colorNewLine(v)
	case token.ExternalFunctionToken:
		return colorExternalFunction(v)
	case token.AsmToken:
		return colorAsm(v)
	case token.LambdaToken:
		return colorLambda(v)
	}
	panic(fmt.Sprintf("ColorToken: unknown type %T", t))
}

func SyntaxColor(code string) (string, error) {
	ioReader := strings.NewReader(code)
	runeReader, err := runestream.NewRuneReader(ioReader, "for coloring")
	if err != nil {
		return "", err
	}
	const enforceStyle = true
	tokenizer, tokenizerErr := tokenize.NewTokenizerInternal(runeReader, enforceStyle)
	if tokenizerErr != nil {
		return "", tokenizerErr
	}

	a := ""
	for {
		token, tokenErr := tokenizer.ReadTermToken()
		if tokenErr != nil {
			return a, tokenErr
		}
		if _, isEOF := token.(*tokenize.EndOfFile); isEOF {
			return a, nil
		}

		a += colorToken(token)
	}
}
