/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package tokenize

import (
	"fmt"
	"strings"

	"github.com/swamp/compiler/src/token"
)

func parseAsm(t *Tokenizer, pos token.SourceFileReference) (token.AsmToken, TokenError) {
	asm := t.ReadStringUntilEndOfLine()
	asm = strings.TrimSpace(asm)
	return token.NewAsmToken(asm, pos), nil
}

func (t *Tokenizer) ParseSpecialKeyword(pos token.PositionToken) (token.Token, TokenError) {
	var a string

	for {
		ch := t.nextRune()
		if !isSymbol(ch) {
			t.unreadRune()
			break
		}
		a += string(ch)
	}
	switch a {
	case "asm":
		return parseAsm(t, t.MakeSourceFileReference(pos))
	case "externalfn":
		return token.NewExternalFunctionToken(t.MakeSourceFileReference(pos)), nil
	}

	return nil, NewInternalError(fmt.Errorf("not a starting keyword"))
}
