/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"testing"

	"github.com/swamp/compiler/src/ast"
	parerr "github.com/swamp/compiler/src/parser/errors"
)

func TestNumber(t *testing.T) {
	testParseExpression(t,
		`2`,
		`#2`)
}

func TestNewFunction(t *testing.T) {
	testParse(t,
		`
first (startNumber: Int, somethingElse: (String -> List a)) -> (List (List Fixed), Int) =
    startNumber * startNumber
`,
		`
[FnDef $first = [Fn [TypeParamContext [a] +> [FnType [TypeReference $Int], [FnType [TypeReference $String] -> [TypeReference $List [[LocalTypeNameRef [LocalTypeNameDef [LocalTypeName a]]]]]] -> [TupleType [[TypeReference $List [[TypeReference $List [[TypeReference $Fixed]]]]] [TypeReference $Int]]]]] (startNumber, somethingElse) = ($startNumber * $startNumber)]]
`)
}

func TestExternalVar(t *testing.T) {
	testParse(t,
		`
__externalvarfn head : (List a) -> Maybe a
`,
		`
[FnDef $head = [Fn [TypeParamContext [a] +> [FnType [TypeReference $List [[LocalTypeNameRef [LocalTypeNameDef [LocalTypeName a]]]]] -> [TypeReference $Maybe [[LocalTypeNameRef [LocalTypeNameDef [LocalTypeName a]]]]]]] (_) = [FnDeclExpr 2]]]
`)
}

func TestFixedNumber(t *testing.T) {
	testParseExpression(t,
		`2.03`,
		`#!2030`)
}

func TestFixedNumberSmall(t *testing.T) {
	testParseExpression(t,
		`0.4`,
		`#!400`)
}

func xTestMultilineCommentOnOneLine(t *testing.T) {
	testParseExpression(t,
		`{- this is basically a comment -}`,
		`[multilinecomment ' this is basically a comment ']`)
}

func xTestMultilineComment(t *testing.T) {
	testParseExpression(t,
		`{-  
  this is basically 
              a comment -}
`,
		`[multilinecomment '  
  this is basically 
              a comment ']`)
}

func TestChar(t *testing.T) {
	testParseExpression(t,
		`'c'
`,
		`'c'`)
}

func TestConstant(t *testing.T) {
	testParse(t,
		`
fn : String =
    "Hello"
`,
		`
[Constant $fn = 'Hello']
`)
}

func TestStringInterpolation(t *testing.T) {
	testParseExpression(t,
		`
let
    a = "${a} ${b}"
in
a
`,

		`
[Let: [[LetAssign [$a] = '${a} ${b}']] in $a]
`)
}

func TestGuard(t *testing.T) {
	testParseExpression(t,
		`
let
    a =
        | y > 3 -> y
        | y == 4 -> 45
        | _ -> "hello"
in
a
`,

		`
[Let: [[LetAssign [$a] = [Guard [[($y > #3) => $y] [($y == #4) => #45]] [_ => 'hello']]]] in $a]
`)
}

func TestGuardExpression(t *testing.T) {
	testParse(t,
		`
fn : (x: Int) -> Int =
    | x < 2 -> 42
    | x == 4 -> 45
    | _ -> -1
`,

		`
[FnDef $fn = [Fn [FnType [TypeReference $Int] -> [TypeReference $Int]] (x) = [Guard [[($x < #2) => #42] [($x == #4) => #45]] [_ => #-1]]]]
`)
}

func TestWildcardType(t *testing.T) {
	testParse(t,
		`
fn : * -> Int
`,

		`
[FnDef $fn = [Fn [FnType [FnType [AnyMatchingType *] -> [TypeReference $Int]]] () = [FnDeclExpr 0]]]
`)
}

func TestResourceName(t *testing.T) {
	testParseExpression(t,
		`
@some/name/here2.png
`,
		`@some/name/here2.png`)
}

func xTestSingleLineComment(t *testing.T) {
	testParseExpression(t,
		`--   another comment
`,
		` [singlelinecomment '   another comment']`)
}

func TestSingleLineCommentInLet(t *testing.T) {
	testParseExpression(t,
		`
let
    -- this is just a test to check
    x = 3
in
4
`,
		`[Let: [[LetAssign [$x] = #3]] in #4]`)
}

func TestType(t *testing.T) {
	testParseExpression(t,
		`Unknown`,
		`[CCall [TypeReference $Unknown]]`)
}

func TestString(t *testing.T) {
	testParseExpression(t,
		`"hello, world!"`,
		`'hello, world!'`)
}

func TestStringContinuation(t *testing.T) {
	testParseExpression(t,
		`"hello, world! \
  .next line"`,
		`'hello, world! .next line'`)
}

func TestStringTriple(t *testing.T) {
	testParseExpression(t,
		`
"""x
hello, world!
   .next line"""`,
		`'x
hello, world!
   .next line'`)
}

func TestStringContinuationWithIndentation(t *testing.T) {
	testParseExpression(t,
		`"hello, world!  \
               .next line"`,
		`'hello, world!  .next line'`)
}

func TestBool(t *testing.T) {
	testParseExpression(t,
		`true`,
		`€true`)
}

func TestEmptyList(t *testing.T) {
	testParseExpression(t,
		`[]`,
		`[ListLiteral []]`)
}

func TestIf(t *testing.T) {
	testParseExpression(t,
		`if true then 1 else 0`,
		`[If €true then #1 else #0]`)
}

func TestIfWithNewLine(t *testing.T) {
	testParseExpression(t,
		`
if true then
    1
else
    0
`,
		`
[If €true then #1 else #0]
`)
}

func TestIfWithNewLineIf(t *testing.T) {
	testParseExpression(t,
		`
if true then
    if false then
        22
    else
        44
else
    0
`,
		`
[If €true then [If €false then #22 else #44] else #0]
`)
}

func TestIfWithNewLineAndLet(t *testing.T) {
	testParseExpression(t,
		`
if true then
    let
        x = 44
    in
    y
else
    0
`,
		`
[If €true then [Let: [[LetAssign [$x] = #44]] in $y] else #0]
`)
}

func TestIfWithNewLineAndLetMultiAssign(t *testing.T) {
	testParseExpression(t,
		`
if true then
    let
        x = 44

        y = 98
    in
    y
else
    0
`,
		`
[If €true then [Let: [[LetAssign [$x] = #44] [LetAssign [$y] = #98]] in $y] else #0]
`)
}

func TestIfWithNewLineAndLetMultiline(t *testing.T) {
	testParseExpression(t,
		`
if true then
    let
        x =
            case something of
                Kalle -> 2

                Lisa a ->
                    44

        y = 98
    in
    y
else
    0
`,
		`
[If €true then [Let: [[LetAssign [$x] = [CaseCustomType $something of [CaseConsCustomType $Kalle ([]) => #2];[CaseConsCustomType $Lisa ([$a]) => #44]]] [LetAssign [$y] = #98]] in $y] else #0]
`)
}

func TestSimpleAdd(t *testing.T) {
	testParseExpression(t,
		`
2 + 3
`,
		`(#2 + #3)`)
}

func TestString2(t *testing.T) {
	testParseExpression(t,
		`
"Hello"
`,
		`'Hello'`)
}

func TestArithmetic(t *testing.T) {
	testParseExpression(t,
		`
2 + 3 * 4 + 5
`,
		`((#2 + (#3 * #4)) + #5)`)
}

func TestRecordLiteral(t *testing.T) {
	testParseExpression(t,
		`
{ a = 3 }
`,
		`[RecordLiteral [[$a = #3]]]`)
}

func TestRecordLiteralTwoItems(t *testing.T) {
	testParseExpression(t,
		`
{ a = 3, b = 4 }
`,
		`[RecordLiteral [[$a = #3] [$b = #4]]]`)
}

func TestRecordLiteralSeveralLines(t *testing.T) {
	testParseExpression(t,
		`
{ a = 3
, b = 4 }
`,
		`[RecordLiteral [[$a = #3] [$b = #4]]]`)
}

func TestListOneLineFromList(t *testing.T) {
	testParseExpression(t,
		`
Array.fromList [ Array.fromList [ 0, 1, 2, 3 ] Array.fromList [ 8, 9, 10, 11 ] ]
`,
		`[Call Array.$fromList [[ListLiteral [[Call Array.$fromList [[ListLiteral [#0 #1 #2 #3]] Array.$fromList [ListLiteral [#8 #9 #10 #11]]]]]]]]`)
}

func TestContinueOnNextLine(t *testing.T) {
	testParseExpression(t,
		`
callSomething a
    |> nextLine
`,
		`([Call $callSomething [$a]] |> $nextLine)`)
}

func TestContinueOnNextLine2(t *testing.T) {
	testParseExpression(t,
		`
callSomething a
    nextLine
`,
		`[Call $callSomething [$a $nextLine]]`)
}

func TestRecordLiteralSeveralLinesHex(t *testing.T) {
	testParseExpression(t,
		`
{ a = 0x00FF00FF
, b = 4 }
`,
		`[RecordLiteral [[$a = #16711935] [$b = #4]]]`)
}

func TestRecordLiteralSeveralLines2(t *testing.T) {
	testParseExpression(t,
		`
let
    x =
        { a = { scaleX = 100, scaleY = 200 }
        , b = 4
        }
in
x
`,
		`[Let: [[LetAssign [$x] = [RecordLiteral [[$a = [RecordLiteral [[$scaleX = #100] [$scaleY = #200]]]] [$b = #4]]]]] in $x]`)
}

func TestListLiteralTwoItems(t *testing.T) {
	testParseExpression(t,
		`
[ 3, 4 + 5 ]
`,
		`[ListLiteral [#3 (#4 + #5)]]`)
}

func TestArrayLiteralTwoItems(t *testing.T) {
	testParseExpression(t,
		`
[| 3, 4 + 5 |]
`,
		`[ArrayLiteral [#3 (#4 + #5)]]`)
}

func TestSimpleDefinition(t *testing.T) {
	testParse(t,
		`
a =
    3
`,
		`[Constant $a = #3]`)
}

func TestCurrying(t *testing.T) {
	testParse(t,
		`
f (name: String, score: Int) -> Bool =
    if name == "Peter" then
        score * 2
    else
        score


another : (score: Int) -> Bool =
    let
        af = f "Peter"
    in
    af score
`,
		`
[FnDef $f = [Fn [FnType [TypeReference $String], [TypeReference $Int] -> [TypeReference $Bool]] (name, score) = [If ($name == 'Peter') then ($score * #2) else $score]]]
[FnDef $another = [Fn [FnType [TypeReference $Int] -> [TypeReference $Bool]] (score) = [Let: [[LetAssign [$af] = [Call $f ['Peter']]]] in [Call $af [$score]]]]]
`)
}

func TestAnnotationParen(t *testing.T) {
	testParse(t,
		`
something : (String, Bool) -> Int
`, `
[FnDef $something = [Fn [FnType [TypeReference $String], [TypeReference $Bool] -> [TypeReference $Int]] (_, _) = [FnDeclExpr 0]]]
`)
}

func TestAnnotationParen2(t *testing.T) {
	testParse(t,
		`
something : ((String -> Bool)) -> (Int -> Bool)
`, `
[FnDef $something = [Fn [FnType [FnType [TypeReference $String] -> [TypeReference $Bool]] -> [FnType [TypeReference $Int] -> [TypeReference $Bool]]] (_) = [FnDeclExpr 0]]]
`)
}

func TestImport(t *testing.T) {
	testParse(t,
		`
import First
import Second.Sub
`, `
[Import [ModuleRef [First]]]
[Import [ModuleRef [Second Sub]]]
`)
}

func TestImportThird(t *testing.T) {
	testParse(t,
		`
import First
import Character.Damage.Sub
`, `
[Import [ModuleRef [First]]]
[Import [ModuleRef [Character Damage Sub]]]
`)
}

func TestImportExposingEllipsis(t *testing.T) {
	testParse(t,
		`
import First exposing (..)
`, `
[Import [ModuleRef [First]] exposing (..)]
`)
}

func TestImportFail(t *testing.T) {
	testParseError(t,
		`
import First
import Second.sub
`, parerr.ImportMustHaveUppercasePathError{})
}

func xTestImportFail2(t *testing.T) {
	testParseError(t,
		`
import First
import Second:sub
`, parerr.InternalError{})
}

func TestAppend(t *testing.T) {
	testParseExpression(t,
		`
[ 1, 3, 4 ] ++ [ 5, 6, 7, 8 ]
`, "([ListLiteral [#1 #3 #4]] ++ [ListLiteral [#5 #6 #7 #8]])")
}

func TestSimpleTypeDefinition(t *testing.T) {
	testParse(t,
		`
type alias Struct =
    { a : Int
    , b : Boolean
    }
`, `
[AliasType $Struct [RecordType [[Field: $a [TypeReference $Int]] [Field: $b [TypeReference $Boolean]]]]]
`)
}

func TestTypeDefinition(t *testing.T) {
	testParse(t,
		`
type alias Struct =
    { a : Int
    , b : Boolean
    }


a : Bool = somethingElse (b == 3)
`, `
[AliasType $Struct [RecordType [[Field: $a [TypeReference $Int]] [Field: $b [TypeReference $Boolean]]]]]
[FnDef $a = [Fn [FnType [TypeReference $Bool]] () = [Call $somethingElse [($b == #3)]]]]
`)
}

func TestSimpleTypeDefinitionConstructor(t *testing.T) {
	testParse(t,
		`
type alias Struct =
    { a : Int
    , b : Boolean
    }


a : (Int) -> Struct =
    Struct 2 false
`, `
[AliasType $Struct [RecordType [[Field: $a [TypeReference $Int]] [Field: $b [TypeReference $Boolean]]]]]
[FnDef $a = [Fn [FnType [TypeReference $Int] -> [TypeReference $Struct]] (_) = [CCall [TypeReference $Struct] [#2 €false]]]]
`)
}

func TestSimpleTypeDefinitionConstructorRecord(t *testing.T) {
	testParse(t,
		`
type alias Struct =
    { a : Int
    , b : Boolean
    }


a : (a: Int) -> Struct =
    Struct{ a = 2, b = False }
`, `
[AliasType $Struct [RecordType [[Field: $a [TypeReference $Int]] [Field: $b [TypeReference $Boolean]]]]]
[FnDef $a = [Fn [FnType [TypeReference $Int] -> [TypeReference $Struct]] (a) = [CCall [TypeReference $Struct] [[RecordLiteral [[$a = #2] [$b = [CCall [TypeReference $False]]]]]]]]]
`)
}

func TestFuncDeclarationAndDefinition(t *testing.T) {
	testParse(t,
		`
hello : (first: Int, c: Int) -> Int =
    4 + first
`,
		`
[FnDef $hello = [Fn [FnType [TypeReference $Int], [TypeReference $Int] -> [TypeReference $Int]] (first, c) = (#4 + $first)]]
`)
}

func TestSimpleCallWithLookup(t *testing.T) {
	testParseExpression(t,
		`
move sprite.rootPosition delta
`, "[Call $move [[RecordLookups $sprite [$rootPosition]] $delta]]")
}

func TestSimpleCallWithLookupInAssignmentBlock(t *testing.T) {
	testParseExpression(t,
		`
{ rootPosition = move sprite.rootPosition delta }
`, "[RecordLiteral [[$rootPosition = [Call $move [[RecordLookups $sprite [$rootPosition]] $delta]]]]]")
}

func TestOneLineIf(t *testing.T) {
	testParseExpression(t,
		`
if x == 3 then 4 else 5
`, "[If ($x == #3) then #4 else #5]")
}

func TestOneLineIfWithCall(t *testing.T) {
	testParseExpression(t,
		`
if x == 3 then
    extrude (5 * 4)
else
    5
`, "[If ($x == #3) then [Call $extrude [(#5 * #4)]] else #5]")
}

func TestBoolean(t *testing.T) {
	testParseExpression(t, "true", "€true")
}

func TestOneLineIfWithCallAndPipeRight(t *testing.T) {
	testParseExpression(t,
		`
if x == 3 then
    extrude (5 * 4) |> minimize true
else
    5
`, "[If ($x == #3) then ([Call $extrude [(#5 * #4)]] |> [Call $minimize [€true]]) else #5]")
}

func TestConstructorWithOneField(t *testing.T) {
	testParseExpression(t,
		`
{ first = 42 }
`, "[RecordLiteral [[$first = #42]]]")
}

func TestLookup(t *testing.T) {
	testParseExpression(t,
		`
a.b.c + d.e * f.g
`, "([RecordLookups $a [$b $c]] + ([RecordLookups $d [$e]] * [RecordLookups $f [$g]]))")
}

func TestConstructorWithTwoFields(t *testing.T) {
	testParseExpression(t,
		`
{ first = 42, second = 99 }
`, "[RecordLiteral [[$first = #42] [$second = #99]]]")
}

func TestConstructorWithTwoFieldsInDescendingOrder(t *testing.T) {
	testParseExpression(t,
		`
{ z = 42, andor = 99 }
`, "[RecordLiteral [[$z = #42] [$andor = #99]]]")
}

func TestConstructorWithSeveralFields(t *testing.T) {
	testParseExpression(t,
		`
{ first = 42, second = 99, third = 13 }
`, "[RecordLiteral [[$first = #42] [$second = #99] [$third = #13]]]")
}

func TestUnary1(t *testing.T) {
	testParse(t,
		`
a : (x: Bool) -> Bool =
    !x
`, `
[FnDef $a = [Fn [FnType [TypeReference $Bool] -> [TypeReference $Bool]] (x) = (NOT $x)]]
`)
}

func TestUnary2(t *testing.T) {
	testParse(t,
		`
someTest : (a: Bool, b: Bool) -> Bool =
    !a && b
`, `
[FnDef $someTest = [Fn [FnType [TypeReference $Bool], [TypeReference $Bool] -> [TypeReference $Bool]] (a, b) = ((NOT $a) AND $b)]]
`)
}

func TestUnary3(t *testing.T) {
	testParse(t,
		`
someTest : (a: Bool, b: Bool) -> Bool =
    a && !b
`, `
[FnDef $someTest = [Fn [FnType [TypeReference $Bool], [TypeReference $Bool] -> [TypeReference $Bool]] (a, b) = ($a AND (NOT $b))]]
`)
}

func TestList(t *testing.T) {
	testParseExpression(t,
		`
[ 2, 4, b, 101 ]
`, "[ListLiteral [#2 #4 $b #101]]")
}

func TestTuple(t *testing.T) {
	testParseExpression(t,
		`
( 2, 4, 4 )
`, "[TupleLiteral [#2 #4 #4]]")
}

func TestTupleType(t *testing.T) {
	testParse(t,
		`
someFunc : Int -> String -> Int =
    (42, "hi")
`, `
[FnDef $someFunc = [Fn [FnType [FnType [TypeReference $Int], [TypeReference $String] -> [TypeReference $Int]]] () = [TupleLiteral [#42 'hi']]]]
`)
}

func TestModuleReferences(t *testing.T) {
	testParseExpression(t,
		`
FirstModule.SecondMod.someFunc 12
`, "[Call FirstModule.SecondMod.$someFunc [#12]]")
}

func TestModuleReferenceWithType(t *testing.T) {
	testParse(t,
		`
a : (Bool) -> FirstModule.SecondMod.ThisIsAType
`, "[FnDef $a = [Fn [FnType [TypeReference $Bool] -> [ScopedTypeReference FirstModule.SecondMod.$ThisIsAType]] (_) = [FnDeclExpr 0]]]")
}

func TestModuleReferenceWithInitializer(t *testing.T) {
	testParseExpression(t,
		`
{ x = 2 }
`, "[RecordLiteral [[$x = #2]]]")
}

func TestMoreComplexList(t *testing.T) {
	testParseExpression(t,
		`
someFunc [ 2, 4, b, 101 ]
`, "[Call $someFunc [[ListLiteral [#2 #4 $b #101]]]]")
}

func TestTwoListType(t *testing.T) {
	testParse(t,
		`
someFunc : (List Sprite) -> List Another
`, "[FnDef $someFunc = [Fn [FnType [TypeReference $List [[TypeReference $Sprite]]] -> [TypeReference $List [[TypeReference $Another]]]] (_) = [FnDeclExpr 0]]]")
}

func TestMultipleListType(t *testing.T) {
	testParse(t,
		`
someFunc : (List Sprite -> List Another) -> List Something
`, "[FnDef $someFunc = [Fn [FnType [FnType [TypeReference $List [[TypeReference $Sprite]]] -> [TypeReference $List [[TypeReference $Another]]]] -> [TypeReference $List [[TypeReference $Something]]]] (_) = [FnDeclExpr 0]]]")
}

func TestListLiteral(t *testing.T) {
	testParse(t,
		`
type alias Cool =
    { name : String
    }


a : (x: Bool) -> List Cool =
    [ { name = "hi" }, { name = "another" }, { name = "tjoho" } ]
`, `
[AliasType $Cool [RecordType [[Field: $name [TypeReference $String]]]]]
[FnDef $a = [Fn [FnType [TypeReference $Bool] -> [TypeReference $List [[TypeReference $Cool]]]] (x) = [ListLiteral [[RecordLiteral [[$name = 'hi']]] [RecordLiteral [[$name = 'another']]] [RecordLiteral [[$name = 'tjoho']]]]]]]
`)
}

func TestMultipleStatements(t *testing.T) {
	testParse(t,
		`
someFunc : Int =
    let
        i = 3

        j = 4
    in
    if i >= 6 then
        call (i + j)
    else
        3
`, "[FnDef $someFunc = [Fn [FnType [TypeReference $Int]] () = [Let: [[LetAssign [$i] = #3] [LetAssign [$j] = #4]] in [If ($i >= #6) then [Call $call [($i + $j)]] else #3]]]]")
}

func TestMoreMultipleStatements(t *testing.T) {
	testParse(t,
		`
someFunc : Int =
    let
        i = 3

        j = 4
    in
    i + 2


anotherFunc : Int =
    if c + 2 <= 3 then
        18
    else
        19
`, `
[FnDef $someFunc = [Fn [FnType [TypeReference $Int]] () = [Let: [[LetAssign [$i] = #3] [LetAssign [$j] = #4]] in ($i + #2)]]]
[FnDef $anotherFunc = [Fn [FnType [TypeReference $Int]] () = [If (($c + #2) <= #3) then #18 else #19]]]
`)
}

func TestMultilineIf(t *testing.T) {
	testParseExpression(t,
		`
if x == 3 then
    callme (4 + 4)
else
    5
`, "[If ($x == #3) then [Call $callme [(#4 + #4)]] else #5]")
}

func TestSimpleLetOneLine(t *testing.T) {
	testParseExpression(t,
		`
let
    x = 55 + 41 * 4
in
2 + x
`, "[Let: [[LetAssign [$x] = (#55 + (#41 * #4))]] in (#2 + $x)]")
}

func TestSimpleLetInOnAnotherLine(t *testing.T) {
	testParseExpression(t,
		`
let
    x = 55 + 41 * 4
in
2 + x
`, "[Let: [[LetAssign [$x] = (#55 + (#41 * #4))]] in (#2 + $x)]")
}

func TestAnnotationThatIsUnfinished(t *testing.T) {
	testParseError(t,
		`
checkIfBasicAttack : (Character) ->
`, parerr.ExpectedTypeReferenceError{})
}

func TestSimpleLet(t *testing.T) { // --- FIXME
	testParseExpression(t,
		`
let
    x = 55 + 1
in
x + 2
`, "[Let: [[LetAssign [$x] = (#55 + #1)]] in ($x + #2)]")
}

func TestLet(t *testing.T) { // --- FIXME
	testParseExpression(t,
		`
let
    x = callme (4 + 4)

    y = 3

    zarg = 8 * 99
in
x + y
`, "[Let: [[LetAssign [$x] = [Call $callme [(#4 + #4)]]] [LetAssign [$y] = #3] [LetAssign [$zarg] = (#8 * #99)]] in ($x + $y)]")
}

func TestLetSubSameLine(t *testing.T) { // --- FIXME
	testParseExpression(t,
		`
let
    x = callme (4 + 5)
in
let
    bogus = 4
in
x + y * bogus
`, "[Let: [[LetAssign [$x] = [Call $callme [(#4 + #5)]]]] in [Let: [[LetAssign [$bogus] = #4]] in ($x + ($y * $bogus))]]")
}

func TestPyth(t *testing.T) { // --- FIXME
	testParse(t,
		`
pythagoras : (ax: Int, ay: Int, bx: Int, by: Int) -> Int =
    let
        x = ax - bx

        y = ay - by
    in
    Math.sqrt (x * x + y * y)
`, `
[FnDef $pythagoras = [Fn [FnType [TypeReference $Int], [TypeReference $Int], [TypeReference $Int], [TypeReference $Int] -> [TypeReference $Int]] (ax, ay, bx, by) = [Let: [[LetAssign [$x] = ($ax - $bx)] [LetAssign [$y] = ($ay - $by)]] in [Call Math.$sqrt [(($x * $x) + ($y * $y))]]]]]
`)
}

func TestGreat(t *testing.T) {
	testParse(t,
		`
isGreat? : (a: Int) -> Bool =
    a > 99
`, `
[FnDef $isGreat? = [Fn [FnType [TypeReference $Int] -> [TypeReference $Bool]] (a) = ($a > #99)]]
`)
}

func TestCase(t *testing.T) {
	testParseExpression(t,
		`
case x of
    Int i -> itoa i

    String s -> s
`, "[CaseCustomType $x of [CaseConsCustomType $Int ([$i]) => [Call $itoa [$i]]];[CaseConsCustomType $String ([$s]) => $s]]")
}

func TestAlias(t *testing.T) {
	testParse(t,
		`
type alias State =
    { playerX : Int
    , time : Int
    }

`, "[AliasType $State [RecordType [[Field: $playerX [TypeReference $Int]] [Field: $time [TypeReference $Int]]]]]")
}

func TestAliasOnSeveralLines(t *testing.T) {
	testParse(t,
		`
type alias State =
    { playerX : Int
    , time : Int
    }

`, "[AliasType $State [RecordType [[Field: $playerX [TypeReference $Int]] [Field: $time [TypeReference $Int]]]]]")
}

func TestCaseWithDefault(t *testing.T) {
	testParseExpression(t,
		`
case x of
    _ -> 2
`, "[CasePm $x of [CaseConsPm '_' => #2]]")
}

func TestCustomType(t *testing.T) {
	testParse(t,
		`
type SomeEnum =
    First String
    | Anon
    | Second Int

`, "[CustomTypeStatement [CustomType $SomeEnum [[Variant $First [[TypeReference $String]]] [Variant $Anon []] [Variant $Second [[TypeReference $Int]]]]]]")
}

func TestCustomTypeNewFormatting(t *testing.T) {
	testParse(t,
		`
type SomeEnum =
    First String
    | Anon
    | Second Int

`, "[CustomTypeStatement [CustomType $SomeEnum [[Variant $First [[TypeReference $String]]] [Variant $Anon []] [Variant $Second [[TypeReference $Int]]]]]]")
}

func TestModuleVariantConstructor(t *testing.T) { // --- FIXME
	testParseExpression(t,
		`
Imaginary.Module.First "Hello"
`, `
[CCall [ScopedTypeReference Imaginary.Module.$First] ['Hello']]
`)
}

func TestCustomTypeVariantConstructor(t *testing.T) { // --- FIXME
	testParse(t,
		`
type SomeEnum =
    First String
    | Anon
    | Second Int


a : (Bool) -> SomeEnum =
    Imaginary.Module.First "Hello"
`, `
[CustomTypeStatement [CustomType $SomeEnum [[Variant $First [[TypeReference $String]]] [Variant $Anon []] [Variant $Second [[TypeReference $Int]]]]]]
[FnDef $a = [Fn [FnType [TypeReference $Bool] -> [TypeReference $SomeEnum]] (_) = [CCall [ScopedTypeReference Imaginary.Module.$First] ['Hello']]]]
`)
}

func TestCustomTypeWithCase(t *testing.T) {
	testParse(t,
		`
type SomeEnum =
    First String
    | Anon
    | Second Int


a : SomeEnum =
    case x of
        First i -> itoa i

        Second s -> s
`, `
[CustomTypeStatement [CustomType $SomeEnum [[Variant $First [[TypeReference $String]]] [Variant $Anon []] [Variant $Second [[TypeReference $Int]]]]]]
[FnDef $a = [Fn [FnType [TypeReference $SomeEnum]] () = [CaseCustomType $x of [CaseConsCustomType $First ([$i]) => [Call $itoa [$i]]];[CaseConsCustomType $Second ([$s]) => $s]]]]
`)
}

func TestMinusAssociativity(t *testing.T) {
	testParseExpression(t,
		`
5 - 1 - 2
`, "((#5 - #1) - #2)")
}

func TestPipeForward(t *testing.T) { // --- FIXME
	testParseExpression(t,
		`
call param1 param2
    |> anotherCall anotherParam1
    |> third (2 * 3 + 5 * 66)
`, "(([Call $call [$param1 $param2]] |> [Call $anotherCall [$anotherParam1]]) |> [Call $third [((#2 * #3) + (#5 * #66))]])")
}

func TestPipeBackward(t *testing.T) { // --- FIXME
	testParseExpression(t,
		`
call param1 param2 <| laterCall this <| third (2 * 3 + 5 * 66)
`, "(([Call $call [$param1 $param2]] <| [Call $laterCall [$this]]) <| [Call $third [((#2 * #3) + (#5 * #66))]])")
}

func TestPipeRight(t *testing.T) {
	testParseExpression(t,
		`
first (2 + 2) |> second
`, "([Call $first [(#2 + #2)]] |> $second)")
}

func TestNormalParen(t *testing.T) {
	testParseExpression(t,
		`
( 2 + 4 )
`, "(#2 + #4)")
}

func TestAloneUpdate(t *testing.T) {
	testParseExpression(t,
		`
( objectToUpdate | someField = 3 )
`, "($objectToUpdate | ($someField = #3))")
}

func TestUpdate(t *testing.T) {
	testParse(t,
		`
 a =
    { objectToUpdate | someField = 3 }
`, "[Constant $a = [RecordUpdate [[$someField = #3]] ($objectToUpdate)]]")
}

func TestUpdateWithSpace(t *testing.T) {
	testParse(t,
		`
 a =
    { objectToUpdate | someField = 3 }
`, "[Constant $a = [RecordUpdate [[$someField = #3]] ($objectToUpdate)]]")
}

func TestUpdateWithCallAndSpace(t *testing.T) { // --- FIXME
	testParse(t,
		`
 a =
    { inSprite | scale = scaleBoth scaleFactor }
`, "[Constant $a = [RecordUpdate [[$scale = [Call $scaleBoth [$scaleFactor]]]] ($inSprite)]]")
}

func TestComplexUpdate(t *testing.T) {
	testParseExpression(t,
		`
{ inSprite | scale = { scaleX = scale, scaleY = scale } }
`, "[RecordUpdate [[$scale = [RecordLiteral [[$scaleX = $scale] [$scaleY = $scale]]]]] ($inSprite)]")
}

func xTestGenericsError(t *testing.T) {
	testParseError(t,
		`
type alias Tinkering t =
    { solder : Bool
    }
`, &ast.ExtraTypeParameterError{})
}

func xTestGenericsNotDefinedError(t *testing.T) {
	testParseError(t,
		`
type alias Tinkering a =
    { solder : Bool
    , used : a
    , something : t
    }
`, &ast.UndefinedTypeParameterError{})
}

func xTestBasicGenerics(t *testing.T) {
	testParse(t,
		`
type alias Tinkering t =
    { solder : Bool
    , cool : Something t
    }

`, "[AliasType $Tinkering [RecordType [[Field: $solder [TypeReference $Bool]] [Field: $cool [TypeReference $Something [[GenericType [TypeParam $t]]]]]] [[TypeParam $t]]]]")
}

func xTestMultipleGenerics(t *testing.T) {
	testParse(t,
		`
type alias Tinkering t a =
    { solder : Bool
    , cool : List t
    , other : a
    }

`, `
[AliasType $Tinkering [RecordType [[Field: $solder [TypeReference $Bool]] [Field: $cool [TypeReference $List [[GenericType [TypeParam $t]]]]] [Field: $other [GenericType [TypeParam $a]]]] [[TypeParam $t] [TypeParam $a]]]]
`)
}

func TestSkippingModule(t *testing.T) {
	testParse(t,
		`
module Main exposing (main)

import Character
import Characters.Update
`, `
[FnDef $module = [Fn [TypeParamContext [exposing main] +> [FnType [TypeReference $Main [[LocalTypeNameRef [LocalTypeNameDef [LocalTypeName exposing]]] [LocalTypeNameRef [LocalTypeNameDef [LocalTypeName main]]]]]]] () = [FnDeclExpr 0]]]
[Import [ModuleRef [Character]]]
[Import [ModuleRef [Characters Update]]]
`)
}

func xTestGenericsAnnotation(t *testing.T) {
	testParse(t,
		`
type alias Tinkering t =
    { solder : Bool
    , cool : t
    }


f : (a: Tinkering t) -> Int =
    2
`, `
[AliasType $Tinkering [RecordType [[Field: $solder [TypeReference $Bool]] [Field: $cool [GenericType [TypeParam $t]]]] [[TypeParam $t]]]]
[FnDef $f = [Fn ([[Arg $a: [TypeReference $Tinkering [[GenericType [TypeParam $t]]]]]]) => [TypeReference $Int] = #2]]
`)
}

func TestGenericsCustomType(t *testing.T) {
	testParse(t,
		`
type Maybe a =
    Nothing
    | Just a
`, `[CustomTypeStatement [TypeParamContext [a] +> [CustomType $Maybe [[Variant $Nothing []] [Variant $Just [[LocalTypeNameRef [LocalTypeNameDef [LocalTypeName a]]]]]]]]]
`)
}

func TestGenericsAnnotationBegin(t *testing.T) {
	testParse(t,
		`
ownCons : itemType -> List itemType -> List itemType
`, `
[FnDef $ownCons = [Fn [TypeParamContext [itemType] +> [FnType [FnType [LocalTypeNameRef [LocalTypeNameDef [LocalTypeName itemType]]], [TypeReference $List [[LocalTypeNameRef [LocalTypeNameDef [LocalTypeName itemType]]]]] -> [TypeReference $List [[LocalTypeNameRef [LocalTypeNameDef [LocalTypeName itemType]]]]]]]] () = [FnDeclExpr 0]]]
`)
}

func xTestGenericsAnnotationSpecified(t *testing.T) {
	testParse(t,
		`
type alias Tinkering t =
    { solder : Bool
    , secret : t
    }


f : (Tinkering Int) -> Int =
    tinkering.secret
`, `
[AliasType $Tinkering [RecordType [[Field: $solder [TypeReference $Bool]] [Field: $secret [GenericType [TypeParam $t]]]] [[TypeParam $t]]]]
[FnDef $f = [Fn ([]) => [FnType [TypeReference $Tinkering [[TypeReference $Int]]] -> [TypeReference $Int]] = [RecordLookups $tinkering [$secret]]]]
`)
}

func TestGenericsAnnotation2(t *testing.T) {
	testParse(t,
		`
f : Tinkering Int -> Int =
    tinkering.secret
`, `
[FnDef $f = [Fn [FnType [FnType [TypeReference $Tinkering [[TypeReference $Int]]] -> [TypeReference $Int]]] () = [RecordLookups $tinkering [$secret]]]]
`)
}

func xTestGenericsStructInstantiate(t *testing.T) {
	testParse(t,
		`
type alias Tinkering t =
    { solder : Bool
    , secret : t
    }
`, "[AliasType $Tinkering [RecordType [[Field: $solder [TypeReference $Bool]] [Field: $secret [GenericType [TypeParam $t]]]] [[TypeParam $t]]]]")
}

func TestCaseNewLine(t *testing.T) {
	testParseExpression(t,
		`
case a of
    Nothing ->
        None

    Just oldGamepad ->
        2
		`, `
[CaseCustomType $a of [CaseConsCustomType $Nothing ([]) => [CCall [TypeReference $None]]];[CaseConsCustomType $Just ([$oldGamepad]) => #2]]
`)
}

func TestMinimalTypeConstructor(t *testing.T) {
	testParseExpression(t,
		`
case a of
    Nothing ->
        None

    Just oldGamepad ->
        2
		`, `
[CaseCustomType $a of [CaseConsCustomType $Nothing ([]) => [CCall [TypeReference $None]]];[CaseConsCustomType $Just ([$oldGamepad]) => #2]]
`)
}

func xTestRecordAnnotation2(t *testing.T) {
	testParse(t,
		`
type alias Sprite =
    { solder : Bool
    , cool : Int
    }


f : (a: Sprite) -> Int =
    2


g : (a : { solder : Bool, cool : Int }) -> Int =
    f { solder = True, cool = 33 }
`, `
[alias $Sprite [record-type [[field: $solder [type-reference $Bool]] [field: $cool [type-reference $Int]]]]]
[annotation: $f [func-type [type-reference $Sprite] -> [type-reference $Int]]]
[definition: $f = [func ([$a]) -> #2]]
[annotation: $f [func-type [record-type [[field: $solder [type-reference $Bool]] [field: $cool [type-reference $Int]]]] -> [type-reference $Int]]]
[definition: $k = [func ([$a]) -> [call $f [[record-literal: [[$solder = €true] [$cool = #33]]]]]]]
`)
}

func TestRecordAnnotation3(t *testing.T) {
	testParse(t,
		`
f : (a: Int) -> { solder : Bool, cool : Maybe Int } =
    2
`, `
[FnDef $f = [Fn [FnType [TypeReference $Int] -> [RecordType [[Field: $solder [TypeReference $Bool]] [Field: $cool [TypeReference $Maybe [[TypeReference $Int]]]]]]] (a) = #2]]
`)
}

func TestPipeRight2(t *testing.T) {
	testParse(t, `
tester : (b: String) -> Bool =
    first (2 + 2) |> second b |> third
`, `
[FnDef $tester = [Fn [FnType [TypeReference $String] -> [TypeReference $Bool]] (b) = (([Call $first [(#2 + #2)]] |> [Call $second [$b]]) |> $third)]]
`)
}
