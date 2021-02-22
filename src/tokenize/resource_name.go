/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package tokenize

import (
	"github.com/swamp/compiler/src/token"
)

func (t *Tokenizer) parseResourceName(startPos token.PositionToken) (token.Token, error) {
	var a string

	ch := t.nextRune()

	if !isLowerCaseLetter(ch) {
		t.unreadRune()
		return token.ResourceName{}, NewExpectedVariableSymbolError(t.MakeSourceFileReference(startPos), string(ch))
	}

	a += string(ch)

	for {
		ch := t.nextRune()
		if !isLowerCaseLetter(ch) && ch != '/' && !isDigit(ch) && ch != '.' && ch != '_' {
			t.unreadRune()
			break
		}
		a += string(ch)
	}

	return token.NewResourceName(a, t.MakeSourceFileReference(startPos), startPos.Indentation()), nil
}
