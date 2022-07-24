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

// FixedDecimals TODO: Hack. It is a bit of a hack to set the number of decimals in precision here
const FixedDecimals = 3

type ParseNumberType uint8

const (
	IntType ParseNumberType = iota
	FixedType
	HexType
	BinaryType
)

func (t *Tokenizer) ParseNumber(a string) (token.NumberToken, TokenError) {
	startPosition := t.position
	first := false
	parseType := IntType

	for {
		ch := t.nextRune()
		if parseType == HexType {
			if !isHexDigit(ch) {
				t.unreadRune()
				break
			}
		} else if first && ch == '-' {
			first = false
		} else if len(a) == 1 && ch == 'x' && parseType == IntType {
			parseType = HexType
		} else if len(a) == 1 && ch == 'b' && parseType == IntType {
			parseType = BinaryType
		} else if len(a) >= 1 && ch == '.' {
			if parseType == FixedType {
				return token.NumberToken{}, NewInternalError(fmt.Errorf("two dots in fixed"))
			}
			parseType = FixedType
		} else if !isDigit(ch) {
			t.unreadRune()
			break
		}
		a += string(ch)
	}

	total := a
	if parseType == FixedType {
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
		if total == "" {
			total = "0"
		}
	}

	integerValue, integerValueErr := strconv.ParseInt(total, 0, 64)
	if integerValueErr != nil {
		return token.NumberToken{}, NewUnexpectedEatTokenError(t.MakeSourceFileReference(startPosition), ' ', ' ')
	}
	posLen := t.MakeSourceFileReference(startPosition)
	return token.NewNumberToken(a, int32(integerValue), parseType == FixedType, posLen), nil
}
