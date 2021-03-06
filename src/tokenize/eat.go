/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package tokenize

import (
	"github.com/swamp/compiler/src/token"
)

func (t *Tokenizer) EatRune(requiredRune rune) TokenError {
	startPos := t.position
	readRune := t.nextRune()
	if readRune != requiredRune {
		return NewUnexpectedEatTokenError(t.MakeSourceFileReference(startPos), requiredRune, readRune)
	}
	return nil
}

func (t *Tokenizer) checkComment(readRune rune, startPos token.PositionToken) (token.Comment, bool, TokenError) {
	if readRune == '-' {
		nch := t.nextRune()
		if nch == '-' {
			comment := t.ReadSingleLineComment(startPos)
			return comment, true, nil
		} else {
			t.unreadRune()
		}
	} else if readRune == '{' {
		nch := t.nextRune()
		if nch == '-' {
			comment, commentErr := t.ReadMultilineComment(startPos)
			if commentErr != nil {
				return nil, false, commentErr
			}
			return comment, true, nil
		} else {
			t.unreadRune()
		}
	}

	return nil, false, nil
}

func (t *Tokenizer) EatString(requiredString string) TokenError {
	for _, ch := range requiredString {
		eatErr := t.EatRune(ch)
		if eatErr != nil {
			return eatErr
		}
	}
	return nil
}
