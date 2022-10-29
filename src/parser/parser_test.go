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
[FnDef $first = [Fn ([[Arg $startNumber: [TypeReference $Int]] [Arg $somethingElse: [FnType [TypeReference $String] -> [TypeReference $List [[GenericType [TypeParam $a]]]]]]]) => [TupleType [[TypeReference $List [[TypeReference $List [[TypeReference $Fixed]]]]] [TypeReference $Int]]] = ($startNumber * $startNumber)]]
`)
}

func TestExternalVar(t *testing.T) {
	testParse(t,
		`
__externalvarfn head : List a -> Maybe a
`,
		`
[FnDef $head = [Fn ([]) => [FnType [TypeReference $List [[GenericType [TypeParam $a]]]] -> [TypeReference $Maybe [[GenericType [TypeParam $a]]]]] = [EmptyExpression]]]
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

func TestMultilineCommentOnOneLine(t *testing.T) {
	testParseExpression(t,
		`{- this is basically a comment -}`,
		`[multilinecomment ' this is basically a comment ']`)
}

func TestMultilineComment(t *testing.T) {
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
[FnDef $fn = [Fn ([[Arg $x: [TypeReference $Int]]]) => [TypeReference $Int] = [Guard [[($x < #2) => #42] [($x == #4) => #45]] [_ => #-1]]]]
`)
}

func TestWildcardType(t *testing.T) {
	testParse(t,
		`
fn : * -> Int
fn a b x =
    23	
`,

		`
[annotation: $fn [func-type [AnyMatchingType *] -> [TypeReference $Int]]]
[fndefinition: $fn = [func ([$a $b $x]) -> #23]]
`)
}

func TestResourceName(t *testing.T) {
	testParseExpression(t,
		`
@some/name/here2.png
`,
		`@some/name/here2.png`)
}

func TestSingleLineComment(t *testing.T) {
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
		`True`,
		`€true`)
}

func TestEmptyList(t *testing.T) {
	testParseExpression(t,
		`[]`,
		`[list-literal: []]`)
}

func TestIf(t *testing.T) {
	testParseExpression(t,
		`if True then 1 else 0`,
		`[If €true then #1 else #0]`)
}

func TestIfWithNewLine(t *testing.T) {
	testParseExpression(t,
		`
if True then
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
if True then
    if False then
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
if True then
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
if True then
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
if True then
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
		`[Call Array.$fromList [[list-literal: [[Call Array.$fromList [[list-literal: [#0 #1 #2 #3]] Array.$fromList [list-literal: [#8 #9 #10 #11]]]]]]]]`)
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
		`[list-literal: [#3 (#4 + #5)]]`)
}

func TestArrayLiteralTwoItems(t *testing.T) {
	testParseExpression(t,
		`
[| 3, 4 + 5 |]
`,
		`[array-literal: [#3 (#4 + #5)]]`)
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
[FnDef $f = [Fn ([[Arg $name: [TypeReference $String]] [Arg $score: [TypeReference $Int]]]) => [TypeReference $Bool] = [If ($name == 'Peter') then ($score * #2) else $score]]]
[FnDef $another = [Fn ([[Arg $score: [TypeReference $Int]]]) => [TypeReference $Bool] = [Let: [[LetAssign [$af] = [Call $f ['Peter']]]] in [Call $af [$score]]]]]
`)
}

func TestAnnotationParen(t *testing.T) {
	testParse(t,
		`
something : (String, Bool) -> Int
`, `
[annotation: $something [func-type [func-type [type-reference $String] -> [type-reference $Bool]] -> [type-reference $Int]]]
`)
}

func TestAnnotationParen2(t *testing.T) {
	testParse(t,
		`
something : (String -> Bool) -> (Int -> Bool)
`, `
[annotation: $something [func-type [func-type [type-reference $String] -> [type-reference $Bool]] -> [func-type [type-reference $Int] -> [type-reference $Bool]]]]
`)
}

func TestLambdaMultipleParameters(t *testing.T) {
	testParseExpression(t,
		`
\x test b -> hello 2
`,
		`[lambda ([$x $test $b]) -> [call $hello [#2]]]`)
}

func TestImport(t *testing.T) {
	testParse(t,
		`
import First
import Second.Sub
`, `
[import [$First]]
[import [$Second $Sub]]
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
[import [$First] exposing (..)]
`)
}

func TestImportFail(t *testing.T) {
	testParseError(t,
		`
import First
import Second.sub
`, parerr.ImportMustHaveUppercasePathError{})
}

func TestImportFail2(t *testing.T) {
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
`, "([list-literal: [#1 #3 #4]] ++ [list-literal: [#5 #6 #7 #8]])")
}

func TestSimpleTypeDefinition(t *testing.T) {
	testParse(t,
		`
type alias Struct =
    { a : Int
    , b : Boolean
    }
`, `
[alias $Struct [record-type [[field: $a [type-reference $Int]] [field: $b [type-reference $Boolean]]]]]
`)
}

func TestTypeDefinition(t *testing.T) {
	testParse(t,
		`
type alias Struct =
    { a : Int
    , b : Boolean
    }


a =
    somethingElse (b == 3)
`, `
[alias $Struct [record-type [[field: $a [type-reference $Int]] [field: $b [type-reference $Boolean]]]]]
[definition: $a = [func ([]) -> [call $somethingElse [($b == #3)]]]]
`)
}

func TestSimpleTypeDefinitionConstructor(t *testing.T) {
	testParse(t,
		`
type alias Struct =
    { a : Int
    , b : Boolean
    }


a : Int -> Struct
a ignore =
    Struct 2 False
`, `
[AliasType $Struct [record-type [[Field: $a [TypeReference $Int]] [Field: $b [TypeReference $Boolean]]] []]]
[Annotation $a [FuncType [TypeReference $Int] -> [TypeReference $Struct]]]
[FnDef $a = [Func ([$ignore]) -> [CCall [TypeReference $Struct] [#2 €false]]]]
`)
}

func TestSimpleTypeDefinitionConstructorRecord(t *testing.T) {
	testParse(t,
		`
type alias Struct =
    { a : Int
    , b : Boolean
    }


a : Int -> Struct
a ignore =
    Struct{ a = 2, b = False }
`, `
[AliasType $Struct [RecordType [[Field: $a [TypeReference $Int]] [Field: $b [TypeReference $Boolean]]] []]]
[Annotation $a [FnType [TypeReference $Int] -> [TypeReference $Struct]]]
[FnDef $a = [Fn ([$ignore]) -> [CCall [TypeReference $Struct] [[RecordLiteral [[$a = #2] [$b = €false]]]]]]]
`)
}

func TestFuncDeclarationAndDefinition(t *testing.T) {
	testParse(t,
		`
hello : Int -> Int -> Int
hello first c =
    4 + first
`,
		`
[Annotation $hello [FnType [TypeReference $Int] -> [TypeReference $Int] -> [TypeReference $Int]]]
[FnDef $hello = [Fn ([$first $c]) -> (#4 + $first)]]
`)
}

func TestDefinitionWithCallWithRecords(t *testing.T) {
	testParse(t,
		`
main a =
    { sprite = { x = calc 10 } }
`,
		`
[FnDef $main = [Fn ([$a]) -> [RecordLiteral [[$sprite = [RecordLiteral [[$x = [Call $calc [#10]]]]]]]]]]
`)
}

func TestSimpleCall(t *testing.T) {
	testParse(t,
		`
rebecca is fantastic =
    something (3 * 3) (4 + 4)
`, "[FnDef $rebecca = [Fn ([$is $fantastic]) -> [Call $something [(#3 * #3) (#4 + #4)]]]]")
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
`, "[if: ($x == #3) then #4 else #5]")
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
	testParseExpression(t, "True", "€true")
}

func TestOneLineIfWithCallAndPipeRight(t *testing.T) {
	testParseExpression(t,
		`
if x == 3 then
    extrude (5 * 4) |> minimize true
else
    5
`, "[If ($x == #3) then ([Call $extrude [(#5 * #4)]] |> [Call $minimize [$true]]) else #5]")
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
`, "[record-literal: [[$z = #42] [$andor = #99]]]")
}

func TestConstructorWithSeveralFields(t *testing.T) {
	testParseExpression(t,
		`
{ first = 42, second = 99, third = 13 }
`, "[record-literal: [[$first = #42] [$second = #99] [$third = #13]]]")
}

func TestUnary1(t *testing.T) {
	testParse(t,
		`
a : Bool -> Bool
a x =
    !x
`, `
[Annotation $a [FnType [TypeReference $Bool] -> [TypeReference $Bool]]]
[FnDef $a = [Fn ([$x]) -> ( $x)]]
`)
}

func TestUnary2(t *testing.T) {
	testParse(t,
		`
someTest : Bool -> Bool -> Bool
someTest a b =
    !a && b
`, `
[annotation: $someTest [func-type [type-reference $Bool] -> [type-reference $Bool] -> [type-reference $Bool]]]
[definition: $someTest = [func ([$a $b]) -> ((! $a) AND $b)]]
`)
}

func TestUnary3(t *testing.T) {
	testParse(t,
		`
someTest : Bool -> Bool -> Bool
someTest a b =
    a && !b
`, `
[annotation: $someTest [func-type [type-reference $Bool] -> [type-reference $Bool] -> [type-reference $Bool]]]
[definition: $someTest = [func ([$a $b]) -> ($a AND (! $b))]]
`)
}

func TestList(t *testing.T) {
	testParseExpression(t,
		`
[ 2, 4, b, 101 ]
`, "[list-literal: [#2 #4 $b #101]]")
}

func TestTuple(t *testing.T) {
	testParseExpression(t,
		`
( 2, 4, 4 )
`, "[tuple-literal: [#2 #4 #4]]")
}

func TestTupleType(t *testing.T) {
	testParse(t,
		`
someFunc : (Int, String) -> Int
someFunc =
    (42, "hi")
`, `
[Annotation $someFunc [FnType [TupleType [[TypeReference $Int] [TypeReference $String]]] -> [TypeReference $Int]]]
[FnDef $someFunc = [Fn ([]) -> [TupleLiteral [#42 'hi']]]]
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
a : Bool -> FirstModule.SecondMod.ThisIsAType
`, "[Annotation $a [FnType [TypeReference $Bool] -> [ScopedTypeReference FirstModule.SecondMod.$ThisIsAType]]]")
}

func TestAsm(t *testing.T) {
	testParse(t,
		`
a : Bool -> FirstModule.SecondMod.ThisIsAType
a ignore =
    __asm curry 2 4 ([3])
`, `
[annotation: $a [func-type [type-reference $Bool] -> [type-reference FirstModule.SecondMod.$ThisIsAType]]]
[definition: $a = [func ([$ignore]) -> [asm: curry 2 4 ([3])]]]
`)
}

func TestModuleReferenceWithInitializer(t *testing.T) {
	testParseExpression(t,
		`
{ x = 2 }
`, "[record-literal: [[$x = #2]]]")
}

func TestMoreComplexList(t *testing.T) {
	testParseExpression(t,
		`
someFunc [ 2, 4, b, 101 ]
`, "[call $someFunc [[list-literal: [#2 #4 $b #101]]]]")
}

func TestTwoListType(t *testing.T) {
	testParse(t,
		`
someFunc : List Sprite -> List Another
`, "[annotation: $someFunc [func-type [type-reference $List [[type-reference $Sprite]]] -> [type-reference $List [[type-reference $Another]]]]]")
}

func TestMultipleListType(t *testing.T) {
	testParse(t,
		`
someFunc : List Sprite -> List Another -> List Something
`, "[annotation: $someFunc [func-type [type-reference $List [[type-reference $Sprite]]] -> [type-reference $List [[type-reference $Another]]] -> [type-reference $List [[type-reference $Something]]]]]")
}

func TestListLiteral(t *testing.T) {
	testParse(t,
		`
type alias Cool =
    { name : String
    }


a : Bool -> List Cool
a x =
    [ { name = "hi" }, { name = "another" }, { name = "tjoho" } ]
`, `
[alias $Cool [record-type [[field: $name [type-reference $String]]]]]
[annotation: $a [func-type [type-reference $Bool] -> [type-reference $List [[type-reference $Cool]]]]]
[definition: $a = [func ([$x]) -> [list-literal: [[record-literal: [[$name = 'hi']]] [record-literal: [[$name = 'another']]] [record-literal: [[$name = 'tjoho']]]]]]]
`)
}

func TestMultipleStatements(t *testing.T) {
	testParse(t,
		`
someFunc a b =
    let
        i = 3

        j = 4
    in
    if i >= 6 then
        call (i + j)
    else
        3
`, "[definition: $someFunc = [func ([$a $b]) -> [let: [[LetAssign $i = #3] [letassign $j = #4]] in [if: ($i >= #6) then [call $call [($i + $j)]] else #3]]]]")
}

func TestMoreMultipleStatements(t *testing.T) {
	testParse(t,
		`
someFunc a b =
    let
        i = 3

        j = 4
    in
    i + 2


anotherFunc c =
    if c + 2 <= 3 then
        18
    else
        19
`, `
[definition: $someFunc = [func ([$a $b]) -> [let: [[letassign $i = #3] [letassign $j = #4]] in ($i + #2)]]]
[definition: $anotherFunc = [func ([$c]) -> [if: (($c + #2) <= #3) then #18 else #19]]]
`)
}

func TestMultilineIf(t *testing.T) {
	testParseExpression(t,
		`
if x == 3 then
    callme (4 + 4)
else
    5
`, "[if: ($x == #3) then [call $callme [(#4 + #4)]] else #5]")
}

func TestSimpleLetOneLine(t *testing.T) {
	testParseExpression(t,
		`
let
    x = 55 + 41 * 4
in
2 + x
`, "[let: [[letassign $x = (#55 + (#41 * #4))]] in (#2 + $x)]")
}

func TestSimpleLetInOnAnotherLine(t *testing.T) {
	testParseExpression(t,
		`
let
    x = 55 + 41 * 4
in
2 + x
`, "[let: [[letassign $x = (#55 + (#41 * #4))]] in (#2 + $x)]")
}

func xTestAnnotationThatIsUnfinished(t *testing.T) {
	testParseError(t,
		`
checkIfBasicAttack : Character ->
`, parerr.ExpectedTypeReferenceError{})
}

func TestSimpleLet(t *testing.T) { // --- FIXME
	testParseExpression(t,
		`
let
    x = 55 + 1
in
x + 2
`, "[let: [[letassign $x = (#55 + #1)]] in ($x + #2)]")
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
`, "[let: [[letassign $x = [call $callme [(#4 + #4)]]] [letassign $y = #3] [letassign $zarg = (#8 * #99)]] in ($x + $y)]")
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
`, "[let: [[letassign $x = [call $callme [(#4 + #5)]]]] in [let: [[letassign $bogus = #4]] in ($x + ($y * $bogus))]]")
}

func TestPyth(t *testing.T) { // --- FIXME
	testParse(t,
		`
pythagoras : Int -> Int -> Int -> Int -> Int
pythagoras ax ay bx by =
    let
        x = ax - bx

        y = ay - by
    in
    Math.sqrt (x * x + y * y)
`, `
[annotation: $pythagoras [func-type [type-reference $Int] -> [type-reference $Int] -> [type-reference $Int] -> [type-reference $Int] -> [type-reference $Int]]]
[definition: $pythagoras = [func ([$ax $ay $bx $by]) -> [let: [[letassign $x = ($ax - $bx)] [letassign $y = ($ay - $by)]] in [call Math.$sqrt [(($x * $x) + ($y * $y))]]]]]`)
}

func TestGreat(t *testing.T) {
	testParse(t,
		`
isGreat? : Int -> Bool
isGreat? a =
    a > 99
`, `
[annotation: $isGreat? [func-type [type-reference $Int] -> [type-reference $Bool]]]
[definition: $isGreat? = [func ([$a]) -> ($a > #99)]]`)
}

func TestCase(t *testing.T) {
	testParseExpression(t,
		`
case x of
    Int i -> itoa i

    String s -> s
`, "[case: $x of [casecons $Int ([$i]) => [call $itoa [$i]]];[casecons $String ([$s]) => $s]]")
}

func TestAlias(t *testing.T) {
	testParse(t,
		`
type alias State =
    { playerX : Int
    , time : Int
    }

`, "[alias $State [record-type [[field: $playerX [type-reference $Int]] [field: $time [type-reference $Int]]]]]")
}

func TestAliasOnSeveralLines(t *testing.T) {
	testParse(t,
		`
type alias State =
    { playerX : Int
    , time : Int
    }

`, "[alias $State [record-type [[field: $playerX [type-reference $Int]] [field: $time [type-reference $Int]]]]]")
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

`, "[CustomType $SomeEnum [[CustomTypeVariant $First [[TypeReference $String]]] [CustomTypeVariant $Anon []] [CustomTypeVariant $Second [[TypeReference $Int]]]]]")
}

func TestCustomTypeNewFormatting(t *testing.T) {
	testParse(t,
		`
type SomeEnum =
    First String
    | Anon
    | Second Int

`, "[CustomType $SomeEnum [[CustomTypeVariant $First [[TypeReference $String]]] [CustomTypeVariant $Anon []] [CustomTypeVariant $Second [[TypeReference $Int]]]]]")
}

func TestExternalFunction(t *testing.T) { // --- FIXME
	testParse(t,
		`
__externalfn coreListMap 2
`, `
[external function: coreListMap 2]
`)
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


a : Bool -> SomeEnum
a dummy =
    Imaginary.Module.First "Hello"
`, `
[CustomType $SomeEnum [[CustomTypeVariant $First [[TypeReference $String]]] [CustomTypeVariant $Anon []] [CustomTypeVariant $Second [[TypeReference $Int]]]]]
[Annotation $a [FnType [TypeReference $Bool] -> [TypeReference $SomeEnum]]]
[FnDef $a = [Fn ([$dummy]) -> [CCall [ScopedTypeReference Imaginary.Module.$First] ['Hello']]]]
`)
}

func TestCustomTypeWithCase(t *testing.T) {
	testParse(t,
		`
type SomeEnum =
    First String
    | Anon
    | Second Int


a =
    case x of
        First i -> itoa i

        Second s -> s
`, `
[custom-type-statement [custom-type $SomeEnum [[variant $First [[type-reference $String]]] [variant $Anon] [variant $Second [[type-reference $Int]]]]]]
[definition: $a = [func ([]) -> [case: $x of [casecons $First ([$i]) => [call $itoa [$i]]];[casecons $Second ([$s]) => $s]]]]
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
`, "[definition: $a = [func ([]) -> [record-literal: [[$someField = #3]] ($objectToUpdate)]]]")
}

func TestUpdateWithSpace(t *testing.T) {
	testParse(t,
		`
 a =
    { objectToUpdate | someField = 3 }
`, "[definition: $a = [func ([]) -> [record-literal: [[$someField = #3]] ($objectToUpdate)]]]")
}

func TestUpdateWithCallAndSpace(t *testing.T) { // --- FIXME
	testParse(t,
		`
 a =
    { inSprite | scale = scaleBoth scaleFactor }
`, "[definition: $a = [func ([]) -> [record-literal: [[$scale = [call $scaleBoth [$scaleFactor]]]] ($inSprite)]]]")
}

func TestComplexUpdate(t *testing.T) {
	testParseExpression(t,
		`
{ inSprite | scale = { scaleX = scale, scaleY = scale } }
`, "[record-literal: [[$scale = [record-literal: [[$scaleX = $scale] [$scaleY = $scale]]]]] ($inSprite)]")
}

func TestGenericsError(t *testing.T) {
	testParseError(t,
		`
type alias Tinkering t =
    { solder : Bool
    }
`, &ast.ExtraTypeParameterError{})
}

func TestGenericsNotDefinedError(t *testing.T) {
	testParseError(t,
		`
type alias Tinkering a =
    { solder : Bool
    , used : a
    , something : t
    }
`, &ast.UndefinedTypeParameterError{})
}

func TestBasicGenerics(t *testing.T) {
	testParse(t,
		`
type alias Tinkering t =
    { solder : Bool
    , cool : Something t
    }

`, "[AliasType $Tinkering [RecordType [[Field: $solder [TypeReference $Bool]] [Field: $cool [TypeReference $Something [[GenericType [GenericParam $t]]]]]] [[GenericParam $t]]]]")
}

func TestMultipleGenerics(t *testing.T) {
	testParse(t,
		`
type alias Tinkering t a =
    { solder : Bool
    , cool : List t
    , other : a
    }

`, `
[alias $Tinkering [wrapped-type [type-param-context [$t $a]] [record-type [[field: $solder [type-reference $Bool]] [field: $cool [type-reference $List [[local-type: [type-param $t]]]]] [field: $other [local-type: [type-param $a]]]]]]]
`)
}

func TestSkippingModule(t *testing.T) {
	testParse(t,
		`
module Main exposing (main)

import Character
import Characters.Update
`, `
[import [$Character]]
[import [$Characters $Update]]
`)
}

func TestGenericsAnnotation(t *testing.T) {
	testParse(t,
		`
type alias Tinkering t =
    { solder : Bool
    , cool : t
    }


f : Tinkering t -> Int
f a =
    2
`, `
[alias $Tinkering [wrapped-type [type-param-context [$t]] [record-type [[field: $solder [type-reference $Bool]] [field: $cool [local-type: [type-param $t]]]]]]]
[annotation: $f [wrapped-type [type-param-context [$t]] [func-type [type-reference $Tinkering [[local-type: [type-param $t]]]] -> [type-reference $Int]]]]
[definition: $f = [func ([$a]) -> #2]]
`)
}

func TestGenericsCustomType(t *testing.T) {
	testParse(t,
		`
type Maybe a =
    Nothing
    | Just a
`, `[custom-type-statement [wrapped-type [type-param-context [$a]] [custom-type $Maybe [[variant $Nothing] [variant $Just [[local-type: [type-param $a]]]]]]]]
`)
}

func TestGenericsAnnotationBegin(t *testing.T) {
	testParse(t,
		`
ownCons : itemType -> List itemType -> List itemType
`, `
[annotation: $ownCons [wrapped-type [type-param-context [$itemType]] [func-type [local-type: [type-param $itemType]] -> [type-reference $List [[local-type: [type-param $itemType]]]] -> [type-reference $List [[local-type: [type-param $itemType]]]]]]]
`)
}

func TestGenericsAnnotationSpecified(t *testing.T) {
	testParse(t,
		`
type alias Tinkering t =
    { solder : Bool
    , secret : t
    }


f : Tinkering Int -> Int
f tinkering =
    tinkering.secret
`, `
[alias $Tinkering [wrapped-type [type-param-context [$t]] [record-type [[field: $solder [type-reference $Bool]] [field: $secret [local-type: [type-param $t]]]]]]]
[annotation: $f [func-type [type-reference $Tinkering [[type-reference $Int]]] -> [type-reference $Int]]]
[definition: $f = [func ([$tinkering]) -> [lookups $tinkering [$secret]]]]
`)
}

func TestGenericsAnnotation2(t *testing.T) {
	testParse(t,
		`
f : Tinkering Int -> Int
f tinkering =
    tinkering.secret
`, `
[annotation: $f [func-type [type-reference $Tinkering [[type-reference $Int]]] -> [type-reference $Int]]]
[definition: $f = [func ([$tinkering]) -> [lookups $tinkering [$secret]]]]
`)
}

func TestGenericsStructInstantiate(t *testing.T) {
	testParse(t,
		`
type alias Tinkering t =
    { solder : Bool
    , secret : t
    }
`, "[alias $Tinkering [wrapped-type [type-param-context [$t]] [record-type [[field: $solder [type-reference $Bool]] [field: $secret [local-type: [type-param $t]]]]]]]")
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
[case: $a of [casecons $Nothing ([]) => [ccall $None]];[casecons $Just ([$oldGamepad]) => #2]]
`)
}

func TestRecordAnnotation(t *testing.T) {
	testParse(t,
		`
f : { solder : Bool, cool : Int } -> Int
f a =
    2
`, `
[annotation: $f [func-type [record-type [[field: $solder [type-reference $Bool]] [field: $cool [type-reference $Int]]]] -> [type-reference $Int]]]
[definition: $f = [func ([$a]) -> #2]]
`)
}

func TestRecordAnnotation2(t *testing.T) {
	testParse(t,
		`
type alias Sprite =
    { solder : Bool
    , cool : Int
    }


f : (a: Sprite) -> Int =
    2


f : (a : { solder : Bool, cool : Int }) -> Int =
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
[FnDef $f = [Fn ([[Arg $a: [TypeReference $Int]]]) => [RecordType [[Field: $solder [TypeReference $Bool]] [Field: $cool [TypeReference $Maybe [[TypeReference $Int]]]]] []] = #2]]
`)
}

func TestPipeRight2(t *testing.T) {
	testParse(t, `
tester : (b: String) -> Bool =
    first (2 + 2) |> second b |> third
`, `
[FnDef $tester = [Fn ([[Arg $b: [TypeReference $String]]]) => [TypeReference $Bool] = (([Call $first [(#2 + #2)]] |> [Call $second [$b]]) |> $third)]]
`)
}
