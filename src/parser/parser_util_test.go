/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/swamp/compiler/src/ast"
	parerr "github.com/swamp/compiler/src/parser/errors"
	testutil "github.com/swamp/compiler/src/test"
	"github.com/swamp/compiler/src/verbosity"
)

func testParseInternal(code string, errorsAsWarnings ReportAsSeverity) (*ast.SourceFile, ParseStream, parerr.ParseError) {
	code = strings.TrimSpace(code)
	const enforceStyle = true
	tokenizer, tokenizerErr := testutil.Setup(code)
	if tokenizerErr != nil {
		ShowError(tokenizer, "test.swamp", tokenizerErr, verbosity.Mid, errorsAsWarnings)
		return nil, nil, tokenizerErr
	}
	p := NewParser(tokenizer, enforceStyle)
	program, programErr := p.Parse()
	if programErr != nil {
		ShowWarningOrError(tokenizer, programErr)
	}
	return program, p.stream, programErr
}

func testParseExpressionInternal(code string, enforceStyle bool) (*ast.SourceFile, ParseStream, parerr.ParseError) {
	const errorsAsWarnings = ReportAsSeverityNote
	code = strings.TrimSpace(code)
	tokenizer, tokenizerErr := testutil.Setup(code)
	if tokenizerErr != nil {
		ShowError(tokenizer, "test.swamp", tokenizerErr, verbosity.Mid, errorsAsWarnings)
		return nil, nil, tokenizerErr
	}

	p := NewParser(tokenizer, enforceStyle)
	program, programErr := p.ParseExpression()
	if programErr != nil {
		ShowError(tokenizer, "test.swamp", programErr, verbosity.Mid, errorsAsWarnings)
	}
	return program, p.stream, programErr
}

func testParse(t *testing.T, code string, ast string) {
	const errorsAsWarnings = ReportAsSeverityNote
	program, _, programErr := testParseInternal(code, errorsAsWarnings)
	if programErr != nil {
		t.Fatal(programErr)
	}
	ast = strings.TrimSpace(ast)
	generatedAst := strings.TrimSpace(program.String())
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

func testParseError(t *testing.T, code string, expectedErrorType interface{}) {
	const errorsAsWarnings = ReportAsSeverityNote

	_, _, programErr := testParseInternal(code, errorsAsWarnings)
	if programErr == nil {
		t.Fatal("it should have failed, but it didn't")
		return
	}

	if expectedErrorType != nil && (reflect.TypeOf(expectedErrorType).Name() != reflect.TypeOf(programErr).Name()) {
		t.Fatalf("expected %T but received %T %v", expectedErrorType, programErr, programErr)
	}
}
