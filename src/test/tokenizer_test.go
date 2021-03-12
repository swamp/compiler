/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package testutil

import (
	"testing"

	"github.com/swamp/compiler/src/token"
	"github.com/swamp/compiler/src/tokenize"
)

func expectTypeSymbol(t *testing.T, tokenizer *tokenize.Tokenizer, expectedString string) {
	hopefullySymbolToken, hopefullySymbolTokenErr := tokenizer.ReadTermToken()
	if hopefullySymbolTokenErr != nil {
		t.Error(hopefullySymbolTokenErr)
	}
	_, ok := hopefullySymbolToken.(token.TypeSymbolToken)
	if !ok {
		t.Errorf("Wrong type. Expected TypeSymbol but was %v", hopefullySymbolToken)
	}
}

func expectVariableSymbol(t *testing.T, tokenizer *tokenize.Tokenizer, expectedString string) {
	// tokenizer.SkipAnyWhitespace()

	_, hopefullySymbolTokenErr := tokenizer.ParseVariableSymbol()
	if hopefullySymbolTokenErr != nil {
		t.Error(hopefullySymbolTokenErr)
	}
}

func expectOperator(t *testing.T, tokenizer *tokenize.Tokenizer, expectedString string) {
	hopefullyOperatorToken, hopefullyOperatorTokenErr := tokenizer.ReadTermToken()
	if hopefullyOperatorTokenErr != nil {
		t.Error(hopefullyOperatorTokenErr)
	}
	_, ok := hopefullyOperatorToken.(token.OperatorToken)
	if !ok {
		t.Errorf("Wrong type. Expected operator but was %v", hopefullyOperatorToken)
	}
}

func TestType(t *testing.T) {
	tokenizer, tokenErr := Setup(
		`
hello:Int->String

`)
	if tokenErr != nil {
		t.Fatal(tokenErr)
	}
	expectVariableSymbol(t, tokenizer, "hello")
	expectOperator(t, tokenizer, ":")
	expectTypeSymbol(t, tokenizer, "Int")
	expectOperator(t, tokenizer, "->")
	expectTypeSymbol(t, tokenizer, "String")
}
