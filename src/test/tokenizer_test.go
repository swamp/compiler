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
	//tokenizer.SkipAnyWhitespace()

	_, hopefullySymbolTokenErr := tokenizer.ParseVariableSymbol()
	if hopefullySymbolTokenErr != nil {
		t.Error(hopefullySymbolTokenErr)
	}
}

func expectLambda(t *testing.T, tokenizer *tokenize.Tokenizer) {
	hopefullyLambdaToken, hopefullyLambdaTokenErr := tokenizer.ReadTermToken()
	if hopefullyLambdaTokenErr != nil {
		t.Error(hopefullyLambdaTokenErr)
	}
	_, ok := hopefullyLambdaToken.(token.LambdaToken)
	if !ok {
		t.Errorf("Wrong type. Expected Lambda start but was %v", hopefullyLambdaToken)
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

func expectNumber(t *testing.T, tokenizer *tokenize.Tokenizer, expectedValue float64) {
	hopefullySymbolToken, hopefullySymbolTokenErr := tokenizer.ReadTermToken()
	if hopefullySymbolTokenErr != nil {
		t.Error(hopefullySymbolTokenErr)
	}
	_, ok := hopefullySymbolToken.(token.NumberToken)
	if !ok {
		t.Errorf("Wrong type. Expected Symbol but was %v", hopefullySymbolToken)
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

func TestLambda(t *testing.T) {
	tokenizer, tokenErr := Setup(
		`
\x->22

`)
	if tokenErr != nil {
		t.Fatal(tokenErr)
	}
	expectLambda(t, tokenizer)
	expectVariableSymbol(t, tokenizer, "x")
	expectOperator(t, tokenizer, "->")
	expectNumber(t, tokenizer, 22)
}
