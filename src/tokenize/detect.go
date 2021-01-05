/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package tokenize

func (t *Tokenizer) DetectRune(requiredRune rune) bool {
	readRune := t.nextRune()
	wasCorrect := readRune == requiredRune
	t.unreadRune()
	return wasCorrect
}

func (t *Tokenizer) DetectString(s string) bool {
	save := t.Tell()
	var correct bool
	for _, expectedRune := range s {
		readRune := t.nextRune()
		correct = readRune == expectedRune
		if !correct {
			break
		}
	}
	t.Seek(save)

	return correct
}
