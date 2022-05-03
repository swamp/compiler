package tokenize

import "github.com/swamp/compiler/src/token"

func (t *Tokenizer) ParseStringInterpolationString(startStringRune rune, startPosition token.PositionToken) (token.StringInterpolationString, TokenError) {
	stringToken, stringErr := t.ParseString(startStringRune, startPosition)
	if stringErr != nil {
		return token.StringInterpolationString{}, stringErr
	}

	return token.NewStringInterpolationString(stringToken), nil
}
