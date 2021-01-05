/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package tokenize

import (
	"github.com/swamp/compiler/src/token"
)

func (t *Tokenizer) ParseVariableSymbol() (token.VariableSymbolToken, TokenError) {
	var a string

	startPos := t.position
	ch := t.nextRune()
	if !isLowerCaseLetter(ch) {
		t.unreadRune()
		return token.VariableSymbolToken{}, NewExpectedVariableSymbolError(t.MakePositionLength(startPos), string(ch))
	}

	a += string(ch)

	for {
		ch := t.nextRune()
		if !isSymbol(ch) {
			t.unreadRune()
			break
		}
		a += string(ch)
	}
	return token.NewVariableSymbolToken(a, t.MakePositionLength(startPos), startPos.Indentation()), nil
}
