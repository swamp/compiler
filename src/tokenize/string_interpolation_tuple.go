package tokenize

import "github.com/swamp/compiler/src/token"

func (t *Tokenizer) ParseStringInterpolationTuple(startStringRune rune, startPosition token.PositionToken) (token.StringInterpolationTuple, TokenError) {
	stringToken, stringErr := t.ParseString(startStringRune, startPosition)
	if stringErr != nil {
		return token.StringInterpolationTuple{}, stringErr
	}

	return token.NewStringInterpolationTuple(stringToken), nil
}
