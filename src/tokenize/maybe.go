/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package tokenize

func (t *Tokenizer) MaybeRune(requiredRune rune) bool {
	_, wasRune := t.WasRune(requiredRune)
	return wasRune
}

func (t *Tokenizer) MaybeSpacingRune(requiredRune rune) bool {
	_, wasRune := t.WasSpacingRune(requiredRune)
	return wasRune
}

func (t *Tokenizer) MaybeString(requiredString string) bool {
	save := t.Tell()
	for _, ch := range requiredString {
		wasRune := t.MaybeRune(ch)
		if !wasRune {
			t.Seek(save)
			return false
		}
	}
	return true
}

func (t *Tokenizer) MaybeEOF() bool {
	isEOF := t.MaybeRune(0)
	if isEOF {
		t.unreadRune()
	}
	return isEOF
}

func (t *Tokenizer) MaybeOneNewLine() bool {
	return t.MaybeSpacingRune('\n') || t.MaybeEOF()
}

func (t *Tokenizer) MaybeAssign() bool {
	return t.MaybeRune('=')
}

func (t *Tokenizer) MaybeAccessor() bool {
	return t.MaybeRune('.')
}

func (t *Tokenizer) MaybeColon() bool {
	return t.MaybeRune(':')
}
