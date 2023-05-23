/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/swamp/compiler/src/runestream"
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/compiler/src/tokenize"
)

func Test(t *testing.T) {
	// TODO: Ability.swamp:143 and 156
	/*
	   %"It is not something with charges, so swap places. Pickup: {giverEntity} \
	   and existing in wallet is: {currentAbilityEntity}"
	*/

	str := strings.NewReader(
		`%"It is not something with charges, so swap places. Pickup: {giverEntity} \
        and existing in wallet is: {currentAbilityEntity}"
`,
	)
	r, _ := runestream.NewRuneReader(str, "hello")
	tokenizer, _ := tokenize.NewTokenizerInternal(r, true)

	parsedToken, _ := tokenizer.ReadTermToken()
	fmt.Println(parsedToken.String())
	tuple := parsedToken.(token.StringInterpolationTuple)
	stringToken := tuple.StringToken()
	stringLines := stringToken.StringLines()

	log.Printf("Range for example: %v\n", stringLines)
	log.Printf("totalRange : %v\n", stringToken.Range)

	if len(stringLines) != 2 {
		t.Fail()
	}

	first := stringLines[0]
	second := stringLines[1]

	assert.Equal(t, 0, first.LocalOctetOffset, "Must have correct local offset %v", first.LocalOctetOffset)

	assert.Equal(t, 72, second.LocalOctetOffset, "Must have correct local offset %v", second.LocalOctetOffset)

	totalRange := stringToken.Range

	assert.Equal(t, 0, totalRange.Start().Line(), "Must have correct start line %v", totalRange.Start())
}

func Test2(t *testing.T) {
	// TODO: Ability.swamp:143 and 156
	/*
	   %"It is not something with charges, so swap places. Pickup: {giverEntity} \
	   and existing in wallet is: {currentAbilityEntity}"
	*/

	str := strings.NewReader(
		`%"It has charges with equal id removing {giverEntity} and adding \
                                charges to {currentAbilityEntity}"
`,
	)
	r, _ := runestream.NewRuneReader(str, "hello")
	tokenizer, _ := tokenize.NewTokenizerInternal(r, true)

	parsedToken, _ := tokenizer.ReadTermToken()
	fmt.Println(parsedToken.String())
	tuple := parsedToken.(token.StringInterpolationTuple)
	stringToken := tuple.StringToken()
	stringLines := stringToken.StringLines()

	log.Printf("Range for example: %v\n", stringLines)
	log.Printf("totalRange : %v\n", stringToken.Range)

	if len(stringLines) != 2 {
		t.Fail()
	}

	first := stringLines[0]
	second := stringLines[1]

	assert.Equal(t, 0, first.LocalOctetOffset, "Must have correct local offset %v", first.LocalOctetOffset)

	assert.Equal(t, 63, second.LocalOctetOffset, "Must have correct local offset %v", second.LocalOctetOffset)

	totalRange := stringToken.Range

	assert.Equal(t, 0, totalRange.Start().Line(), "Must have correct start line %v", totalRange.Start())
	assert.Equal(t, 1, totalRange.Start().Column(), "Must have correct start line %v", totalRange.Start())

	assert.Equal(t, 1, totalRange.End().Line(), "Must have correct start line %v", totalRange.Start())
	assert.Equal(t, 65, totalRange.End().Column(), "Must have correct start line %v", totalRange.Start())
}
