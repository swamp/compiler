/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package tokenize

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

func DetectUppercaseBoolean(t token.TypeSymbolToken) (token.BooleanToken, error) {
	switch t.Name() {
	case "True":
		return token.NewBooleanToken(t.Raw(), true, t.FetchPositionLength()), nil
	case "False":
		return token.NewBooleanToken(t.Raw(), false, t.FetchPositionLength()), nil
	}

	return token.BooleanToken{}, fmt.Errorf("not found")
}
