package tokenize

import "github.com/swamp/compiler/src/token"

func (t *Tokenizer) ParseStringInterpolationString(stringToken token.StringToken) (token.StringInterpolationString, TokenError) {
	return token.NewStringInterpolationString(stringToken), nil
}
