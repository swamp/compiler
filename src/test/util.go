/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package testutil

import (
	"strings"

	"github.com/swamp/compiler/src/runestream"
	"github.com/swamp/compiler/src/tokenize"
)

func Setup(x string) (*tokenize.Tokenizer, tokenize.TokenError) {
	x = strings.TrimSpace(x)
	ioReader := strings.NewReader(x)
	runeReader, runeReaderErr := runestream.NewRuneReader(ioReader, "for test")
	if runeReaderErr != nil {
		return nil, tokenize.NewInternalError(runeReaderErr)
	}
	const enforceStyle = true
	tokenizer, tokenizerErr := tokenize.NewTokenizerInternal(runeReader, enforceStyle)
	if tokenizerErr != nil {
		return tokenizer, tokenizerErr
	}
	return tokenizer, nil
}
