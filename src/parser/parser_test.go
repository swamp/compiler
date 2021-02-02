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



func TestConstant(t *testing.T) {
	testParse(t,
		`
fn : String
fn =
    "Hello"
`,
		`
[annotation: $fn [type-reference $String]]
[definition: $fn = [func ([]) -> 'Hello']]
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
[let: [[letassign $a = ((('' ++ [call Debug.$toString [$a]]) ++ ' ') ++ [call Debug.$toString [$b]])]] in $a]
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
[let: [[letassign $a = [guard: [{($y > #3) $y} {($y == #4) #45}] 'hello']]] in $a]
`)
}


func TestGuardExpression(t *testing.T) {
	testParse(t,
		`
fn : Int -> Int
fn x =
    | x < 2 -> 42
    | x == 4 -> 45
    | _ -> -1
`,

		`
[annotation: $fn [func-type [type-reference $Int] -> [type-reference $Int]]]
[definition: $fn = [func ([$x]) -> [guard: [{($x < #2) #42} {($x == #4) #45}] #-1]]]
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

func TestType(t *testing.T) {
	testParseExpression(t,
		`Unknown`,
		`[ccall $Unknown]`)
}

func TestString(t *testing.T) {
	testParseExpression(t,
		`"hello, world!"`,
		`'hello, world!'`)
}

func TestStringContinuation(t *testing.T) {
	testParseExpression(t,
		`"hello, world!
  .next line"`,
		`'hello, world!  .next line'`)
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
		`[if: €true then #1 else #0]`)
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
[if: €true then #1 else #0]
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
[if: €true then [if: €false then #22 else #44] else #0]
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
[if: €true then [let: [[letassign $x = #44]] in $y] else #0]
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
[if: €true then [let: [[letassign $x = #44] [letassign $y = #98]] in $y] else #0]
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
[if: €true then [let: [[letassign $x = [case: $something of [casecons $Kalle ([]) => #2];[casecons $Lisa ([$a]) => #44]]] [letassign $y = #98]] in $y] else #0]
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
		`[record-literal: [[$a = #3]]]`)
}


func TestRecordLiteralTwoItems(t *testing.T) {
	testParseExpression(t,
		`
{ a = 3, b = 4 }
`,
		`[record-literal: [[$a = #3] [$b = #4]]]`)
}

func TestRecordLiteralSeveralLines(t *testing.T) {
	testParseExpression(t,
		`
{ a = 3
, b = 4 }
`,
		`[record-literal: [[$a = #3] [$b = #4]]]`)
}

func TestListOneLineFromList(t *testing.T) {
	testParseExpression(t,
		`
Array.fromList [ Array.fromList [ 0, 1, 2, 3 ] Array.fromList [ 8, 9, 10, 11 ] ]
`,
		`[call Array.$fromList [[list-literal: [[call Array.$fromList [[list-literal: [#0 #1 #2 #3]] Array.$fromList [list-literal: [#8 #9 #10 #11]]]]]]]]`)
}

func TestContinueOnNextLine(t *testing.T) {
	testParseExpression(t,
		`
callSomething a
    |> nextLine
`,
		`[call $nextLine [[call $callSomething [$a]]]]`)
}

func TestRecordLiteralSeveralLinesHex(t *testing.T) {
	testParseExpression(t,
		`
{ a = 0x00FF00FF
, b = 4 }
`,
		`[record-literal: [[$a = #16711935] [$b = #4]]]`)
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
		`[let: [[letassign $x = [record-literal: [[$a = [record-literal: [[$scaleX = #100] [$scaleY = #200]]]] [$b = #4]]]]] in $x]`)
}

func TestListLiteralTwoItems(t *testing.T) {
	testParseExpression(t,
		`
[ 3, 4 + 5 ]
`,
		`[list-literal: [#3 (#4 + #5)]]`)
}


func TestSimpleAnnotation(t *testing.T) {
	testParse(t,
		`
a : Int
`,
		`[annotation: $a [type-reference $Int]]`)
}

func TestSimpleDefinition(t *testing.T) {
	testParse(t,
		`
a =
    3
`,
		`[definition: $a = [func ([]) -> #3]]`)
}

func TestCurrying(t *testing.T) {
	testParse(t,
		`
f : String -> Int -> Bool
f name score =
    if name == "Peter" then
        score * 2
    else
        score


another : Int -> Bool
another score =
    let
        af = f "Peter"
    in
    af score
`,
		`
[annotation: $f [func-type [type-reference $String] -> [type-reference $Int] -> [type-reference $Bool]]]
[definition: $f = [func ([$name $score]) -> [if: ($name == 'Peter') then ($score * #2) else $score]]]
[annotation: $another [func-type [type-reference $Int] -> [type-reference $Bool]]]
[definition: $another = [func ([$score]) -> [let: [[letassign $af = [call $f ['Peter']]]] in [call $af [$score]]]]]
`)

}

func TestLambda(t *testing.T) {
	testParseExpression(t,
		`
\x -> hello + 2
`,
		`[lambda ([$x]) -> ($hello + #2)]`)
}

func TestLambdaThree(t *testing.T) {
	testParseExpression(t,
		`
\x another andLast -> hello + 2
`,
		`[lambda ([$x $another $andLast]) -> ($hello + #2)]`)
}

func TestAnnotation(t *testing.T) {
	testParse(t,
		`
something : Bool -> Int
`, `
[annotation: $something [func-type [type-reference $Bool] -> [type-reference $Int]]]
`)
}

func TestAnnotationParen(t *testing.T) {
	testParse(t,
		`
something : (String -> Bool) -> Int
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
[import [$First]]
[import [$Character $Damage $Sub]]
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
[alias $Struct [record-type [[field: $a [type-reference $Int]] [field: $b [type-reference $Boolean]]] []]]
[annotation: $a [func-type [type-reference $Int] -> [type-reference $Struct]]]
[definition: $a = [func ([$ignore]) -> [ccall $Struct [#2 €false]]]]
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
[alias $Struct [record-type [[field: $a [type-reference $Int]] [field: $b [type-reference $Boolean]]] []]]
[annotation: $a [func-type [type-reference $Int] -> [type-reference $Struct]]]
[definition: $a = [func ([$ignore]) -> [ccall $Struct [[record-literal: [[$a = #2] [$b = €false]]]]]]]
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
[annotation: $hello [func-type [type-reference $Int] -> [type-reference $Int] -> [type-reference $Int]]]
[definition: $hello = [func ([$first $c]) -> (#4 + $first)]]
`)
}

func TestDefinitionWithCallWithRecords(t *testing.T) {
	testParse(t,
		`
main a =
    { sprite = { x = calc 10 } }
`,
		`
		[definition: $main = [func ([$a]) -> [record-literal: [[$sprite = [record-literal: [[$x = [call $calc [#10]]]]]]]]]]
`)
}

func TestSimpleCall(t *testing.T) {
	testParse(t,
		`
rebecca is fantastic =
    something (3 * 3) (4 + 4)
`, "[definition: $rebecca = [func ([$is $fantastic]) -> [call $something [(#3 * #3) (#4 + #4)]]]]")
}

func TestSimpleCallWithLookup(t *testing.T) {
	testParseExpression(t,
		`
move sprite.rootPosition delta
`, "[call $move [[lookups $sprite [$rootPosition]] $delta]]")
}

func TestSimpleCallWithLookupInAssignmentBlock(t *testing.T) {
	testParseExpression(t,
		`
{ rootPosition = move sprite.rootPosition delta }
`, "[record-literal: [[$rootPosition = [call $move [[lookups $sprite [$rootPosition]] $delta]]]]]")
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
`, "[if: ($x == #3) then [call $extrude [(#5 * #4)]] else #5]")
}

func TestBoolean(t *testing.T) {
	testParseExpression(t, "True", "€true")
}

func TestOneLineIfWithCallAndPipe(t *testing.T) {
	testParseExpression(t,
		`
if x == 3 then
    extrude (5 * 4) |> minimize true
else
    5
`, "[if: ($x == #3) then [call $minimize [$true [call $extrude [(#5 * #4)]]]] else #5]")
}

func TestConstructorWithOneField(t *testing.T) {
	testParseExpression(t,
		`
{ first = 42 }
`, "[record-literal: [[$first = #42]]]")

}

func TestLookup(t *testing.T) {
	testParseExpression(t,
		`
a.b.c + d.e * f.g
`, "([lookups $a [$b $c]] + ([lookups $d [$e]] * [lookups $f [$g]]))")

}

func TestConstructorWithTwoFields(t *testing.T) {
	testParseExpression(t,
		`
{ first = 42, second = 99 }
`, "[record-literal: [[$first = #42] [$second = #99]]]")

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

func TestList(t *testing.T) {
	testParseExpression(t,
		`
[ 2, 4, b, 101 ]
`, "[list-literal: [#2 #4 $b #101]]")

}

func TestModuleReferences(t *testing.T) {
	testParseExpression(t,
		`
FirstModule.SecondMod.someFunc 12
`, "[call FirstModule.SecondMod.$someFunc [#12]]")

}

func TestModuleReferenceWithType(t *testing.T) {
	testParse(t,
		`
a : Bool -> FirstModule.SecondMod.ThisIsAType
`, "[annotation: $a [func-type [type-reference $Bool] -> [type-reference FirstModule.SecondMod.$ThisIsAType]]]")

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
`, "[definition: $someFunc = [func ([$a $b]) -> [let: [[letassign $i = #3] [letassign $j = #4]] in [if: ($i >= #6) then [call $call [($i + $j)]] else #3]]]]")

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
`, "[case: $x of [casecons $_ ([]) => #2]]")
}

func TestCustomType(t *testing.T) {
	testParse(t,
		`
type SomeEnum =
    First String
    | Anon
    | Second Int

`, "[custom-type-statement [custom-type $SomeEnum [[variant $First [[type-reference $String]]] [variant $Anon] [variant $Second [[type-reference $Int]]]]]]")
}

func TestCustomTypeNewFormatting(t *testing.T) {
	testParse(t,
		`
type SomeEnum =
    First String
    | Anon
    | Second Int

`, "[custom-type-statement [custom-type $SomeEnum [[variant $First [[type-reference $String]]] [variant $Anon] [variant $Second [[type-reference $Int]]]]]]")
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
[ccall Imaginary.Module.$First ['Hello']]
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
[custom-type-statement [custom-type $SomeEnum [[variant $First [[type-reference $String]]] [variant $Anon] [variant $Second [[type-reference $Int]]]]]]
[annotation: $a [func-type [type-reference $Bool] -> [type-reference $SomeEnum]]]
[definition: $a = [func ([$dummy]) -> [ccall Imaginary.Module.$First ['Hello']]]]
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
`, "[call $third [((#2 * #3) + (#5 * #66)) [call $anotherCall [$anotherParam1 [call $call [$param1 $param2]]]]]]")
}

func TestPipeBackward(t *testing.T) { // --- FIXME
	testParseExpression(t,
		`
call param1 param2 <| laterCall this <| third (2 * 3 + 5 * 66)
`, "[call $call [$param1 $param2 [call $laterCall [$this [call $third [((#2 * #3) + (#5 * #66))]]]]]]")
}

func TestPipeRight(t *testing.T) {
	testParseExpression(t,
		`
first (2 + 2) |> second
`, "[call $second [[call $first [(#2 + #2)]]]]")
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

`, "[alias $Tinkering [wrapped-type [type-param-context [$t]] [record-type [[field: $solder [type-reference $Bool]] [field: $cool [type-reference $Something [[local-type: [type-param $t]]]]]]]]]")
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
[case: $a of [casecons $Nothing ([]) => [ccall $None]];[casecons $Just ([$oldGamepad]) => #2]]
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


f : Sprite -> Int
f a =
    2


f : { solder : Bool, cool : Int } -> Int
k a =
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
f : Int -> { solder : Bool, cool : Maybe Int }
f a =
    2
`, `
[annotation: $f [func-type [type-reference $Int] -> [record-type [[field: $solder [type-reference $Bool]] [field: $cool [type-reference $Maybe [[type-reference $Int]]]]]]]]
[definition: $f = [func ([$a]) -> #2]]
`)
}


func TestPipeRight2(t *testing.T) {
	testParse(t, `
tester : String -> Bool
tester b =
    first (2 + 2) |> second b |> third
`, `
[annotation: $tester [func-type [type-reference $String] -> [type-reference $Bool]]]
[definition: $tester = [func ([$b]) -> [call $third [[call $second [$b [call $first [(#2 + #2)]]]]]]]]
`)
}
