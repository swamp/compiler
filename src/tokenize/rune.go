/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package tokenize

import (
	"strings"
)

func isIndentation(ch rune) bool {
	return ch == ' '
}

func isNewLine(ch rune) bool {
	return ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isUpperCaseLetter(ch rune) bool {
	return (ch >= 'A' && ch <= 'Z')
}

func isLowerCaseLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z')
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

func isHexDigit(ch rune) bool {
	return isDigit(ch) || (ch >= 'A' && ch <= 'F') //(ch >= 'a' && ch <= 'f')
}

func isSymbol(ch rune) bool {
	return isLetter(ch) || isDigit(ch) || ch == '?'
}

func isStartString(ch rune) bool {
	return ch == '\'' || ch == '"'
}

func isOperator(ch rune) bool {
	return strings.Contains(":|<=>!-+*/.&^~", string(ch))
}
