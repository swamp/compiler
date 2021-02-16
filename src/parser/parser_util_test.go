/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/swamp/compiler/src/ast"
	parerr "github.com/swamp/compiler/src/parser/errors"
	testutil "github.com/swamp/compiler/src/test"
)

func testParseInternal(code string, errorsAsWarnings bool) (*ast.Program, ParseStream, parerr.ParseError) {
	code = strings.TrimSpace(code)
	const enforceStyle = true
	tokenizer, tokenizerErr := testutil.Setup(code)
	if tokenizerErr != nil {
		ShowError(tokenizer, "test.swamp", tokenizerErr, true, errorsAsWarnings)
		return nil, nil, tokenizerErr
	}
	p := NewParser(tokenizer, enforceStyle)
	program, programErr := p.Parse()
	if programErr != nil {
		ShowError(tokenizer, "test.swamp", programErr, true, errorsAsWarnings)
	}
	return program, p.stream, programErr
}

func testParseExpressionInternal(code string, enforceStyle bool) (*ast.Program, ParseStream, parerr.ParseError) {
	const errorsAsWarnings = false
	code = strings.TrimSpace(code)
	tokenizer, tokenizerErr := testutil.Setup(code)
	if tokenizerErr != nil {
		fmt.Printf("error:%p\n", tokenizerErr)
		ShowError(tokenizer, "test.swamp", tokenizerErr, true, errorsAsWarnings)
		return nil, nil, tokenizerErr
	}

	p := NewParser(tokenizer, enforceStyle)
	program, programErr := p.ParseExpression()
	if programErr != nil {
		fmt.Printf("error:%v\n", programErr)
		ShowError(tokenizer, "test.swamp", programErr, true, errorsAsWarnings)
	}
	return program, p.stream, programErr
}

func testParse(t *testing.T, code string, ast string) {
	const errorsAsWarnings = false
	program, _, programErr := testParseInternal(code, errorsAsWarnings)
	if programErr != nil {
		t.Fatal(programErr)
	}
	ast = strings.TrimSpace(ast)
	generatedAst := strings.TrimSpace(program.String())
	// fmt.Printf("ast:%v\n", generatedAst)
	if ast != generatedAst {
		t.Errorf("wrong ast generated:\n\n%v\nexpected\n\n%v\n", generatedAst, ast)
	}
}

func testParseExpressionHelper(t *testing.T, code string, ast string, enforceStyle bool) {
	program, _, programErr := testParseExpressionInternal(code, enforceStyle)
	if programErr != nil {
		t.Fatal(programErr)
	}
	ast = strings.TrimSpace(ast)
	generatedAst := strings.TrimSpace(program.String())
	if ast != generatedAst {
		t.Errorf("wrong ast generated:\n\n%v\nexpected\n\n%v\n", generatedAst, ast)
	}
}

func testParseExpression(t *testing.T, code string, ast string) {
	testParseExpressionHelper(t, code, ast, true)
}

func testParseExpressionNoStyle(t *testing.T, code string, ast string) {
	testParseExpressionHelper(t, code, ast, false)
}

func testParseError(t *testing.T, code string, expectedErrorType interface{}) {
	const errorsAsWarnings = true
	_, _, programErr := testParseInternal(code, errorsAsWarnings)
	if programErr == nil {
		t.Fatal("it should have failed, but it didn't")
		return
	}

	if expectedErrorType != nil && (reflect.TypeOf(expectedErrorType).Name() != reflect.TypeOf(programErr).Name()) {
		t.Fatalf("expected %T but received %T %v", expectedErrorType, programErr, programErr)
	}
}
