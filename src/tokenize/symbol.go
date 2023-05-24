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
		return typeSymbol, nil
	} else {

	}
	variableSymbol, variableSymbolErr := t.ParseVariableSymbol()
	if variableSymbolErr != nil {
		return nil, variableSymbolErr
	}

	booleanSymbol, booleanSymbolErr := DetectLowerCaseBoolean(variableSymbol)
	if booleanSymbolErr == nil {
		return booleanSymbol, nil
	}

	if variableSymbol.Name() == "not" {
		t.EatOneSpace()
		return token.NewOperatorToken(token.OperatorUnaryNot, variableSymbol.FetchPositionLength(),
			variableSymbol.Raw(), "NOT"), nil
	}
	keywordSymbol, keywordSymbolErr := DetectLowercaseKeyword(variableSymbol)
	if keywordSymbolErr == nil {
		return keywordSymbol, nil
	}

	return variableSymbol, nil
}
