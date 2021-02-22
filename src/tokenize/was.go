/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package tokenize

import "github.com/swamp/compiler/src/token"

func (t *Tokenizer) WasDefaultSymbol() (token.RuneToken, bool) {
	return t.WasRune('_')
}

func (t *Tokenizer) WasRune(requiredRune rune) (token.RuneToken, bool) {
	startPos := t.position
	readRune := t.nextRune()
	wasCorrect := readRune == requiredRune
	if !wasCorrect {
		t.unreadRune()
		return token.RuneToken{}, false
	}
	return token.NewRuneToken(readRune, t.MakeSourceFileReference(startPos)), true
}

func (t *Tokenizer) WasSpacingRune(requiredRune rune) (token.RuneToken, bool) {
	startPos := t.position
	readRune := t.nextRune()
	wasCorrect := readRune == requiredRune
	if !wasCorrect {
		_, detectedComment, _ := t.checkComment(readRune, startPos)
		if detectedComment {
			return t.WasSpacingRune(requiredRune)
		}
		t.unreadRune()
		return token.RuneToken{}, false
	}
	return token.NewRuneToken(readRune, t.MakeSourceFileReference(startPos)), true
}
