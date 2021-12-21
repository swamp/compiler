package generate_sp

import (
	"testing"
)

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

func TestBooleanOperator(t *testing.T) {
	testGenerate(t,
		`
main : Int -> Bool
main _ =
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

func TestArrayLiteral(t *testing.T) {
	testGenerate(t,
		`
type alias Cool =
    { name : String
    }


a : Bool -> Array Cool
a x =
    [| { name = "hello" }, { name = "world" }, { name = "ossian" } |]
`, `
func [function a 5 1 [[constant1 hi #5] [constant2 another #6] [constant3 tjoho #7]]]
00: crs 2 [5]
04: crs 3 [6]
08: crs 4 [7]
0c: crs 0 [2 3 4]
12: ret
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

func TestMaybeInt(t *testing.T) {
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

func TestIfStatement(t *testing.T) {
	testGenerate(t,
		`
main : String -> Bool
main name =
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
isWinner : String -> Int -> Bool
isWinner name score =
    if name == "Ossian" then
        score * 2 > 100
    else
        score > 100


main : Int -> Bool
main score =
    let
        checkScoreFn = isWinner "Ossian"
    in
    checkScoreFn score
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

func TestGuardLetInChar(t *testing.T) {
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
`)
}

func TestCasePatternMatchingString(t *testing.T) {
	testGenerate(t,
		`
some : String -> Int
some a =
    case a of
        "hello" -> 0

        "something else" -> 1

        _ -> -1
	`, `

`)
}
