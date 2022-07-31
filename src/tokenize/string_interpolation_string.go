/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package tokenize

import "github.com/swamp/compiler/src/token"

func (t *Tokenizer) ParseStringInterpolationString(stringToken token.StringToken) (token.StringInterpolationString, TokenError) {
	return token.NewStringInterpolationString(stringToken), nil
}
