/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package tokenize

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

func DetectLowercaseBooleanLiteral(t token.VariableSymbolToken) (token.BooleanToken, error) {
	switch t.Name() {
	case "true":
		return token.NewBooleanToken(t.Raw(), true, t.SourceFileReference), nil
	case "false":
		return token.NewBooleanToken(t.Raw(), false, t.SourceFileReference), nil
	}

	return token.BooleanToken{}, fmt.Errorf("not found")
}
