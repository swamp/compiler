package tokenize

import "github.com/swamp/compiler/src/token"

func (t *Tokenizer) ParseStringInterpolationTuple(stringToken token.StringToken) (token.StringInterpolationTuple, TokenError) {
	return token.NewStringInterpolationTuple(stringToken), nil
}
