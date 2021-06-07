/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/ast/codewriter"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	parerr "github.com/swamp/compiler/src/parser/errors"
	"github.com/swamp/compiler/src/runestream"
	"github.com/swamp/compiler/src/tokenize"
	"github.com/swamp/compiler/src/verbosity"
)

func compileToProgram(moduleName string, x string, enforceStyle bool, verbose verbosity.Verbosity) (*tokenize.Tokenizer, *ast.SourceFile, decshared.DecoratedError) {
	ioReader := strings.NewReader(x)
	runeReader, runeReaderErr := runestream.NewRuneReader(ioReader, "style test")
	if runeReaderErr != nil {
		return nil, nil, decorated.NewInternalError(runeReaderErr)
	}
	tokenizer, tokenizerErr := tokenize.NewTokenizer(runeReader, enforceStyle)
	if tokenizerErr != nil {
		const errorsAsWarnings = false
		ShowError(tokenizer, moduleName, tokenizerErr, verbose, errorsAsWarnings)
		return tokenizer, nil, tokenizerErr
	}
	p := NewParser(tokenizer, enforceStyle)

	program, programErr := p.Parse()
	if programErr != nil {
		return tokenizer, nil, programErr
	}

	return tokenizer, program, nil
}

func testStyle(code string, useCores bool) (string, string, string, error) {
	fmt.Println("==== Unformatted ====")
	code = strings.TrimSpace(code)
	fmt.Println(code)
	fmt.Println("========")
	const verbose = verbosity.Mid

	const enforceStyle = true
	const doNotEnforceStyle = false
	// Intentionally enforce style and make sure that it fails
	_, _, shouldFailErr := compileToProgram("Main", code, enforceStyle, verbose)
	if shouldFailErr == nil {
		return "", "", "", fmt.Errorf("it is supposed to fail with unformatted text")
	}

	_, program, compileErr := compileToProgram("Main", code, doNotEnforceStyle, verbose)
	if compileErr != nil {
		return "", "", "", compileErr
	}

	withColor, colorErr := codewriter.WriteCode(program, true)
	if colorErr != nil {
		return "", "", "", colorErr
	}
	codeWithoutColor, withoutColorErr := codewriter.WriteCode(program, false)
	if withoutColorErr != nil {
		return "", "", "", withoutColorErr
	}

	fmt.Println("=== reformatted ===")
	fmt.Println(codeWithoutColor)
	fmt.Println("=== reformatted (with color) ===")
	fmt.Println(withColor)
	fmt.Println("======")

	_, reProgram, reCompileErr := compileToProgram("Main", codeWithoutColor, enforceStyle, verbose)
	if reCompileErr != nil {
		return "", "", "", fmt.Errorf("recompile error:%v", reCompileErr)
	}

	recodeWithoutColor, recodeWithoutColorErr := codewriter.WriteCode(reProgram, false)
	if recodeWithoutColorErr != nil {
		return "", "", "", recodeWithoutColorErr
	}

	fmt.Println("==========")
	return withColor, codeWithoutColor, recodeWithoutColor, nil
}

func testStyleInternal(t *testing.T, code string, expected string) {
	const useCores = false
	resultWithColor, result, recodedResult, err := testStyle(code, useCores)
	if err != nil {
		t.Error(err)
	}

	if recodedResult != result {
		fmt.Printf("recompile mismatch expected:\n%v\n\nBut received:\n%v\n", result, recodedResult)
		t.Errorf("recompile mismatch:\n'%v'\n\n'%v'", result, recodedResult)
	}

	expected = strings.TrimSpace(expected)
	if result != expected {
		fmt.Printf("%v\n%v\n", resultWithColor, expected)
		t.Errorf("mismatch:\n'%v'\n\nexpected:\n\n'%v'", result, expected)
	}
}

func testStyleInternalErr(t *testing.T, code string, expectedError error) {
	const useCores = false
	_, _, _, testErr := testStyle(code, useCores)
	if testErr == nil {
		t.Error("test was supposed to fail, but didn't")
		return
	}

	if errors.Is(testErr, expectedError) {
		t.Errorf("unexpected fail: %v %T but expected %T", testErr, testErr, expectedError)
	}
}

func TestFormatRecordTypeAlias(t *testing.T) {
	testStyleInternal(t, `
type alias Sprite = {  y    :   Int  ,
 x : Int }`,
		`
type alias Sprite =
    { y : Int
    , x : Int
    }
`)
}

func TestFormatBadContinuation(t *testing.T) {
	testStyleInternalErr(t, `
type alias Sprite = {  x    :   Int  ,
y : Int }`,
		parerr.ExpectedContinuationLineOrOneSpace{})
}

func TestFormatCase(t *testing.T) {
	testStyleInternal(t, `
f a b =
    case    x   of
        Something ->
            3
        Other ->
            2
`,
		`
f a b =
    case x of
        Something ->
            3

        Other ->
            2
`)
}

func TestFormatLet(t *testing.T) {
	testStyleInternal(t, `
f a b =
    let
         x =
             13

         y =
             32
    in
           someFunc x * 3
`,
		`
f a b =
    let
        x =
            13

        y =
            32
    in
    someFunc x * 3
`)
}

func TestFormatLet2(t *testing.T) {
	testStyleInternal(t, `
f a b =
    let
        x = callme (4 + 4)
         y = 3
        zarg = 8 * 99
    in
    x+y`, `

f a b =
    let
        x =
            callme (4 + 4)

        y =
            3

        zarg =
            8 * 99
    in
    x + y

`)
}

func TestFormatIf(t *testing.T) {
	testStyleInternal(t, `
someFunction testing c =
    if c > 39 then
        testing
    else
        c+99
`,
		`
someFunction testing c =
    if c > 39 then
        testing
    else
        c + 99
`)
}

func TestFormatImport(t *testing.T) {
	testStyleInternal(t, `
import AnotherFile.AndSubFolder

import YetAnother

sample:Int     ->
  Int
sample a =
    AnotherFile.MakeSomething -10
`, `
import AnotherFile.AndSubFolder
import YetAnother


sample : Int -> Int
sample a =
    AnotherFile.MakeSomething -10`)
}

func TestFormatMultipleDefinitions(t *testing.T) {
	testStyleInternal(t, `
import AnotherFile.AndSubFolder

sample:Int     ->
  Int
sample a =
    AnotherFile.makeSomething -10

something:    Int->  String

something a =
      "hello  "
`, `
import AnotherFile.AndSubFolder


sample : Int -> Int
sample a =
    AnotherFile.makeSomething -10


something : Int -> String
something a =
    "hello  "
`)
}

func TestFormatCustomType(t *testing.T) {
	testStyleInternal(t, `
type MyMaybe a =
  Just a
       | Nothing

`, `
type MyMaybe a =
    Just a
    | Nothing
`)
}

func TestFormatUnary(t *testing.T) {
	testStyleInternal(t, `
test a =
    if not   True then
        "hello"
    else
            "asd"

`, `
test a =
    if not True then
        "hello"
    else
        "asd"
`)
}

func TestFormatTypeRef(t *testing.T) {
	testStyleInternal(t, `
doIt:Int    ->  List
    Sprite
doIt    a      =
       b
`, `
doIt : Int -> List Sprite
doIt a =
    b
`)
}

func TestFormatAlias(t *testing.T) {
	testStyleInternal(t, `
type alias Tinkering a t =
   { solder : Bool
   , used : a
   , something : t
   }
`, `
type alias Tinkering a t =
    { solder : Bool
    , used : a
    , something : t
    }`)
}

func TestFormatAnnotation(t *testing.T) {
	testStyleInternal(t, `
annotation :List a -> List b
	`, `
annotation : List a -> List b
`)
}

func TestFormatAnnotationFunc(t *testing.T) {
	testStyleInternal(t, `
annotation :(a -> b)-> List a -> List b
	`, `
annotation : (a -> b) -> List a -> List b
`)
}

func TestFormatAnnotationFuncMulti(t *testing.T) {
	testStyleInternal(t, `
annotation:(a ->  b ->  b) ->
  b
	`, `
annotation : (a -> b -> b) -> b
`)
}

func TestFormatCustomTypeAgain(t *testing.T) {
	testStyleInternal(t, `
type Direction =
    NotMoving
    | Right
    |  Down
    |  Left
     | Up
	`, `
type Direction =
    NotMoving
    | Right
    | Down
    | Left
    | Up
`)
}

func TestFormatCustomTypeAgain2(t *testing.T) {
	testStyleInternal(t, `
type   Status =
    Unknown
     | Something Int
	`, `
type Status =
    Unknown
    | Something Int
`)
}

func TestFormatOperator(t *testing.T) {
	testStyleInternal(t, `
a : Int -> Int
a x =
    let
       x = callme (4 + 5)
    in
    let
             bogus =4
    in
      x+y*bogus`, `
a : Int -> Int
a x =
    let
        x =
            callme (4 + 5)
    in
    let
        bogus =
            4
    in
    x + y * bogus
`)
}

func TestFormatPipeRight(t *testing.T) {
	testStyleInternal(t,
		`
tester:String -> Bool
tester b =
    first (2 + 2)
        |>     second b
        |> third
`, `
tester : String -> Bool
tester b =
    first (2 + 2)
        |> second b
        |> third
`)
}

func TestFormatCustomType3(t *testing.T) {
	testStyleInternal(t,
		`
type Maybe    a =
    Nothing
     | Just a


a : Int -> { x : Maybe Int }

`, `
type Maybe a =
    Nothing
    | Just a


a : Int -> { x : Maybe Int
}
`)
}
