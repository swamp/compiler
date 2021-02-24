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

func parseExternalFunction(raw string, t *Tokenizer, pos token.SourceFileReference) (token.ExternalFunctionToken, TokenError) {
	comment, afterKeywordSpaceErr := t.EatOneSpace()
	if afterKeywordSpaceErr != nil {
		return token.ExternalFunctionToken{}, afterKeywordSpaceErr
	}
	functionName, err := t.ParseVariableSymbol()
	if err != nil {
		return token.ExternalFunctionToken{}, err
	}
	_, beforeNumberErr := t.EatOneSpace()
	if beforeNumberErr != nil {
		return token.ExternalFunctionToken{}, beforeNumberErr
	}
	numberOfParams, numberOfParamsErr := t.ParseNumber("")
	if numberOfParamsErr != nil {
		return token.ExternalFunctionToken{}, numberOfParamsErr
	}
	return token.NewExternalFunctionToken(raw, functionName.Name(), uint(numberOfParams.Value()), comment, pos), nil
}

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
		return parseExternalFunction(a, t, t.MakeSourceFileReference(pos))
	}

	return nil, NewInternalError(fmt.Errorf("not a starting keyword"))
}
