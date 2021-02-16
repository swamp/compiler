/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate

import (
	"testing"

	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func TestResourceName(t *testing.T) {
	testGenerate(t,
		`
a : Bool -> List ResourceName
a x =
    [ @this/is/cool, @as/well/as/that ]
`, `
func [function a 3 1 [[constant1 @this/is/cool #4] [constant2 @as/well/as/that #5]]]
00: lr 2,4
03: lr 3,5
06: crl 0 [2 3]
0b: ret
`)
}

func TestCharacter(t *testing.T) {
	testGenerate(t,
		`
a : Bool -> Char
a x =
    'A'
`, `
func [function a 2 1 [[constant1 int:65 #2]]]
00: lr 0,2
03: ret
`)
}

func TestTypeIdRef(t *testing.T) {
	testGenerate(t,
		`
__externalfn someLoad 1


load : TypeRef a -> a
load ignore =
    __asm callexternal 00 someLoad 01


main : Bool -> List Int
main x =
    load $(List Int)
`, `
func [function load 1 1 [[constant1 funcExternal:someLoad #2]]]
00: callexternal 0 2 ([1])
05: ret

func [function main 5 1 [[constant1 int:4 #3] [constant2 func:load #4]]]
00: lr 2,3
03: call 0 4 ([2])
08: ret
`)
}

func TestBooleanOperator(t *testing.T) {
	testGenerate(t,
		`
main : Int -> Bool
main x =
    let
        a = True
        b = False
    in
    a && b
`, `
func [function load 1 1 [[constant1 funcExternal:someLoad #2]]]
00: callexternal 0 2 ([1])
05: ret

func [function main 5 1 [[constant1 int:4 #3] [constant2 func:load #4]]]
00: lr 2,3
03: call 0 4 ([2])
08: ret
`)
}

func TestRecordTypeGenerics(t *testing.T) {
	testGenerate(t,
		`
type alias MyTypeRef a =
    { ignore : a
    }


__externalfn someLoad 1


load : MyTypeRef a -> a
load ignore =
    __asm callexternal 00 someLoad 01


main : MyTypeRef Int -> Int
main x =
    load x
`, `
func [function load 1 1 [[constant1 funcExternal:someLoad #2]]]
00: callexternal 0 2 ([1])
05: ret

func [function main 3 1 [[constant1 func:load #2]]]
00: call 0 2 ([1])
05: ret
`)
}

func TestGuard(t *testing.T) {
	testGenerate(t,
		`
isBeautiful : String -> Bool
isBeautiful name =
    | name == "Rebecca" -> True
    | _ -> False
`, `
func [function isBeautiful 2 1 [[constant1 False #3] [constant2 Rebecca #4] [constant3 True #5]]]
00: cpeq 2,1,4
04: brne 2 [label @0c]
07: lr 0,5
0a: jmp [label @0f]
0c: lr 0,3
0f: ret
`)
}

func TestStringEqual(t *testing.T) {
	testGenerate(t,
		`
isBeautiful : String -> Bool
isBeautiful name =
    if name == "Rebecca" then True else False
`, `
func [function isBeautiful 2 1 [[constant1 Rebecca #3] [constant2 True #4] [constant3 False #5]]]
00: cpve 2,1,3
04: brfa 2 [label @0c]
07: lr 0,4
0a: jmp [label @0f]
0c: lr 0,5
0f: ret
`)
}

func TestIntEqual(t *testing.T) {
	testGenerate(t,
		`
isCold : Int -> Bool
isCold temp =
    temp == -1
`, `
func [function isCold 2 1 [[constant1 int:-1 #2]]]
00: cpeq 0,1,2
04: ret
`)
}

func TestUnary2(t *testing.T) {
	testGenerate(t,
		`
someTest : Bool -> Bool -> Bool
someTest a b =
    !a && b
`, `
func [function someTest 1 2 []]
00: not 0,1
03: brfa 0 [label @09]
06: lr 0,2
09: ret
`)
}

func TestGuardMultiple(t *testing.T) {
	testGenerate(t,
		`
howGreat : Int -> String
howGreat value =
    | value == 0 -> "Zero"
    | value > 100 || value == 40 -> "Crisis"
    | _ -> "No idea"
`, `
func [function howGreat 2 1 [[constant1 No idea #4] [constant2 int:0 #5] [constant3 Zero #6] [constant4 int:100 #7] [constant5 int:40 #8] [constant6 Crisis #9]]]
00: cpeq 2,1,5
04: brne 2 [label @0c]
07: lr 0,6
0a: jmp [label @22]
0c: cpg 3,1,7
10: brne 3 [label @17]
13: cpeq 3,1,8
17: brne 3 [label @1f]
1a: lr 0,9
1d: jmp [label @22]
1f: lr 0,4
22: ret
`)
}

func TestListLiteral(t *testing.T) {
	testGenerate(t,
		`
type alias Cool =
    { name : String
    }


a : Bool -> List Cool
a x =
    [ { name = "hi" }, { name = "another" }, { name = "tjoho" } ]
`, `
func [function a func(Bool -> List<Cool>) 1 [[constant1 hi #5] [constant2 another #6] [constant3 tjoho #7]]]
00: crs 2 [5]
04: crs 3 [6]
08: crs 4 [7]
0c: crl 0 [2 3 4]
12: ret
`)
}

func TestListLiteral2(t *testing.T) {
	testGenerate(t,
		`
type alias Cool =
    { name : String
    }


a : Bool -> List Int
a x =
    [ 2, 44, 55 ]
`, `
func [function a func(Bool -> List<Int>) 1 [[constant1 int:2 #2] [constant2 int:44 #3] [constant3 int:55 #4]]]
00: crl 0 [2 3 4]
06: ret`)
}

func TestListHead(t *testing.T) {
	testGenerate(t,
		`
a : Bool -> Maybe Int
a x =
    List.head [ 2, 44, 55 ]
`, `
func [function a func(Bool -> Maybe<Int>) 1 [[constant1 int:2 #3] [constant2 int:44 #4] [constant3 int:55 #5] [constant4 func:List.head #6]]]
00: crl 2 [3 4 5]
06: call 0 6 ([2])
0b: ret
`)
}

func TestBlobInRecordAlias(t *testing.T) {
	testGenerate(t,
		`
type alias Cool =
    { name : Blob
    }


a : Cool -> Int
a blob =
    2
`, `
func [function a func(Cool -> Int) 1 [[constant1 int:2 #2]]]
00: lr 0,2
03: ret
`)
}

func TestEmptyListWithSpace(t *testing.T) {
	testGenerate(t,
		`
type alias Cool =
    { name : String
    }


a : Bool -> List Int
a x =
    []
`, `
func [function a func(Bool -> List<Int>) 1 []]
00: crl 0 []
03: ret
`)
}

func TestEmptyList(t *testing.T) {
	testGenerate(t,
		`
type alias Cool =
    { name : String
    }


a : Bool -> List Int
a x =
    []
`, `
func [function a func(Bool -> List<Int>) 1 []]
00: crl 0 []
03: ret
`)
}

func TestCustomTypeConstructor(t *testing.T) {
	testGenerate(t,
		`
type SomeEnum =
    First String
    | Anon
    | Second Int


a : Bool -> SomeEnum
a dummy =
    First "Hello"
`, `
func [function a func(Bool -> SomeEnum) 1 [[constant1 Hello #2]]]
00: createenum 0 0 ([2])
05: ret
`)
}

func TestCustomTypeVariantConstructorSecond(t *testing.T) {
	testGenerate(t,
		`
type SomeEnum =
    First String
    | Anon
    | Second Int


a : Bool -> SomeEnum
a dummy =
    Second 1
`, `
func [function a func(Bool -> SomeEnum) 1 [[constant1 int:1 #2]]]
00: createenum 0 2 ([2])
05: ret
`)
}

func TestCustomTypeGenerics(t *testing.T) {
	testGenerate(t,
		`
main : Bool -> Maybe Int
main ignored =
    Just 3
`, `
func [function main func(Bool -> Maybe<Int>) 1 [[constant1 int:3 #2]]]
00: createenum 0 1 ([2])
05: ret
`)
}

func TestCasePatternMatching(t *testing.T) {
	testGenerate(t,
		`
some : String -> Int
some a =
    case a of
        "hello" -> 0

        _ -> -1
	`, `
func [function some 2 1 [[constant1 hello #2] [constant2 int:0 #3] [constant3 int:-1 #4]]]
00: csep 0 1 [[2 [label @08]] [0 [label offset @0d]]]
08: lr 0,3
0b: jmp [label @10]
0d: lr 0,4
10: ret
`)
}

func TestSimple(t *testing.T) {
	testGenerate(t,
		`
main : Bool -> Maybe String
main ignored =
    Just "hi"
`, `
func [function main func(Bool -> Maybe<String>) 1 [[constant1 hi #2]]]
00: createenum 0 1 ([2])
05: ret
`)
}

func TestToFixed(t *testing.T) {
	testGenerate(t,
		`
main : Int -> Fixed
main a =
    Int.toFixed a
`, `
func [function main func(Int -> Fixed) 1 [[constant1 func:Int.toFixed #2]]]
00: call 0 2 ([1])
05: ret
`)
}

func TestRound(t *testing.T) {
	testGenerate(t,
		`
main : Fixed -> Int
main a =
    Int.round a
`, `
func [function main func(Fixed -> Int) 1 [[constant1 func:Int.round #2]]]
00: call 0 2 ([1])
05: ret
`)
}

func TestFixedMul(t *testing.T) {
	testGenerate(t,
		`
main : Int -> Fixed
main a =
    ( Int.toFixed 23 ) * 4.2
`, `
func [function main func(Int -> Fixed) 1 [[constant1 int:23 #4] [constant2 func:Int.toFixed #5] [constant3 int:420 #6]]]
00: call 2 5 ([4])
05: lr 3,6
08: fxmul 0,2,3
0c: ret
`)
}

func TestFixedDiv(t *testing.T) {
	testGenerate(t,
		`
main : Int -> Fixed
main a =
    -4.98 / 4.2
`, `
func [function main func(Int -> Fixed) 1 [[constant1 int:-498 #4] [constant2 int:420 #5]]]
00: lr 2,4
03: lr 3,5
06: fxdiv 0,2,3
0a: ret
`)
}

func TestCustomTypeGenericsFail(t *testing.T) {
	testGenerateFail(t,
		`
main : Bool -> Maybe Int
main ignored =
    Just "hi"
`, &decorated.UnMatchingFunctionReturnTypesInFunctionValue{})
}

func TestListWithFunc(t *testing.T) {
	testGenerate(t,
		`
type alias Cool =
    { name : String
    }


another : String -> Int
another name =
    22


a : Bool -> List Int
a x =
    [ another "tjoho", 4 ]
`, `
func [function a func(Bool -> List<Int>) 1 [[constant1 tjoho #3] [constant2 func:another #4] [constant3 int:4 #5]]]
00: call 2 4 ([3])
05: crl 0 [2 5]
0a: ret

func [function another func(String -> Int) 1 [[constant1 int:22 #2]]]
00: lr 0,2
03: ret

`)
}

func TestListType(t *testing.T) {
	testGenerate(t,
		`
type alias Position =
    { x : Int
    , y : Int
    }


type alias Sprite =
    { pos : Position
    }


move : Position -> Position -> Position
move pos delta =
    let
        newX = pos.x + delta.x

        newY = pos.y - delta.y
    in
    { x = newX, y = newY }
`, `
func [function move func(Position -> Position -> Position) 2 []]
00: get 4, 1, [#0]
05: get 5, 2, [#0]
0a: add 3,4,5
0e: get 5, 1, [#1]
13: get 6, 2, [#1]
18: sub 4,5,6
1c: crs 0 [3 4]
21: ret
`)
}

func TestListTypeWithoutName(t *testing.T) {
	testGenerate(t,
		`
type alias Position =
    { x : Int
    , y : Int
    }


type alias Sprite =
    { pos : Position
    }


move : Position -> Position -> Position
move pos delta =
    let
        newX = pos.x + delta.x

        newY = pos.y - delta.y
    in
    { x = newX, y = newY }
`, `
func [function move func(Position -> Position -> Position) 2 []]
00: get 4, 1, [#0]
05: get 5, 2, [#0]
0a: add 3,4,5
0e: get 5, 1, [#1]
13: get 6, 2, [#1]
18: sub 4,5,6
1c: crs 0 [3 4]
21: ret
`)
}

func TestSprite(t *testing.T) {
	testGenerate(t,
		`
type alias Sprite =
    { x : Int
    }


calc : Int -> Int
calc a =
    a * a + 2


main : Bool -> Sprite
main a =
    { x = calc 10 }
`, `
func [function calc func(Int -> Int) 1 [[constant1 int:2 #3]]]
00: mul 2,1,1
04: add 0,2,3
08: ret

func [function main func(Bool -> Sprite) 1 [[constant1 int:10 #3] [constant2 func:calc #4]]]
00: call 2 4 ([3])
05: crs 0 [2]
09: ret
`)
}

func TestInfoSprite(t *testing.T) {
	testGenerate(t,
		`
type alias Sprite =
    { x : Int
    }


calc : Int -> Int
calc a =
    a * a + 2


something : Sprite -> Sprite
something s =
    s


main : Bool -> Sprite
main a =
    something { x = calc 10 }
`, `
func [function calc func(Int -> Int) 1 [[constant1 int:2 #3]]]
00: mul 2,1,1
04: add 0,2,3
08: ret

func [function main func(Bool -> Sprite) 1 [[constant1 int:10 #4] [constant2 func:calc #5] [constant3 func:something #6]]]
00: call 3 5 ([4])
05: crs 2 [3]
09: call 0 6 ([2])
0e: ret

func [function something func(Sprite -> Sprite) 1 []]
00: lr 0,1
03: ret

`)
}

func TestInfoSpriteSmall(t *testing.T) {
	testGenerate(t,
		`
type alias Sprite =
    { x : Int
    }


calc : Int -> Int
calc a =
    a * a + 2


something : Sprite -> Sprite
something s =
    s


main : Bool -> Sprite
main a =
    something { x = calc 10 }
`, `
func [function calc func(Int -> Int) 1 [[constant1 int:2 #3]]]
00: mul 2,1,1
04: add 0,2,3
08: ret

func [function main func(Bool -> Sprite) 1 [[constant1 int:10 #4] [constant2 func:calc #5] [constant3 func:something #6]]]
00: call 3 5 ([4])
05: crs 2 [3]
09: call 0 6 ([2])
0e: ret

func [function something func(Sprite -> Sprite) 1 []]
00: lr 0,1
03: ret

`)
}

func TestInfoSpriteSprite(t *testing.T) {
	testGenerate(t,
		`
type alias Sprite =
    { x : Int
    }


type alias Info =
    { sprite : Sprite
    }


calc : Int -> Int
calc a =
    a * a + 2


something : Info -> Sprite
something info =
    info.sprite


main : Bool -> Sprite
main a =
    something { sprite = { x = calc 10 } }
`, `
func [function calc func(Int -> Int) 1 [[constant1 int:2 #3]]]
00: mul 2,1,1
04: add 0,2,3
08: ret

func [function main func(Bool -> Sprite) 1 [[constant1 int:10 #5] [constant2 func:calc #6] [constant3 func:something #7]]]
00: call 4 6 ([5])
05: crs 3 [4]
09: crs 2 [3]
0d: call 0 7 ([2])
12: ret

func [function something func(Info -> Sprite) 1 []]
00: get 0, 1, [#0]
05: ret


`)
}

// -- just a comment
func TestListType2(t *testing.T) {
	testGenerate(t,
		`
type alias Tinkering =
    { solder : Bool
    }


type alias Studying =
    { subject : String
    }


type alias Work =
    { web : Int
    }


type Child =
    Aron Tinkering
    | Alexandra
    | Alma Studying
    | Isabelle Work


some : Child -> String
some child =
    case child of
        Aron x ->
            if True then
                "Aron"
            else
                "strange"

        Alexandra -> "Alexandris"

        _ -> "Unknown"


main : Bool -> String
main ignored =
    some ( Aron { solder = True } )
`, `
func [function main func(Bool -> String) 1 [[constant1 True #4] [constant2 func:some #5]]]
00: crs 3 [4]
04: createenum 2 0 ([3])
09: call 0 5 ([2])
0e: ret

func [function some func(Child -> String) 1 [[constant1 True #3] [constant2 Aron #4] [constant3 strange #5] [constant4 Alexandris #6] [constant5 Unknown #7]]]
00: cse 0 1 [[0 [2] [label @0e]] [1 [] [label offset @1a]] [255 [] [label offset @1e]]]
0e: brne 3 [label @16]
11: lr 0,4
14: jmp [label @19]
16: lr 0,5
19: ret
1a: lr 0,6
1d: ret
1e: lr 0,7
21: ret

`)
}

func TestListType4(t *testing.T) {
	testGenerate(t,
		`
isBeautiful : String -> Bool
isBeautiful name =
    if name == "Rebecca" then
        let
            x = "Rebecca"
        in
        True
    else
        let
            x = "3"
        in
        False


main : Bool -> Bool
main a =
    isBeautiful "John"

`, `
func [function isBeautiful func(String -> Bool) 1 [[constant1 Rebecca #4] [constant2 True #5] [constant3 3 #6] [constant4 False #7]]]
00: cpeq 2,1,4
04: brne 2 [label @0f]
07: lr 3,4
0a: lr 0,5
0d: jmp [label @15]
0f: lr 3,6
12: lr 0,7
15: ret

func [function main func(Bool -> Bool) 1 [[constant1 John #2] [constant2 func:isBeautiful #3]]]
00: call 0 3 ([2])
05: ret
`)
}

func TestCurry(t *testing.T) {
	testGenerate(t,
		`
f : String -> Int -> Bool
f name score =
    if name == "Peter" then
        score * 2 > 100
    else
        score > 100


another : Int -> Bool
another score =
    let
        af = f "Peter"
    in
    af score
`, `
func [function another func(Int -> Bool) 1 [[constant1 Peter #3] [constant2 func:f #4]]]
00: curry 2 4 ([3])
05: call 0 2 ([1])
0a: ret

func [function f func(String -> Int -> Bool) 2 [[constant1 Peter #5] [constant2 int:2 #6] [constant3 int:100 #7]]]
00: cpeq 3,1,5
04: brne 3 [label @11]
07: mul 4,2,6
0b: cpg 0,4,7
0f: jmp [label @15]
11: cpg 0,2,7
15: ret

`)
}

func TestAppend(t *testing.T) {
	testGenerate(t,
		`
a : Int -> List Int
a x =
    [ 1, 3, 4 ] ++ [ 5, 6, 7, 8 ] ++ [ 9 ]
`, `
func [function a func(Int -> List<Int>) 1 [[constant1 int:1 #5] [constant2 int:3 #6] [constant3 int:4 #7] [constant4 int:5 #8] [constant5 int:6 #9] [constant6 int:7 #10] [constant7 int:8 #11] [constant8 int:9 #12]]]
00: crl 3 [5 6 7]
06: crl 4 [8 9 10 11]
0d: listappend 2,3,4
11: crl 3 [12]
15: listappend 0,2,3
19: ret
`)
}

func TestOwnAppend(t *testing.T) {
	testGenerate(t,
		`
ownAppender : List Int -> List Int -> List Int
ownAppender lista listb =
    lista ++ listb ++ [ 9 ]


main : Bool -> List Int
main a =
    ownAppender [ 1, 2 ] [ 3, 4 ]
`, `
func [function main func(Bool -> List<Int>) 1 [[constant1 int:1 #4] [constant2 int:2 #5] [constant3 int:3 #6] [constant4 int:4 #7] [constant5 func:ownAppender #8]]]
00: crl 2 [4 5]
05: crl 3 [6 7]
0a: call 0 8 ([2 3])
10: ret

func [function ownAppender func(List<Int> -> List<Int> -> List<Int>) 2 [[constant1 int:9 #5]]]
00: listappend 3,1,2
04: crl 4 [5]
08: listappend 0,3,4
0c: ret


`)
}

func TestCons(t *testing.T) {
	testGenerate(t,
		`
a : Int -> List Int
a x =
    99 :: [ 1, 3, 4 ]
`, `
func [function a func(Int -> List<Int>) 1 [[constant1 int:99 #3] [constant2 int:1 #4] [constant3 int:3 #5] [constant4 int:4 #6]]]
00: crl 2 [4 5 6]
06: conj 0 3 2
0a: ret
`)
}

func TestMath(t *testing.T) {
	testGenerate(t,
		`
math : Int -> Int
math a =
    72 + 4 * 3 + a


main : Bool -> Bool
main a =
    (math 32) > 116
`, `

func [function main func(Bool -> Bool) 1 [[constant1 int:32 #3] [constant2 func:math #4] [constant3 int:116 #5]]]
00: call 2 4 ([3])
05: cpg 0,2,5
09: ret

func [function math func(Int -> Int) 1 [[constant1 int:72 #4] [constant2 int:4 #5] [constant3 int:3 #6]]]
00: mul 3,5,6
04: add 2,4,3
08: add 0,2,1
0c: ret

`)
}

func TestListType3(t *testing.T) {
	testGenerate(t,
		`
isBeautiful : String -> Bool
isBeautiful name =
    if name == "Rebecca" then
        let
            x = "Rebecca"
        in
        True
    else
        let
            x = "3"
        in
        False


main : Bool -> Bool
main a =
    isBeautiful "Peter"
`, `
func [function isBeautiful func(String -> Bool) 1 [[constant1 Rebecca #4] [constant2 True #5] [constant3 3 #6] [constant4 False #7]]]
00: cpeq 2,1,4
04: brne 2 [label @0f]
07: lr 3,4
0a: lr 0,5
0d: jmp [label @15]
0f: lr 3,6
12: lr 0,7
15: ret

func [function main func(Bool -> Bool) 1 [[constant1 Peter #2] [constant2 func:isBeautiful #3]]]
00: call 0 3 ([2])
05: ret
`)
}

/*
func [function first func(Int -> Int) 1 []]
00: mul 0,1,1
04: ret

func [function second func(String -> Int -> Bool) 2 [[constant1 int:25 #3]]]
00: cpg 0,2,3 ; a > 25
04: ret

func [function tester func(String -> Bool) 1 [[constant1 int:2 #5] [constant2 func:first #6] [constant3 func:second #7] [constant4 func:third #8]]]
00: add 4,5,5             ; 2 + 2
04: call 3 6 ([4])        ; first (2+2) -> @3
09: call 2 7 ([1 3])      ; second b ( first (2+2) ) -> @2
0f: call 0 8 ([2])        ; return third (second b ( first ( 2 + 2 ) ) )
14: ret

func [function third func(Bool -> Bool) 1 []]
00: lr 0,1 ; return a
03: ret
*/
const expectedPipeAsm = `
func [function first 2 2 []]
00: mul 0,1,1
04: ret

func [function second 4 2 [[constant1 int:25 #3]]]
00: cpg 0,2,3
04: ret

func [function something 5 0 [[constant1 hello #1]]]
00: lr 0,1
03: ret

func [function tester 6 1 [[constant1 int:2 #6] [constant2 func:something #7] [constant3 func:first #8] [constant4 func:second #9] [constant5 func:third #10]]]
00: add 4,6,6
04: call 5 7 ([])
08: call 3 8 ([4 5])
0e: call 2 9 ([1 3])
14: call 0 10 ([2])
19: ret

func [function third 7 1 []]
00: lr 0,1
03: ret
`

func TestOperatorPipeRight(t *testing.T) {
	testGenerate(t,
		`
first : Int -> Int
first a =
    a * a


second : String -> Int -> Bool
second str a =
    a > 25


third : Bool -> Bool
third a =
    a


tester : String -> Bool
tester b =
    first (2 + 2) |> second b |> third
`, expectedPipeAsm)
}

func TestUnaryMinus(t *testing.T) {
	testGenerate(t,
		`
second : Int -> Bool
second a =
    a < 0


tester : Int -> Bool
tester b =
    second -b
`, `
func [function second 2 1 [[constant1 int:0 #2]]]
00: cpl 0,1,2
04: ret

func [function tester 2 1 [[constant1 func:second #3]]]
00: neg 2,1
03: call 0 3 ([2])
08: ret`)
}

func TestUnaryMinusInNotACall(t *testing.T) {
	testGenerate(t,
		`
type alias Data = { end : Int }


second : Int -> Bool
second a =
    a < 0


tester : Data -> Int -> Bool
tester data b =
    second (data.end - b)
`, `
func [function second 2 1 [[constant1 int:0 #2]]]
00: cpl 0,1,2
04: ret

func [function tester 5 2 [[constant1 func:second #5]]]
00: get 4, 1, [#0]
05: sub 3,4,2
09: call 0 5 ([3])
0e: ret
`)
}

func TestGuardLet2(t *testing.T) {
	testGenerate(t,
		`
tester : Int -> Char
tester x =
    let
        existingTile = 'a'
        isUpperLeft = False
    in
    | existingTile == '_' -> '@'
    | isUpperLeft -> '/'
    | _ -> '2'
`, `
func [function tester 2 1 [[constant1 int:97 #5] [constant2 False #6] [constant3 int:50 #7] [constant4 int:95 #8] [constant5 int:64 #9] [constant6 int:47 #10]]]
00: lr 2,5
03: lr 3,6
06: cpve 4,2,8
0a: brfa 4 [label @12]
0d: lr 0,9
10: jmp [label @1d]
12: brfa 3 [label @1a]
15: lr 0,10
18: jmp [label @1d]
1a: lr 0,7
1d: ret
`)
}

func TestUnaryMinusInNotACall2(t *testing.T) {
	testGenerate(t,
		`
second : Int -> Int -> Bool
second time startTime =
    let
        timeInAnimation = time - startTime
    in
    timeInAnimation > 10


tester : Int -> Bool
tester b =
    second -2 -3
`, `
func [function second 2 2 [[constant1 int:10 #4]]]
00: sub 3,1,2
04: cpg 0,3,4
08: ret

func [function tester 3 1 [[constant1 int:2 #4] [constant2 int:3 #5] [constant3 func:second #6]]]
00: neg 2,4
03: neg 3,5
06: call 0 6 ([2 3])
0c: ret
`)
}

func TestDivide(t *testing.T) {
	testGenerate(t,
		`
type alias State = { time : Int }


second : State -> Bool
second state =
    let
        ft = (Int.toFixed state.time) / 50.0
    in
    ft > 10.0


tester : Int -> Bool
tester b =
    second { time = 42 }
`, `
func [function second 4 1 [[constant1 func:Int.toFixed #6] [constant2 int:50000 #7] [constant3 int:10000 #8]]]
00: get 4, 1, [#0]
05: call 3 6 ([4])
0a: lr 5,7
0d: fxdiv 2,3,5
11: lr 3,8
14: add 0,2,3
18: ret

func [function tester 5 1 [[constant1 int:42 #3] [constant2 func:second #4]]]
00: crs 2 [3]
04: call 0 4 ([2])
09: ret

`)
}

func TestUnaryNot(t *testing.T) {
	testGenerate(t,
		`
second : Bool -> Int
second a =
    if a then 3 else 4


tester : Bool -> Int
tester b =
    second !b
`, `
func [function second 2 1 [[constant1 int:3 #2] [constant2 int:4 #3]]]
00: brfa 1 [label @08]
03: lr 0,2
06: jmp [label @0b]
08: lr 0,3
0b: ret

func [function tester 2 1 [[constant1 func:second #3]]]
00: not 2,1
03: call 0 3 ([2])
08: ret
`)
}

func TestBinaryMinus(t *testing.T) {
	testGenerate(t,
		`
second : Int -> Bool
second a =
    a < 0


tester : Int -> Bool
tester b =
    let
        x = 42 - 2
    in
    second (2 - b)

`, `
func [function second 2 1 [[constant1 int:0 #2]]]
00: cpl 0,1,2
04: ret

func [function tester 2 1 [[constant1 int:42 #4] [constant2 int:2 #5] [constant3 func:second #6]]]
00: sub 2,4,5
04: sub 3,5,1
08: call 0 6 ([3])
0d: ret
`)
}

func TestOperatorPipeLeft(t *testing.T) {
	testGenerate(t,
		`
first : Int -> Int
first a =
    a * a


second : String -> Int -> Bool
second str a =
    a > 25


third : Bool -> Bool
third a =
    a


tester : String -> Bool
tester b =
    third <| second b <| first (2 + 2)
`, expectedPipeAsm)
}

func TestOperatorPipeLeft4(t *testing.T) {
	testGenerate(t,
		`
something : String
something =
    "hello"


first : Int -> String -> Int
first a b =
    a * a


second : String -> Int -> Bool
second str a =
    a > 25


third : Bool -> Bool
third a =
    a


tester : String -> Bool
tester b =
    third <| second b <| first (2 + 2) <| something()
`, expectedPipeAsm)
}

func TestOperatorPipeLeftConstructor(t *testing.T) {
	testGenerate(t,
		`
type alias Vector2 = { x : Int, y : Int }


first : Int -> Vector2 -> Int
first a vec =
    42


second : String -> Int -> Bool
second str a =
    a > 25


third : Bool -> Bool
third a =
    a


tester : String -> Bool
tester b =
    third <| second b <| first (2 + 2) <| Vector2 0 23
`, `
func [function first 3 2 [[constant1 int:42 #3]]]
00: lr 0,3
03: ret

func [function second 6 2 [[constant1 int:25 #3]]]
00: cpg 0,2,3
04: ret

func [function tester 7 1 [[constant1 int:2 #6] [constant2 int:0 #7] [constant3 int:23 #8] [constant4 func:first #9] [constant5 func:second #10] [constant6 func:third #11]]]
00: add 4,6,6
04: crs 5 [7 8]
09: call 3 9 ([4 5])
0f: call 2 10 ([1 3])
15: call 0 11 ([2])
1a: ret

func [function third 8 1 []]
00: lr 0,1
03: ret
`)
}

func TestListType5(t *testing.T) {
	testGenerate(t,
		`
tester : String -> Bool
tester name =
    name == "Peter" && name == "Rebecca" && name == "Alma"
`, `
func [function tester func(String -> Bool) 1 [[constant1 Peter #2] [constant2 Rebecca #3] [constant3 Alma #4]]]
00: cpeq 0,1,2 ; name == "Peter"
04: brne 0 [label @0b] ; if name != "Peter"
07: cpeq 0,1,3 				 ; name is "Peter": compare name == "Rebecca"
0b: brne 0 [label @12] ; name != "Peter"
0e: cpeq 0,1,4			 	 ; name == "Alma"?
12: ret
`)
}

func TestUpdateSprite(t *testing.T) {
	testGenerate(t,
		`
type alias Scale2 =
    { scaleX : Int
    , scaleY : Int
    }


type alias Sprite =
    { dummy : Int
    , scale : Scale2
    }


updateSprite : Sprite -> Int -> Sprite
updateSprite inSprite newScale =
    { inSprite | scale = { scaleX = newScale, scaleY = newScale } }
`, `
func [function updateSprite func(Sprite -> Int -> Sprite) 2 []]
00: crs 3 [2 2]
05: update 0 1 [#1<-3]
0b: ret
`)
}

func TestSpriteListOrAny(t *testing.T) {
	testGenerate(t,
		`
type alias Scale2 =
    { scaleX : Int
    , scaleY : Int
    }


type alias Sprite =
    { dummy : Int
    , scale : Scale2
    }


spritesToDraw : Bool -> List Sprite
spritesToDraw doIt =
    if doIt then
        [ { dummy = 0, scale = { scaleX = 10, scaleY = 10 } } ]
    else
        []


main : Bool -> List Sprite
main ignore =
    spritesToDraw True

`, `
func [function main func(Bool -> List<Sprite>) 1 [[constant1 True #2] [constant2 func:spritesToDraw #3]]]
00: call 0 3 ([2])
05: ret

func [function spritesToDraw func(Bool -> List<Sprite>) 1 [[constant1 int:0 #4] [constant2 int:10 #5]]]
00: brne 1 [label @13]
03: crs 3 [5 5]
08: crs 2 [4 3]
0d: crl 0 [2]
11: jmp [label @16]
13: crl 0 []
16: ret

`)
}

func TestArray(t *testing.T) {
	testGenerate(t,
		`
type alias Scale2 =
    { scaleX : Int
    , scaleY : Int
    }


type alias Sprite =
    { dummy : Int
    , scale : Scale2
    }


spritesToDraw : Bool -> Array Sprite
spritesToDraw doIt =
    Array.fromList [ { dummy = 0, scale = { scaleX = 10, scaleY = 10 } } ]


main : Bool -> Sprite
main ignore =
    let
        mightBeSprite = Array.get 0 (spritesToDraw True)
    in
    case mightBeSprite of
        Just sprite -> sprite

        Nothing -> { dummy = 0, scale = { scaleX = 10, scaleY = 10 } }
`, `
func [function main func(Bool -> Sprite) 1 [[constant1 int:0 #5] [constant2 True #6] [constant3 func:spritesToDraw #7] [constant4 func:Array.get #8] [constant5 int:10 #9]]]
00: call 3 7 ([6])
05: call 2 8 ([5 3])
0b: cse 0 2 [[1 [4] [label @16]] [0 [] [label offset @1a]]]
16: lr 0,4
19: ret
1a: crs 4 [9 9]
1f: crs 0 [5 4]
24: ret

func [function spritesToDraw func(Bool -> Array<Sprite>) 1 [[constant1 int:0 #5] [constant2 int:10 #6] [constant3 func:Array.fromList #7]]]
00: crs 4 [6 6]
05: crs 3 [5 4]
0a: crl 2 [3]
0e: call 0 7 ([2])
13: ret

`)
}

func TestCustomTypeReturn(t *testing.T) {
	testGenerate(t,
		`
type PlayerAction =
    None
    | Jump


type alias Gamepad =
    { a : Int
    }


checkGamepad : Gamepad -> Gamepad -> PlayerAction
checkGamepad oldGamepad gamepad =
    if (gamepad.a != 0 && (oldGamepad.a != gamepad.a)) then
        Jump
    else
        None


checkGamepadMaybe : Maybe Gamepad -> Maybe Gamepad -> PlayerAction
checkGamepadMaybe oldGamepad gamepad =
    case oldGamepad of
        Nothing -> None

        Just a -> case gamepad of
            Nothing -> None

            Just b -> checkGamepad a b
`,
		`func [function checkGamepad func(Gamepad -> Gamepad -> PlayerAction) 2 [[constant1 int:0 #6]]]
00: get 4, 2, [#0]
05: cpne 3,4,6
09: brne 3 [label @1a]
0c: get 4, 1, [#0]
11: get 5, 2, [#0]
16: cpne 3,4,5
1a: brne 3 [label @23]
1d: createenum 0 1 ([])
21: jmp [label @27]
23: createenum 0 0 ([])
27: ret

func [function checkGamepadMaybe func(Maybe<Gamepad> -> Maybe<Gamepad> -> PlayerAction) 2 [[constant1 func:checkGamepad #5]]]
00: cse 0 1 [[0 [] [label @0b]] [1 [3] [label offset @10]]]
0b: createenum 0 0 ([])
0f: ret
10: cse 0 2 [[0 [] [label @1b]] [1 [4] [label offset @20]]]
1b: createenum 0 0 ([])
1f: ret
20: call 0 5 ([3 4])
26: ret
`)
}

//-- ignore this

func TestBadToken(t *testing.T) {
	testGenerate(t,

		`
type alias Gamepad =
    { a : Int
    }


type PlayerAction =
    None
    | Jump


type alias Touch =
    { dummy : Int
    }


type alias UserInput =
    { gamepads : Array Gamepad
    , touches : List Touch
    }


type alias ScreenInfo =
    { dummy : Int
    }


type alias PlayerInputs =
    { inputs : List PlayerAction
    }


checkGamepad : Gamepad -> Gamepad -> PlayerAction
checkGamepad oldGamepad gamepad =
    if gamepad.a != 0 && (oldGamepad.a != gamepad.a) then
        Jump
    else
        None


checkMaybeGamepad : Maybe Gamepad -> Maybe Gamepad -> PlayerAction
checkMaybeGamepad a b =
    case a of
        Just actual -> Jump

        Nothing -> case b of
            Just other -> None

            Nothing -> None


main : UserInput -> UserInput -> ScreenInfo -> PlayerInputs
main oldUserInputs userInputs screenInfo =
    { inputs = [] }
`,
		`
func [function checkGamepad func(Gamepad -> Gamepad -> PlayerAction) 2 [[constant1 int:0 #6]]]
00: get 4, 2, [#0]
05: cpne 3,4,6
09: brne 3 [label @1a]
0c: get 4, 1, [#0]
11: get 5, 2, [#0]
16: cpne 3,4,5
1a: brne 3 [label @23]
1d: createenum 0 1 ([])
21: jmp [label @27]
23: createenum 0 0 ([])
27: ret

func [function checkMaybeGamepad func(Maybe<Gamepad> -> Maybe<Gamepad> -> PlayerAction) 2 []]
00: cse 0 1 [[1 [3] [label @0b]] [0 [] [label offset @10]]]
0b: createenum 0 1 ([])
0f: ret
10: cse 0 2 [[1 [3] [label @1b]] [0 [] [label offset @20]]]
1b: createenum 0 0 ([])
1f: ret
20: createenum 0 0 ([])
24: ret

func [function main func(UserInput -> UserInput -> ScreenInfo -> PlayerInputs) 3 []]
00: crl 4 []
03: crs 0 [4]
07: ret
`)
}

func TestBadTokenMinimal(t *testing.T) {
	testGenerate(t,

		`
checkMaybe : Maybe Int -> Maybe Int -> Int
checkMaybe a b =
    case a of
        Just v -> 0

        Nothing ->
            case b of
                Just v -> v

                Nothing -> 1


main : Int -> Int
main a =
    2
`,
		`
func [function checkMaybe func(Maybe<Int> -> Maybe<Int> -> Int) 2 [[constant1 int:0 #4] [constant2 int:1 #5]]]
00: cse 0 1 [[1 [3] [label @0b]] [0 [] [label offset @0f]]]
0b: lr 0,4
0e: ret
0f: cse 0 2 [[1 [3] [label @1a]] [0 [] [label offset @1e]]]
1a: lr 0,3
1d: ret
1e: lr 0,5
21: ret

func [function main func(Int -> Int) 1 [[constant1 int:2 #2]]]
00: lr 0,2
03: ret
    `)
}

func TestBadTokenMinimal2(t *testing.T) {
	testGenerate(t,

		`
type PlayerAction =
    None
    | Jump


type alias Gamepad =
    { a : Int
    }


type alias UserInput =
    { gamepads : Array Gamepad
    }


type alias PlayerInputs =
    { inputs : List PlayerAction
    }


checkMaybeGamepad : Maybe Gamepad -> PlayerAction
checkMaybeGamepad a =
    None


main : UserInput -> PlayerInputs
main oldUserInputs =
    let
        gamepads = oldUserInputs.gamepads

        maybeOld = Array.get 0 gamepads

        playerAction = checkMaybeGamepad maybeOld
    in
    { inputs = [ playerAction ] }
    `,
		`
func [function checkMaybeGamepad func(Maybe<Gamepad> -> PlayerAction) 1 []]
00: createenum 0 0 ([])
04: ret

func [function main func(UserInput -> PlayerInputs) 1 [[constant1 int:0 #6] [constant2 func:Array.get #7] [constant3 func:checkMaybeGamepad #8]]]
00: get 2, 1, [#0]
05: call 3 7 ([6 2])
0b: call 4 8 ([3])
10: crl 5 [4]
14: crs 0 [5]
18: ret
    `)
}
