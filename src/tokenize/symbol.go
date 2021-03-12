/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package tokenize

import (
	"github.com/swamp/compiler/src/token"
)

func (t *Tokenizer) parseAnySymbol(startPosition token.PositionToken) (token.Token, TokenError) {
	ch := t.nextRune()
	t.unreadRune()
	if isUpperCaseLetter(ch) {
		typeSymbol, typeSymbolErr := t.ParseTypeSymbol()
		if typeSymbolErr != nil {
			return nil, typeSymbolErr
		}
		booleanSymbol, booleanSymbolErr := DetectUppercaseBoolean(typeSymbol)
		if booleanSymbolErr == nil {
			return booleanSymbol, nil
		}
		return typeSymbol, nil
	}
	variableSymbol, variableSymbolErr := t.ParseVariableSymbol()
	if variableSymbolErr != nil {
		return nil, variableSymbolErr
	}
	if variableSymbol.Name() == "module" {
		t.MaybeOneNewLine()
		t.MaybeOneNewLine()
		t.MaybeOneNewLine()

		return token.NewMultiLineCommentToken("ignore", "ignore", false, t.MakeSourceFileReference(startPosition)), nil
	}
	if variableSymbol.Name() == "not" {
		t.EatOneSpace()
		return token.NewOperatorToken(token.OperatorUnaryNot, variableSymbol.FetchPositionLength(), variableSymbol.Raw(), "NOT"), nil
	}
	keywordSymbol, keywordSymbolErr := DetectLowercaseKeyword(variableSymbol)
	if keywordSymbolErr == nil {
		return keywordSymbol, nil
	}

	return variableSymbol, nil
}
