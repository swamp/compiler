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
    [| { name = "hi" }, { name = "another" }, { name = "tjoho" } |]
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
