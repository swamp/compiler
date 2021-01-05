/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package tokenize

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/swamp/compiler/src/token"
)

// TODO: Hack. It is a bit of a hack to set the number of decimals in precision here
const FixedDecimals = 3

func (t *Tokenizer) ParseNumber(a string) (token.NumberToken, TokenError) {
	startPosition := t.position
	first := false
	isHex := false
	isFixed := false
	for {
		ch := t.nextRune()
		if isHex {
			if !isHexDigit(ch) {
				t.unreadRune()
				break
			}
		} else if first && ch == '-' {
			first = false
		} else if len(a) == 1 && ch == 'x' {
			isHex = true
		} else if len(a) >= 1 && ch == '.' {
			if isFixed {
				return token.NumberToken{}, NewInternalError(fmt.Errorf("two dots in fixed"))
			}
			isFixed = true
		} else if !isDigit(ch) {
			t.unreadRune()
			break
		}
		a += string(ch)
	}

	total := a
	if isFixed {
		parts := strings.Split(a, ".")
		if len(parts) > 2 {
			return token.NumberToken{}, NewInternalError(fmt.Errorf("illegal fixed value"))
		}

		decimals := parts[1]
		numberOfDecimals := len(decimals)
		padCount := FixedDecimals - numberOfDecimals
		if padCount < 0 {
			return token.NumberToken{}, NewInternalError(fmt.Errorf("illegal fixed value"))
		}
		decimals += strings.Repeat("0", padCount)

		integerPart := parts[0]
		integerPart = strings.TrimLeft(integerPart, "0")

		if integerPart == "" {
			decimals = strings.TrimLeft(decimals, "0")
		}
		total = integerPart + decimals

	}

	integerValue, integerValueErr := strconv.ParseInt(total, 0, 32)
	if integerValueErr != nil {
		return token.NumberToken{}, NewUnexpectedEatTokenError(t.MakePositionLength(startPosition), ' ', ' ')
	}
	posLen := t.MakePositionLength(startPosition)
	return token.NewNumberToken(a, int32(integerValue), isFixed, posLen), nil
}
