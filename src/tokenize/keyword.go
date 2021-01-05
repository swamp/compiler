/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package tokenize

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

func DetectLowercaseKeyword(t token.VariableSymbolToken) (token.Token, error) {
	pos := t.FetchPositionLength()
	raw := t.Name()
	switch raw {
	case "if":
		return token.NewKeyword(raw, token.If, pos), nil
	case "then":
		return token.NewKeyword(raw, token.Then, pos), nil
	case "else":
		return token.NewKeyword(raw, token.Else, pos), nil
	case "let":
		return token.NewKeyword(raw, token.Let, pos), nil
	case "in":
		return token.NewKeyword(raw, token.In, pos), nil
	case "case":
		return token.NewKeyword(raw, token.Case, pos), nil
	case "type":
		return token.NewKeyword(raw, token.TypeDef, pos), nil
	case "of":
		return token.NewKeyword(raw, token.Of, pos), nil
	case "alias":
		return token.NewKeyword(raw, token.Alias, pos), nil
	case "import":
		return token.NewKeyword(raw, token.Import, pos), nil
	}

	return token.Keyword{}, fmt.Errorf("unknown keyword")
}
