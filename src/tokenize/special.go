/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package tokenize

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

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
	case "externalfn":
		return token.NewExternalFunctionToken(t.MakeSourceFileReference(pos)), nil
	case "externalvarfn":
		return token.NewExternalVarFunction(t.MakeSourceFileReference(pos)), nil
	case "externalvarexfn":
		return token.NewExternalVarExFunction(t.MakeSourceFileReference(pos)), nil
	}

	return nil, NewInternalError(fmt.Errorf("not a starting keyword '%s'", a))
}
