/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package deccy

import (
	"testing"

	decorated "github.com/swamp/compiler/src/decorated/expression"
	parerr "github.com/swamp/compiler/src/parser/errors"
)

func TestComment(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
--| ignore this
first : Int -> Int
first a =
    a * a
`,
		`
[ModuleDef $first = [FunctionValue ([[arg $a = [PrimitiveTypeVariantRef NamedDefTypeRef [Primitive Int]]]]) -> (arithmetic [FunctionParamRef $a [arg $a = [PrimitiveTypeVariantRef NamedDefTypeRef [Primitive Int]]]] MULTIPLY [FunctionParamRef $a [arg $a = [PrimitiveTypeVariantRef NamedDefTypeRef [Primitive Int]]]])]]
`)
}

func TestConstant(t *testing.T) {
	testDecorateWithoutDefault(t, `
fn =
    "Hello"


another : String
another =
    fn


`, `
[ModuleDef $another = [FunctionValue ([]) -> [ConstantRef [NamedDefinitionReference :$fn]]]]
[ModuleDef $fn = [Constant [str Hello]]]
`)
}

func TestBooleanLookup(t *testing.T) {
	testDecorateWithoutDefault(t, `
another : Bool
another =
    let
        a = True
        b = False
    in
    a && b
`, `
[ModuleDef $another = [FunctionValue ([]) -> [Let [[LetAssign [[LetVar $a]] = [Bool true]] [LetAssign [[LetVar $b]] = [Bool false]]] in [Logical [LetVarRef [LetVar $a]] and [LetVarRef [LetVar $b]]]]]]`)
}

func TestCast(t *testing.T) {
	testDecorate(t, `
type alias Something = Int


another : Bool
another =
    let
        b = 32 : Something
    in
    b >= 32
`, `
[ModuleDef $another = [FunctionValue ([]) -> [Let [[LetAssign [[LetVar $b]] = [Cast [Integer 32] [AliasRef [LetVarRef [Alias Something [PrimitiveTypeVariantRef NamedDefTypeRef [Primitive Int]]]]]]]] in [BoolOp [LetVarRef [LetVar $b]] GRE [Integer 32]]]]]
`)
}

func TestResourceName(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
--| ignore this
first : ResourceName -> Int
first _ =
    2


main : Bool -> Int
main _ =
    first @this/is/something.txt
`,
		`
func(ResourceName -> Int) : [func  [primitive ResourceName] [primitive Int]]
func(Bool -> Int) : [func  [primitive Bool] [primitive Int]]

first = [functionvalue ([[arg $a = [primitive ResourceName]]]) -> [integer 2]]
main = [functionvalue ([[arg $ignore = [primitive Bool]]]) -> [fcall [getvar $first [primitive Int]] [[resource name this/is/something.txt]]]]
`)
}

func TestAnyMatchingType(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
fn : * -> Int
fn _ _ _ =
    23
`,
		`
[ModuleDef $fn = [FunctionValue ([[arg $_ = [AnyMatching AnyMatchingType]]]) -> [Integer 23]]]
`)
}

func TestAnyMatchingTypeCall(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
fn : * -> Int
fn _ _ _ =
    23


main : Int
main =
    fn 23 "hello" 42.0
`,
		`
[mdefx $fn = [functionvalue ([[arg $a = [anymatching types AnyMatchingType]]]) -> [integer 23]]]
[mdefx $main = [functionvalue ([]) -> [fcall [functionref named definition reference [functionvalue ([[arg $a = [anymatching types AnyMatchingType]]]) -> [integer 23]]] [[integer 23] [str hello] [integer 42000]]]]]
`)
}

func TestAnyMatchingTypeCallMiddle(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
fn : String -> * -> Int
fn _ _ _ =
    23


main : Int
main =
    fn "hello" 42.0
`,
		`
[ModuleDef $fn = [FunctionValue ([[arg $_ = [PrimitiveTypeVariantRef [NamedDefTypeRef <nil>:[type-reference $String]] [Primitive String]]] [arg $_ = [AnyMatching AnyMatchingType]]]) -> [Integer 23]]]
[ModuleDef $main = [FunctionValue ([]) -> [FnCall [FunctionRef [NamedDefinitionReference :$fn]] [[String hello] [Fixed 42000]]]]]
`)
}

func TestAnyMatchingTypeCallMiddleLocalType(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
fn : String -> * -> a -> List a
fn _ _ _ =
    [ 23.0 ]


main : List Fixed
main =
    fn "hello" -23939 42.0
`,
		`
[ModuleDef $fn = [FunctionValue ([[Arg $_ : [PrimitiveTypeVariantRef [NamedDefTypeRef :[TypeReference $String]] [Primitive String]]] [Arg $_ : [AnyMatching AnyMatchingType]] [Arg $_ : [GenericParam a]]]) -> [ListLiteral List [[Fixed 23000]]]]]
[ModuleDef $main = [FunctionValue ([]) -> [FnCall [FunctionRef [NamedDefinitionReference :$fn]] [[String hello] [Integer -23939] [Fixed 42000]]]]]
`)
}

func TestAnyMatchingTypeCallMiddleLocalTypeFn(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
fn : String -> (* -> a) -> List a
fn _ _ =
    [ 23.0 ]


someOther : Int -> Fixed
someOther _ =
    3.5


main : List Fixed
main =
    fn "hello" someOther
`,
		`
[ModuleDef $fn = [FunctionValue ([[Arg $_ : [PrimitiveTypeVariantRef [NamedDefTypeRef :[TypeReference $String]]]] [Arg $_ : [FunctionTypeRef [FunctionType [[AnyMatching AnyMatchingType] [GenericParam a]]]]]) -> [ListLiteral [[Fixed 23000]]]]]
[ModuleDef $main = [FunctionValue ([]) -> [FnCall [FunctionRef [NamedDefinitionReference :$fn]] [[String hello] [FunctionRef [NamedDefinitionReference :$someOther]]]]]]
[ModuleDef $someOther = [FunctionValue ([[Arg $_ : [PrimitiveTypeVariantRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Fixed 3500]]]
`)
}

func TestFixed(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
--| ignore this
first : Fixed -> Fixed
first a =
    a * a
`,
		`
[ModuleDef $first = [FunctionValue ([[Arg $a : [PrimitiveTypeVariantRef [NamedDefTypeRef :[TypeReference $Fixed]]]]]) -> (arithmetic [FunctionParamRef [Arg $a : [PrimitiveTypeVariantRef [NamedDefTypeRef :[TypeReference $Fixed]]]]] FMULTIPLY [FunctionParamRef [Arg $a : [PrimitiveTypeVariantRef [NamedDefTypeRef :[TypeReference $Fixed]]]]])]]
`)
}

func TestChar(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
--| ignore this
first : Char -> Int
first _ =
    2


main : Bool -> Int
main _ =
    first 'c'
`,
		`
[ModuleDef $first = [FunctionValue ([[Arg $_ : [PrimitiveTypeVariantRef [NamedDefTypeRef :[TypeReference $Char]]]]]) -> [Integer 2]]]
[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeVariantRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference :$first]] [[Char 99]]]]]
`)
}

func TestFixedConvert(t *testing.T) {
	testDecorate(t,
		`
--| ignore this
first : Int -> Fixed
first a =
    Int.toFixed a


another : Int -> Fixed
another _ =
    first 2
`,
		`
[ModuleDef $another = [FunctionValue ([[Arg $_ : [PrimitiveTypeVariantRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /first]] [[Integer 2]]]]]
[ModuleDef $first = [FunctionValue ([[Arg $a : [PrimitiveTypeVariantRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference Int/toFixed]] [[FunctionParamRef [Arg $a : [PrimitiveTypeVariantRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]]
`)
}

func TestFixedConvertRound(t *testing.T) {
	testDecorate(t,
		`
--| ignore this
first : Fixed -> Int
first a =
    Int.round a


another : Int -> Int
another _ =
    first 2.3
`,
		`
[ModuleDef $another = [FunctionValue ([[Arg $_ : [PrimitiveTypeVariantRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /first]] [[Fixed 2300]]]]]
[ModuleDef $first = [FunctionValue ([[Arg $a : [PrimitiveTypeVariantRef [NamedDefTypeRef :[TypeReference $Fixed]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference Int/round]] [[FunctionParamRef [Arg $a : [PrimitiveTypeVariantRef [NamedDefTypeRef :[TypeReference $Fixed]]]]]]]]]
`)
}

func TestFixedConvertLet(t *testing.T) {
	testDecorate(t,
		`
--| ignore this
first : Int -> Fixed
first _ =
    0.3


another : Int -> Fixed
another _ =
    let
        x = first 2
    in
    x
`,
		`
[ModuleDef $another = [FunctionValue ([[Arg $_ : [PrimitiveTypeVariantRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Let [[LetAssign [[LetVar $x]] = [FnCall [FunctionRef [NamedDefinitionReference /first]] [[Integer 2]]]]] in [LetVarRef [LetVar $x]]]]]
[ModuleDef $first = [FunctionValue ([[Arg $_ : [PrimitiveTypeVariantRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Fixed 300]]]
`)
}

func TestFixedConvertRecordSet(t *testing.T) {
	testDecorate(t,
		`
type alias State =
    { playerX : Int
    }


another : Int -> State
another _ =
    let
        state = { playerX = 0 }
    in
    { state | playerX = 22 }
`,
		`
[ModuleDef $another = [FunctionValue ([[Arg $_ : [PrimitiveTypeVariantRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Let [[LetAssign [[LetVar $state]] = [record-literal record-type [[record-type-field $playerX [Primitive Int] (0)]][]] [0 = [Integer 0]]]]] in [record-literal record-type [[record-type-field $playerX [Primitive Int] (0)]][]] [0 = [Integer 22]]]]]]
`)
}

func TestCommentMulti(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
{-
   ignore this
      for sure

   -}
first : Int -> Int
first a =
    a * a
`,
		`
func(Int -> Int) : [func  [primitive Int] [primitive Int]]

first = [functionvalue ([[arg $a = [primitive Int]]]) -> (arithmetic [getvar $a [primitive Int]] MULTIPLY [getvar $a [primitive Int]])]
`)
}

func TestCommentMultiDoc(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
{-|
   ignore this
      for sure

   -}
first : Int -> Int
first a =
    a * a
`,
		`
func(Int -> Int) : [func  [primitive Int] [primitive Int]]

first = [functionvalue ([[arg $a = [primitive Int]]]) -> (arithmetic [getvar $a [primitive Int]] MULTIPLY [getvar $a [primitive Int]])]
`)
}

func TestSpacingMultiDoc(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
{-|
   ignore this
      for sure

   -}
first : Int -> Int
first a =
    a * a


{-
   multiline

-}
second : Int -> Int
second b =
    b + b


-- single line
third : Int -> Int
third c =
    c - c
`,
		`
func(Int -> Int) : [func  [primitive Int] [primitive Int]]

first = [functionvalue ([[arg $a = [primitive Int]]]) -> (arithmetic [getvar $a [primitive Int]] MULTIPLY [getvar $a [primitive Int]])]
second = [functionvalue ([[arg $b = [primitive Int]]]) -> (arithmetic [getvar $b [primitive Int]] PLUS [getvar $b [primitive Int]])]
third = [functionvalue ([[arg $c = [primitive Int]]]) -> (arithmetic [getvar $c [primitive Int]] MINUS [getvar $c [primitive Int]])]
`)
}

func xTestWrongSpacingMultiDoc(t *testing.T) {
	testDecorateFail(t,
		`
{-|
   ignore this
      for sure

   -}
first : Int -> Int
first a =
    a * a

{-
   something else

   -}
second : Int -> Int
second b =
    b + b
`, parerr.ExpectedTwoLinesAfterStatement{})
}

func TestSimpleErr(t *testing.T) {
	testDecorateFail(t,
		`
someFunc : Int -> String
someFunc name =
    "2"


another : Int -> String
another ignore =
    someFunc "2"
`,
		&decorated.FunctionArgumentTypeMismatch{})
}

func TestSomething2(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
first : Int -> Int
first a =
    a * a
`,
		`
func(Int -> Int) : [func  [primitive Int] [primitive Int]]

first = [functionvalue ([[arg $a = [primitive Int]]]) -> (arithmetic [getvar $a [primitive Int]] MULTIPLY [getvar $a [primitive Int]])]
`)
}

func TestCustomTypeVariantLiteral(t *testing.T) {
	testDecorate(t,
		`
type Status =
    Unknown
    | Something Int


receiveStatus : Status -> Status
receiveStatus status =
    status


someFunc : String -> Status
someFunc _ =
    receiveStatus Unknown
`,
		`
[mdefx $receiveStatus = [functionvalue ([[arg $status = [customtypevariantref named definition reference [custom-type [[variant $Unknown] [variant $Something [[customtypevariantref named definition reference [primitive Int]]]]]]]]]) -> [functionparamref $status [arg $status = [customtypevariantref named definition reference [custom-type [[variant $Unknown] [variant $Something [[customtypevariantref named definition reference [primitive Int]]]]]]]]]]]
[mdefx $someFunc = [functionvalue ([[arg $name = [customtypevariantref named definition reference [primitive String]]]]) -> [fcall [functionref named definition reference [functionvalue ([[arg $status = [customtypevariantref named definition reference [custom-type [[variant $Unknown] [variant $Something [[customtypevariantref named definition reference [primitive Int]]]]]]]]]) -> [functionparamref $status [arg $status = [customtypevariantref named definition reference [custom-type [[variant $Unknown] [variant $Something [[customtypevariantref named definition reference [primitive Int]]]]]]]]]]] [[variant-constructor [variant $Unknown] []]]]]]

`)
}

func TestCustomTypeVariantEqual(t *testing.T) {
	testDecorate(t,
		`
someFunc : String -> Bool
someFunc name =
    name == "Something"
`,
		`
func(String -> Bool) : [func  [primitive String] [primitive Bool]]

someFunc = [functionvalue ([[arg $name = [primitive String]]]) -> (boolop [getvar $name [primitive String]] EQ [str Something])]
`)
}

func TestUnknownAnnotationType(t *testing.T) {
	testDecorateFail(t, `
    someFunc : Position2
    `, &decorated.UnknownAnnotationTypeReference{})
}

func TestCustomTypeVariantLiteralWithParameters(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
type Status =
    Unknown
    | Something Int


someFunc : String -> Status
someFunc name =
    Something 42
`, `
Unknown : [variantconstr [variant $Unknown]]
Something : [variantconstr [variant $Something [[primitive Int]]]]
Status : [custom-type  [variant $Unknown] [variant $Something]]
func(String -> Status) : [func  [primitive String] [custom-type  [variant $Unknown] [variant $Something]]]

someFunc = [functionvalue ([[arg $name = [primitive String]]]) -> [variant-constructor [variant $Something [[primitive Int]]] [[integer 42]]]]
`)
}

func TestListListType(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
type alias Lister t u =
    { fake : t
    , another : u
    }


type alias Sprite =
    { x : Int
    }


someFunc : Lister Sprite Int -> Int
`, `
{another:u;fake:t} : [record-type  [record-field another [localtype u]] [record-field fake [localtype t]]]
{x:Int} : [record-type  [record-field x [primitive Int]]]
Sprite : [alias Sprite {x:Int}]
func(Lister<Sprite,Int> -> Int) : [func  Lister<Sprite,Int> [primitive Int]]
`)
}

func TestPipeRight(t *testing.T) {
	testDecorateWithoutDefault(t,
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
`, `
func(Int -> Int) : [func  [primitive Int] [primitive Int]]
func(String -> Int -> Bool) : [func  [primitive String] [primitive Int] [primitive Bool]]
func(Bool -> Bool) : [func  [primitive Bool] [primitive Bool]]
func(String -> Bool) : [func  [primitive String] [primitive Bool]]

first = [functionvalue ([[arg $a = [primitive Int]]]) -> (arithmetic [getvar $a [primitive Int]] MULTIPLY [getvar $a [primitive Int]])]
second = [functionvalue ([[arg $str = [primitive String]] [arg $a = [primitive Int]]]) -> (boolop [getvar $a [primitive Int]] GR [integer 25])]
tester = [functionvalue ([[arg $b = [primitive String]]]) -> [fcall [getvar $third [primitive Bool]] [[fcall [getvar $second [primitive Bool]] [[getvar $b [primitive String]] [fcall [getvar $first [primitive Int]] [(arithmetic [integer 2] PLUS [integer 2])]]]]]]]
third = [functionvalue ([[arg $a = [primitive Bool]]]) -> [getvar $a [primitive Bool]]]

`)
}

func TestBasicStringInterpolation(t *testing.T) {
	testDecorate(t,
		`
sample : Int -> String
sample a =
    "hello ${a}"
`, `
func(Int -> String) : [func  [primitive Int] [primitive String]]

sample = [functionvalue ([[arg $a = [primitive Int]]]) -> (arithmetic [str hello ] APPEND [fcall [getvar $Debug.toString [primitive String]] [[getvar $a [primitive Int]]]])]
`)
}

func TestArrayVsListFail(t *testing.T) {
	testDecorateFail(t,
		`

updater : Int -> String -> Int
updater a b =
    42


sample : Int -> List Int
sample a =
    let
        intArray = Array.fromList [ 0, 1, 2 ]

        arraySlice = Array.slice 0 2 intArray
    in
    List.map2 updater arraySlice [ "hello", "world", "fail" ]
`, &decorated.CouldNotSmashFunctions{})
}

func TestOwnListMap(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
__externalfn coreListMap 2

map : (a -> b) -> List a -> List b
map x y =
    __asm callexternal 00 coreListMap 01 02
	`, `
func(a -> b) : [func  [localtype a] [localtype b]]
func(func(a -> b) -> List<a> -> List<b>) : [func  [func  [localtype a] [localtype b]] List<a> List<b>]

map = [functionvalue ([[arg $x = [functype [[localtype a] [localtype b]]]] [arg $y = List<a>]]) -> [asm callexternal 00 coreListMap 01 02]]
`)
}

func TestBasicDecorate(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
type alias Position =
    { x : Int
    , y : Int
    }


type alias Sprite =
    { rootPosition : Position
    }


move : Position -> Position -> Position
move pos delta =
    let
        newX = pos.x + delta.x

        newY = pos.y + delta.y
    in
    { x = newX
    , y = newY
    }


moveSprite : Sprite -> Position -> Sprite
moveSprite sprite delta =
    { rootPosition = move sprite.rootPosition delta }
`, `
{x:Int;y:Int} : [record-type  [record-field x [primitive Int]] [record-field y [primitive Int]]]
Position : [alias Position {x:Int;y:Int}]
{rootPosition:Position} : [record-type  [record-field rootPosition [alias Position {x:Int;y:Int}]]]
Sprite : [alias Sprite {rootPosition:Position}]
func(Position -> Position -> Position) : [func  [alias Position {x:Int;y:Int}] [alias Position {x:Int;y:Int}] [alias Position {x:Int;y:Int}]]
func(Sprite -> Position -> Sprite) : [func  [alias Sprite {rootPosition:Position}] [alias Position {x:Int;y:Int}] [alias Sprite {rootPosition:Position}]]

move = [functionvalue ([[arg $pos = [alias Position record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]]]] [arg $delta = [alias Position record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]]]]]) -> [let [[letassign $newX = (arithmetic [lookups [lookupvar $pos ([alias Position record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]]])] [[lookup [record-type-field x [primitive Int] (0)]]]] PLUS [lookups [lookupvar $delta ([alias Position record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]]])] [[lookup [record-type-field x [primitive Int] (0)]]]])] [letassign $newY = (arithmetic [lookups [lookupvar $pos ([alias Position record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]]])] [[lookup [record-type-field y [primitive Int] (1)]]]] PLUS [lookups [lookupvar $delta ([alias Position record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]]])] [[lookup [record-type-field y [primitive Int] (1)]]]])]] in [record-literal record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]] [0 = [getvar $newX [primitive Int]] 1 = [getvar $newY [primitive Int]]]]]]
moveSprite = [functionvalue ([[arg $sprite = [alias Sprite record-type [[record-type-field rootPosition [alias Position record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]]] (0)]]]]] [arg $delta = [alias Position record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]]]]]) -> [record-literal record-type [[record-type-field rootPosition [alias Position record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]]] (0)]]] [0 = [fcall [getvar $move [alias Position record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]]]] [[lookups [lookupvar $sprite ([alias Sprite record-type [[record-type-field rootPosition [alias Position record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]]] (0)]]]])] [[lookup [record-type-field rootPosition [alias Position record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]]] (0)]]]] [getvar $delta [alias Position record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]]]]]]]]]
`)
}

func TestOperatorPipeRight(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
first : Int -> Int
first a =
    a * a


second : Int -> Bool
second a =
    a > 25


tester : String -> Bool
tester a =
    first (2 + 2) |> second
`, `
func(Int -> Int) : [func  [primitive Int] [primitive Int]]
func(Int -> Bool) : [func  [primitive Int] [primitive Bool]]
func(String -> Bool) : [func  [primitive String] [primitive Bool]]

first = [functionvalue ([[arg $a = [primitive Int]]]) -> (arithmetic [getvar $a [primitive Int]] MULTIPLY [getvar $a [primitive Int]])]
second = [functionvalue ([[arg $a = [primitive Int]]]) -> (boolop [getvar $a [primitive Int]] GR [integer 25])]
tester = [functionvalue ([[arg $a = [primitive String]]]) -> [fcall [getvar $second [primitive Bool]] [[fcall [getvar $first [primitive Int]] [(arithmetic [integer 2] PLUS [integer 2])]]]]]

`)
}

func TestOperatorPipeLeft(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
first : Int -> Int
first a =
    a * a


second : Int -> Bool
second a =
    a > 25


tester : String -> Bool
tester a =
    second <| first (2 + 2)
`, `
func(Int -> Int) : [func  [primitive Int] [primitive Int]]
func(Int -> Bool) : [func  [primitive Int] [primitive Bool]]
func(String -> Bool) : [func  [primitive String] [primitive Bool]]

first = [functionvalue ([[arg $a = [primitive Int]]]) -> (arithmetic [getvar $a [primitive Int]] MULTIPLY [getvar $a [primitive Int]])]
second = [functionvalue ([[arg $a = [primitive Int]]]) -> (boolop [getvar $a [primitive Int]] GR [integer 25])]
tester = [functionvalue ([[arg $a = [primitive String]]]) -> [fcall [getvar $second [primitive Bool]] [[fcall [getvar $first [primitive Int]] [(arithmetic [integer 2] PLUS [integer 2])]]]]]

`)
}

func TestOperatorPipeNextLine(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
first : Int -> Int
first a =
    a * a


second : Int -> Bool
second a =
    a > 25


tester : String -> Bool
tester a =
    first (2 + 2)
        |> second
`, `
func(Int -> Int) : [func  [primitive Int] [primitive Int]]
func(Int -> Bool) : [func  [primitive Int] [primitive Bool]]
func(String -> Bool) : [func  [primitive String] [primitive Bool]]

first = [functionvalue ([[arg $a = [primitive Int]]]) -> (arithmetic [getvar $a [primitive Int]] MULTIPLY [getvar $a [primitive Int]])]
second = [functionvalue ([[arg $a = [primitive Int]]]) -> (boolop [getvar $a [primitive Int]] GR [integer 25])]
tester = [functionvalue ([[arg $a = [primitive String]]]) -> [fcall [getvar $second [primitive Bool]] [[fcall [getvar $first [primitive Int]] [(arithmetic [integer 2] PLUS [integer 2])]]]]]

`)
}

func TestBasicCall(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
move : Int -> Int -> Int
move pos delta =
    pos + delta
`, `
func(Int -> Int -> Int) : [func  [primitive Int] [primitive Int] [primitive Int]]

move = [functionvalue ([[arg $pos = [primitive Int]] [arg $delta = [primitive Int]]]) -> (arithmetic [getvar $pos [primitive Int]] PLUS [getvar $delta [primitive Int]])]
`)
}

func TestBasicConstructor(t *testing.T) {
	testDecorate(t,
		`
type alias Constructor =
    { a : Int
    , b : Bool
    }


create : Int -> Constructor
create a =
    { a = 2, b = True }
`, `
{a:Int;b:Bool} : [record-type  [record-field a [primitive Int]] [record-field b [primitive Bool]]]
Constructor : [alias Constructor {a:Int;b:Bool}]
func(Int -> Constructor) : [func  [primitive Int] [alias Constructor {a:Int;b:Bool}]]

create = [functionvalue ([[arg $a = [primitive Int]]]) -> [record-literal record-type [[record-type-field a [primitive Int] (0)] [record-type-field b [primitive Bool] (1)]]] [0 = [integer 2] 1 = [bool true]]]]
`)
}

func TestConstant2(t *testing.T) {
	testDecorate(t,
		`
tileHeight : Int
tileHeight =
    2


create : Int -> Int
create a =
    tileHeight
`, `
func(Int) : [func  [primitive Int]]
func(Int -> Int) : [func  [primitive Int] [primitive Int]]

create = [functionvalue ([[arg $a = [primitive Int]]]) -> [integer 2]]
tileHeight = [functionvalue ([]) -> [integer 2]]
`)
}

func TestCustomTypeConstructor(t *testing.T) {
	testDecorate(t,
		`
type SomeCustomType =
    First String
    | Anon
    | Second Int


a : Bool -> SomeCustomType
a dummy =
    First "Hello"
`, `
First : [variantconstr [variant $First [[primitive String]]]]
Anon : [variantconstr [variant $Anon []]]
Second : [variantconstr [variant $Second [[primitive Int]]]]
SomeCustomType : [custom-type  [variant $First] [variant $Anon] [variant $Second]]
func(Bool -> SomeCustomType) : [func  [primitive Bool] [custom-type  [variant $First] [variant $Anon] [variant $Second]]]

a = [functionvalue ([[arg $dummy = [primitive Bool]]]) -> [variant-constructor [variant $First [[primitive String]]] [[str Hello]]]]
`)
}

func TestBasicCallFail(t *testing.T) {
	testDecorateFail(t,
		`
move : Int -> String -> Int
move pos delta =
    pos + delta
`, &decorated.UnMatchingBinaryOperatorTypes{})
}

func TestBasicCallLet(t *testing.T) {
	testDecorate(t,
		`
move : Int -> Int -> Int
move pos delta =
    let
        tenMore = pos + 10

        tenLess = delta - 10
    in
    tenMore + tenLess
`, `
func(Int -> Int -> Int) : [func  [primitive Int] [primitive Int] [primitive Int]]

move = [functionvalue ([[arg $pos = [primitive Int]] [arg $delta = [primitive Int]]]) -> [let [[letassign $tenMore = (arithmetic [getvar $pos [primitive Int]] PLUS [integer 10])] [letassign $tenLess = (arithmetic [getvar $delta [primitive Int]] MINUS [integer 10])]] in (arithmetic [getvar $tenMore [primitive Int]] PLUS [getvar $tenLess [primitive Int]])]]
`)
}

func TestBasicCallFromVariable(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
move : Int -> Int -> Int
move pos delta =
    pos + delta


main : Bool -> Int
main ignore =
    let
        fn = move
    in
    fn 2 3
`, `
func(Int -> Int -> Int) : [func  [primitive Int] [primitive Int] [primitive Int]]
func(Bool -> Int) : [func  [primitive Bool] [primitive Int]]

main = [functionvalue ([[arg $ignore = [primitive Bool]]]) -> [let [[letassign $fn = [getvar $move [functype [[primitive Int] [primitive Int] [primitive Int]]]]]] in [fcall [getvar $fn [primitive Int]] [[integer 2] [integer 3]]]]]
move = [functionvalue ([[arg $pos = [primitive Int]] [arg $delta = [primitive Int]]]) -> (arithmetic [getvar $pos [primitive Int]] PLUS [getvar $delta [primitive Int]])]

`)
}

func TestBasicAppend(t *testing.T) {
	testDecorate(t,
		`
a : Int -> List Int
a x =
    [ 1, 3, 4 ] ++ [ 5, 6, 7, 8 ]
`, `
func(Int -> List<Int>) : [func  [primitive Int] List<Int>]

a = [functionvalue ([[arg $x = [primitive Int]]]) -> (arithmetic [ListLiteral List<Int> [[integer 1] [integer 3] [integer 4]]] APPEND [ListLiteral List<Int> [[integer 5] [integer 6] [integer 7] [integer 8]]])]
`)
}

func TestBasicCallLetFail(t *testing.T) {
	testDecorateFail(t,
		`
move : Int -> Int -> String
move pos delta =
    let
        tenMore = pos + 10

        tenLess = delta - 10
    in
    tenMore + tenLess
`, &decorated.UnMatchingFunctionReturnTypesInFunctionValue{})
}

func TestBasicLookup(t *testing.T) {
	testDecorate(t,
		`
type alias Position =
    { x : Int
    , y : Int
    }


getx : Position -> Int
getx pos =
    pos.x + pos.y
`, `
{x:Int;y:Int} : [record-type  [record-field x [primitive Int]] [record-field y [primitive Int]]]
Position : [alias Position {x:Int;y:Int}]
func(Position -> Int) : [func  [alias Position {x:Int;y:Int}] [primitive Int]]

getx = [functionvalue ([[arg $pos = [alias Position record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]]]]]) -> (arithmetic [lookups [lookupvar $pos ([alias Position record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]]])] [[lookup [record-type-field x [primitive Int] (0)]]]] PLUS [lookups [lookupvar $pos ([alias Position record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]]])] [[lookup [record-type-field y [primitive Int] (1)]]]])]
`)
}

func TestSubLookup(t *testing.T) {
	testDecorate(t,
		`
type alias Position =
    { x : Int
    , y : Int
    }


type alias Sprite =
    { pos : Position
    }


getx : Sprite -> Int
getx sprite =
    sprite.pos.x
`, `
{x:Int;y:Int} : [record-type  [record-field x [primitive Int]] [record-field y [primitive Int]]]
Position : [alias Position {x:Int;y:Int}]
{pos:Position} : [record-type  [record-field pos [alias Position {x:Int;y:Int}]]]
Sprite : [alias Sprite {pos:Position}]
func(Sprite -> Int) : [func  [alias Sprite {pos:Position}] [primitive Int]]

getx = [functionvalue ([[arg $sprite = [alias Sprite record-type [[record-type-field pos [alias Position record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]]] (0)]]]]]]) -> [lookups [lookupvar $sprite ([alias Sprite record-type [[record-type-field pos [alias Position record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]]] (0)]]]])] [[lookup [record-type-field pos [alias Position record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]]] (0)]] [lookup [record-type-field x [primitive Int] (0)]]]]]
`)
}

func TestIf(t *testing.T) {
	testDecorate(t,
		`
isBestAge : Int -> Int
isBestAge age =
    if age == 50 || age >= 65 then
        100
    else
        0
`, `
func(Int -> Int) : [func  [primitive Int] [primitive Int]]

isBestAge = [functionvalue ([[arg $age = [primitive Int]]]) -> [if [logical (boolop [getvar $age [primitive Int]] EQ [integer 50]) (boolop [getvar $age [primitive Int]] GRE [integer 65]) 1] then [integer 100] else [integer 0]]]`)
}

func TestIfFail(t *testing.T) {
	testDecorateFail(t,
		`
isBestAge : Int -> Bool
isBestAge age =
    if age == 50 then
        100
    else
        0
`, &decorated.UnMatchingFunctionReturnTypesInFunctionValue{})
}

func TestIfPerson(t *testing.T) {
	testDecorate(t,
		`
isBeautiful : String -> Bool
isBeautiful name =
    if name == "Rebecca" then
        True
    else
        False
`, `
func(String -> Bool) : [func  [primitive String] [primitive Bool]]

isBeautiful = [functionvalue ([[arg $name = [primitive String]]]) -> [if (boolop [getvar $name [primitive String]] EQ [str Rebecca]) then [bool true] else [bool false]]]`)
}

func TestGuard(t *testing.T) {
	testDecorate(t,
		`
isBeautiful : String -> Bool
isBeautiful name =
    | name == "Rebecca" -> True
    | _ -> False
`, `
func(String -> Bool) : [func  [primitive String] [primitive Bool]]

isBeautiful = [functionvalue ([[arg $name = [primitive String]]]) -> [dguard: [dguarditem (boolop [getvar $name [primitive String]] EQ [str Rebecca]) [bool true]] default: [bool false]]]`)
}

func TestBoolPerson(t *testing.T) {
	testDecorate(t,
		`
isBeautiful : String -> Bool
isBeautiful name =
    name == "Rebecca"
`, `
func(String -> Bool) : [func  [primitive String] [primitive Bool]]

isBeautiful = [functionvalue ([[arg $name = [primitive String]]]) -> (boolop [getvar $name [primitive String]] EQ [str Rebecca])]`)
}

func TestBoolPersonCall(t *testing.T) {
	testDecorate(t,
		`
isLoveOfMyLife : String -> Bool
isLoveOfMyLife name =
    name == "Rebecca"


main : String -> Bool
main fake =
    isLoveOfMyLife "Lisa"
`, `
func(String -> Bool) : [func  [primitive String] [primitive Bool]]

isLoveOfMyLife = [functionvalue ([[arg $name = [primitive String]]]) -> (boolop [getvar $name [primitive String]] EQ [str Rebecca])]
main = [functionvalue ([[arg $fake = [primitive String]]]) -> [fcall [getvar $isLoveOfMyLife [primitive Bool]] [[str Lisa]]]]

`)
}

func TestBoolPersonCallAgain(t *testing.T) {
	testDecorate(t,
		`
isLoveOfMyLife : String -> Int -> Bool
isLoveOfMyLife name other =
    name == "Rebecca"


main : String -> Bool
main fake =
    isLoveOfMyLife "Lisa" 2
`, `
func(String -> Int -> Bool) : [func  [primitive String] [primitive Int] [primitive Bool]]
func(String -> Bool) : [func  [primitive String] [primitive Bool]]

isLoveOfMyLife = [functionvalue ([[arg $name = [primitive String]] [arg $other = [primitive Int]]]) -> (boolop [getvar $name [primitive String]] EQ [str Rebecca])]
main = [functionvalue ([[arg $fake = [primitive String]]]) -> [fcall [getvar $isLoveOfMyLife [primitive Bool]] [[str Lisa] [integer 2]]]]

`)
}

func TestCustomType(t *testing.T) {
	testDecorate(t,
		`
type Chore =
    Meeting String
    | Running Int
    | Unknown
`, `
Meeting : [variantconstr [variant $Meeting [[primitive String]]]]
Running : [variantconstr [variant $Running [[primitive Int]]]]
Unknown : [variantconstr [variant $Unknown]]
Chore : [custom-type  [variant $Meeting [primitive String]] [variant $Running [primitive Int]] [variant $Unknown
`)
}

func TestCustomTypeGenerics(t *testing.T) {
	testDecorate(t,
		`
type Perhaps a =
    None
    | Actual a
`, `
None : [variantconstr [variant $None]]
Actual : [variantconstr [variant $Actual [[localtype a]]]]
Perhaps : [custom-type  [variant $None] [variant $Actual [localtype a]]]
`)
}

func TestCustomTypeWithStructs(t *testing.T) {
	/*
		{-
		this is
		just some
		enum
		comment
		13890 - -)
	-} */
	testDecorate(t,
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
`, `
{solder:Bool} : [record-type  [record-field solder [primitive Bool]]]
Tinkering : [alias Tinkering {solder:Bool}]
{subject:String} : [record-type  [record-field subject [primitive String]]]
Studying : [alias Studying {subject:String}]
{web:Int} : [record-type  [record-field web [primitive Int]]]
Work : [alias Work {web:Int}]
Aron : [variantconstr [variant $Aron [[alias Tinkering record-type [[record-type-field solder [primitive Bool] (0)]]]]]]]
Alexandra : [variantconstr [variant $Alexandra]]
Alma : [variantconstr [variant $Alma [[alias Studying record-type [[record-type-field subject [primitive String] (0)]]]]]]]
Isabelle : [variantconstr [variant $Isabelle [[alias Work record-type [[record-type-field web [primitive Int] (0)]]]]]]]
Child : [custom-type  [variant $Aron [alias Tinkering {solder:Bool}]] [variant $Alexandra] [variant $Alma [alias Studying {subject:String}]] [variant $Isabelle [alias Work {web:Int}]]]
`)
}

func TestCurrying(t *testing.T) {
	testDecorate(t,
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
func(String -> Int -> Bool) : [func  [primitive String] [primitive Int] [primitive Bool]]
func(Int -> Bool) : [func  [primitive Int] [primitive Bool]]

another = [functionvalue ([[arg $score = [primitive Int]]]) -> [let [[letassign $af = [curry [getvar $f [functype [[primitive Int] [primitive Bool]]]] [[str Peter]]]]] in [fcall [getvar $af [primitive Bool]] [[getvar $score [primitive Int]]]]]]
f = [functionvalue ([[arg $name = [primitive String]] [arg $score = [primitive Int]]]) -> [if (boolop [getvar $name [primitive String]] EQ [str Peter]) then (boolop (arithmetic [getvar $score [primitive Int]] MULTIPLY [integer 2]) GR [integer 100]) else (boolop [getvar $score [primitive Int]] GR [integer 100])]]

`)
}

func TestBlob(t *testing.T) {
	testDecorate(t,
		`
a : Blob -> List Int
a x =
    [ 10, 20, 99 ]
`, `
func(Blob -> List<Int>) : [func  [primitive Blob] List<Int>]

a = [functionvalue ([[arg $x = [primitive Blob]]]) -> [ListLiteral List<Int> [[integer 10] [integer 20] [integer 99]]]]
`)
}

func TestListLiteral2(t *testing.T) {
	testDecorate(t,
		`
a : Bool -> List Int
a x =
    [ 10, 20, 99 ]
`, `
func(Bool -> List<Int>) : [func  [primitive Bool] List<Int>]

a = [functionvalue ([[arg $x = [primitive Bool]]]) -> [ListLiteral List<Int> [[integer 10] [integer 20] [integer 99]]]]
`)
}

func TestTuple(t *testing.T) {
	testDecorate(t,
		`
a : Bool -> (Int, String)
a x =
    (2, "Hello")
`, `
a = [functionvalue ([[arg $x = typeref $Bool [primitive Bool]]]) -> [TupleLiteral [tupletype [[primitive Int] [primitive String]]] [tuple-literal: [#2 'Hello']] [[integer 2] [str Hello]]]]
`)
}

func TestTupleSecond(t *testing.T) {
	testDecorate(t,
		`
a : Bool -> String
a x =
    Tuple.second (2, "Hello")
`, `
a = [functionvalue ([[arg $x = typeref $Bool [primitive Bool]]]) -> [fcall [functionref named definition reference [functionvalue ([[arg $tuple = [tupletype [[localtype a] [localtype b]]]]]) -> [asm callexternal 00 coreTupleSecond 01]]] [[TupleLiteral [tupletype [[primitive Int] [primitive String]]] [tuple-literal: [#2 'Hello']] [[integer 2] [str Hello]]]]]]
`)
}

func TestTupleFirst(t *testing.T) {
	testDecorate(t,
		`
a : Bool -> Int
a x =
    Tuple.first (2, "Hello")
`, `
a = [functionvalue ([[arg $x = typeref $Bool [primitive Bool]]]) -> [fcall [functionref named definition reference [functionvalue ([[arg $tuple = [tupletype [[localtype a] [localtype b]]]]]) -> [asm callexternal 00 coreTupleFirst 01]]] [[TupleLiteral [tupletype [[primitive Int] [primitive String]]] [tuple-literal: [#2 'Hello']] [[integer 2] [str Hello]]]]]]
`)
}

func TestArrayLiteral2(t *testing.T) {
	testDecorate(t,
		`
a : Bool -> Array Int
a x =
    [| 10, 20, 99 |]
`, `
func(Bool -> Array<Int>) : [func  [primitive Bool] Array<Int>]

a = [functionvalue ([[arg $x = [primitive Bool]]]) -> [ArrayLiteral Array<Int> [[integer 10] [integer 20] [integer 99]]]]
`)
}

func TestListLiteral3(t *testing.T) {
	testDecorate(t,
		`
type alias Cool =
    { name : String
    }


a : Bool -> List Cool
a x =
    [ { name = "hi" }, { name = "another" }, { name = "tjoho" } ]
`, `
{name:String} : [record-type  [record-field name [primitive String]]]
Cool : [alias Cool [record-type  [record-field name [primitive String]]]]
func(Bool -> List<Cool>) : [func  [primitive Bool] List<Cool>]

a = [functionvalue ([[arg $x = [primitive Bool]]]) -> [ListLiteral List<{name:String}> [[record-literal record-type [[record-type-field name [primitive String] (0)]]] [0 = [str hi]]] [record-literal record-type [[record-type-field name [primitive String] (0)]]] [0 = [str another]]] [record-literal record-type [[record-type-field name [primitive String] (0)]]] [0 = [str tjoho]]]]]]
`)
}

func TestFailListLiteral4(t *testing.T) {
	testDecorateFail(t,
		`
type alias Cool =
    { name : Int
    }


a : Bool -> List Cool
a x =
    [ { name = "hi" }, { name = "another" }, { name = "tjoho" } ]
`, &decorated.UnMatchingFunctionReturnTypesInFunctionValue{})
}

func TestRecordConstructor(t *testing.T) {
	testDecorate(t,
		`
type alias Cool =
    { name : Int
    }


a : Bool -> List Cool
a x =
    [ Cool{ name = 95 } ]
`, `
{name:Int} : [record-type  [record-field name [primitive Int]]]
Cool : [alias Cool [record-type  [record-field name [primitive Int]]]]
func(Bool -> List<Cool>) : [func  [primitive Bool] List<Cool>]

a = [functionvalue ([[arg $x = [primitive Bool]]]) -> [ListLiteral List<{name:Int}> [[record-constructor-record $Cool [record-literal record-type [[record-type-field name [primitive Int] (0)]]] [0 = [integer 95]]]]]]]
`)
}

func TestRecordConstructorValues(t *testing.T) {
	testDecorate(t,
		`
type alias Cool =
    { name : Int
    }


a : Bool -> List Cool
a x =
    [ Cool 2 ]
`, `
{name:Int} : [record-type  [record-field name [primitive Int]]]
Cool : [alias Cool [record-type  [record-field name [primitive Int]]]]
func(Bool -> List<Cool>) : [func  [primitive Bool] List<Cool>]

a = [functionvalue ([[arg $x = [primitive Bool]]]) -> [ListLiteral List<{name:Int}> [[record-constructor $Cool [0 = [integer 2]]]]]]

`)
}

func TestRecordConstructorValuesWrong(t *testing.T) {
	testDecorateFail(t,
		`
type alias Cool =
    { name : Int
    }


a : Bool -> List Cool
a x =
    [ Cool "2" ]
`, &decorated.WrongTypeForRecordConstructorField{})
}

func TestCaseDefault(t *testing.T) {
	testDecorate(t, // -- just a comment
		`
type CustomType =
    First
    | Second


some : CustomType -> String
some a =
    case a of
        _ -> ""
`, `
First : [variantconstr [variant $First]]
Second : [variantconstr [variant $Second]]
CustomType : [custom-type  [variant $First] [variant $Second]]
func(CustomType -> String) : [func  [custom-type  [variant $First] [variant $Second]] [primitive String]]

some = [functionvalue ([[arg $a = [custom-type [[variant $First] [variant $Second]]]]]) -> [dcase: [getvar $a [custom-type [[variant $First] [variant $Second]]]] of  default: [str ]]]
`)
}

func TestCaseCustomTypeWithStructs(t *testing.T) {
	testDecorate(t, // -- just a comment
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
        Aron x -> "Aron"

        Alexandra -> "Alexandris"

        _ -> ""
`, `
{solder:Bool} : [record-type  [record-field solder [primitive Bool]]]
Tinkering : [alias Tinkering {solder:Bool}]
{subject:String} : [record-type  [record-field subject [primitive String]]]
Studying : [alias Studying {subject:String}]
{web:Int} : [record-type  [record-field web [primitive Int]]]
Work : [alias Work {web:Int}]
Aron : [variantconstr [variant $Aron [[alias Tinkering record-type [[record-type-field solder [primitive Bool] (0)]]]]]]]
Alexandra : [variantconstr [variant $Alexandra []]]
Alma : [variantconstr [variant $Alma [[alias Studying record-type [[record-type-field subject [primitive String] (0)]]]]]]]
Isabelle : [variantconstr [variant $Isabelle [[alias Work record-type [[record-type-field web [primitive Int] (0)]]]]]]]
Child : [custom-type  [variant $Aron] [variant $Alexandra] [variant $Alma] [variant $Isabelle]]
func(Child -> String) : [func  [custom-type  [variant $Aron] [variant $Alexandra] [variant $Alma] [variant $Isabelle]] [primitive String]]

some = [functionvalue ([[arg $child = [custom-type [[variant $Aron [[alias Tinkering record-type [[record-type-field solder [primitive Bool] (0)]]]]]] [variant $Alexandra []] [variant $Alma [[alias Studying record-type [[record-type-field subject [primitive String] (0)]]]]]] [variant $Isabelle [[alias Work record-type [[record-type-field web [primitive Int] (0)]]]]]]]]]]) -> [dcase: [getvar $child [custom-type [[variant $Aron [[alias Tinkering record-type [[record-type-field solder [primitive Bool] (0)]]]]]] [variant $Alexandra []] [variant $Alma [[alias Studying record-type [[record-type-field subject [primitive String] (0)]]]]]] [variant $Isabelle [[alias Work record-type [[record-type-field web [primitive Int] (0)]]]]]]]]] of [dcasecons $Aron ([[dcaseparm $x type:[alias Tinkering record-type [[record-type-field solder [primitive Bool] (0)]]]]]]) => [str Aron]];[dcasecons $Alexandra ([]) => [str Alexandris]] default: [str ]]]`)
}

func TestCaseStringAndDefault(t *testing.T) {
	testDecorate(t, // -- just a comment
		`
some : String -> Int
some a =
    case a of
        "hello" -> 0

        _ -> -1
`, `
func(String -> Int) : [func  [primitive String] [primitive Int]]

some = [functionvalue ([[arg $a = [primitive String]]]) -> [dpmcase: [getvar $a [primitive String]] of [dpmcasecons [str hello] => [integer 0]] default: [integer -1]]]
`)
}

func TestGenerics(t *testing.T) {
	testDecorate(t,
		`
type alias Tinkering t =
    { solder : Bool
    , secret : t
    }


f : Tinkering Int -> Int
f tinkering =
    tinkering.secret
`, `
Tinkering : [alias Tinkering RecordType([t])]
func(Tinkering<Int> -> Int) : [func  Tinkering<Int> [primitive Int]]
{secret:Int;solder:Bool} : [record-type  [record-field secret [primitive Int]] [record-field solder [primitive Bool]]]

f = [functionvalue ([[arg $tinkering = Tinkering<Int>]]) -> [lookups [lookupvar $tinkering (Tinkering<Int>)] [[lookup [record-type-field secret [primitive Int] (0)]]]]]
`)
}

func TestGenericsFail(t *testing.T) {
	testDecorateFail(t,
		`
type alias Tinkering t =
    { solder : Bool
    , secret : t
    }


f : Tinkering String -> Int
f tinkering =
    tinkering.secret
`, &decorated.UnMatchingFunctionReturnTypesInFunctionValue{})
}

func TestGenericsStructInstantiate(t *testing.T) {
	testDecorate(t,
		`
type alias Tinkering t =
    { solder : Bool
    , secret : t
    }
`, `
{secret:localtype;solder:Bool} : [record-type  [record-field secret [localtype [type-param $t]]] [record-field solder [primitive Bool]]]
Tinkering : [alias Tinkering {secret:localtype;solder:Bool}]
`)
}

func TestRecordList(t *testing.T) {
	testDecorate(t, //-- just a comment
		`
type alias Enemy =
    { values : List Int
    }


type alias World =
    { enemies : List Enemy
    }


updateEnemy : List Enemy -> Bool
updateEnemy enemies =
    True


updateWorld : World -> Bool
updateWorld w =
    updateEnemy w.enemies


main : Bool -> Bool
main ignore =
    updateWorld { enemies = [ { values = [ 1, 3 ] } ] }
`, `
{values:List<Int>} : [record-type  [record-field values List<Int>]]
Enemy : [alias Enemy {values:List<Int>}]
{enemies:List<Enemy>} : [record-type  [record-field enemies List<Enemy>]]
World : [alias World {enemies:List<Enemy>}]
func(List<Enemy> -> Bool) : [func  List<Enemy> [primitive Bool]]
func(World -> Bool) : [func  [alias World {enemies:List<Enemy>}] [primitive Bool]]
func(Bool -> Bool) : [func  [primitive Bool] [primitive Bool]]
{enemies:List<{values:List<Int>}>} : [record-type  [record-field enemies [concrcolltypes List [record-type  [record-field values List<Int>]]]]]

main = [functionvalue ([[arg $ignore = [primitive Bool]]]) -> [fcall [getvar $updateWorld [primitive Bool]] [[record-literal record-type [[record-type-field enemies [concrcolltype [List [a]] [record-type [[record-type-field values List<Int> (0)]]]]] (0)]]] [0 = [ListLiteral List<{values:List<Int>}> [[record-literal record-type [[record-type-field values List<Int> (0)]]] [0 = [ListLiteral List<Int> [[integer 1] [integer 3]]]]]]]]]]]]
updateEnemy = [functionvalue ([[arg $enemies = List<Enemy>]]) -> [bool true]]
updateWorld = [functionvalue ([[arg $w = [alias World record-type [[record-type-field enemies List<Enemy> (0)]]]]]]) -> [fcall [getvar $updateEnemy [primitive Bool]] [[lookups [lookupvar $w ([alias World record-type [[record-type-field enemies List<Enemy> (0)]]]])] [[lookup [record-type-field enemies List<Enemy> (0)]]]]]]]
`)
}

func TestRecordListInList(t *testing.T) {
	testDecorate(t,

		`
type alias Sprite =
    { x : Int
    , y : Int
    }


type alias World =
    { drawTasks : List (List Sprite)
    }


drawSprite : Sprite -> Bool
drawSprite sprite =
    True


drawSprites : List Sprite -> List Bool
drawSprites sprites =
    List.map drawSprite sprites


drawWorld : World -> List (List Bool)
drawWorld world =
    List.map drawSprites world.drawTasks


main : Bool -> List (List Bool)
main ignore =
    drawWorld { drawTasks = [ [ { x = 10, y = 20 }, { x = 44, y = 98 } ], [ { x = 99, y = 98 } ] ] }
`, `
{x:Int;y:Int} : [record-type  [record-field x [primitive Int]] [record-field y [primitive Int]]]
Sprite : [alias Sprite {x:Int;y:Int}]
{drawTasks:List<List<Sprite>>} : [record-type  [record-field drawTasks List<List<Sprite>>]]
World : [alias World {drawTasks:List<List<Sprite>>}]
func(Sprite -> Bool) : [func  [alias Sprite {x:Int;y:Int}] [primitive Bool]]
func(List<Sprite> -> List<Bool>) : [func  List<Sprite> List<Bool>]
func(World -> List<List<Bool>>) : [func  [alias World {drawTasks:List<List<Sprite>>}] List<List<Bool>>]
func(Bool -> List<List<Bool>>) : [func  [primitive Bool] List<List<Bool>>]
{drawTasks:List<List<{x:Int;y:Int}>>} : [record-type  [record-field drawTasks [concrcolltypes List [concrcolltypes List [record-type  [record-field x [primitive Int]] [record-field y [primitive Int]]]]]]]

drawSprite = [functionvalue ([[arg $sprite = [alias Sprite record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]]]]]) -> [bool true]]
drawSprites = [functionvalue ([[arg $sprites = List<Sprite>]]) -> [fcall [getvar List.$map [concrcolltype [List [a]] [[primitive Bool]]]] [[getvar $drawSprite [functype [[alias Sprite record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]]] [primitive Bool]]]] [getvar $sprites List<Sprite>]]]]
drawWorld = [functionvalue ([[arg $world = [alias World record-type [[record-type-field drawTasks List<List<Sprite>> (0)]]]]]]) -> [fcall [getvar List.$map [concrcolltype [List [a]] [List<Bool>]]] [[getvar $drawSprites [functype [List<Sprite> List<Bool>]]] [lookups [lookupvar $world ([alias World record-type [[record-type-field drawTasks List<List<Sprite>> (0)]]]])] [[lookup [record-type-field drawTasks List<List<Sprite>> (0)]]]]]]]
main = [functionvalue ([[arg $ignore = [primitive Bool]]]) -> [fcall [getvar $drawWorld List<List<Bool>>] [[record-literal record-type [[record-type-field drawTasks [concrcolltype [List [a]] [[concrcolltype [List [a]] [record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]]]]]] (0)]]] [0 = [ListLiteral List<List<{x:Int;y:Int}>> [[ListLiteral List<{x:Int;y:Int}> [[record-literal record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]] [0 = [integer 10] 1 = [integer 20]]] [record-literal record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]] [0 = [integer 44] 1 = [integer 98]]]]] [ListLiteral List<{x:Int;y:Int}> [[record-literal record-type [[record-type-field x [primitive Int] (0)] [record-type-field y [primitive Int] (1)]]] [0 = [integer 99] 1 = [integer 98]]]]]]]]]]]]

`)
}

func TestAppliedAnnotation(t *testing.T) {
	testDecorateWithoutDefault(t,

		`
type alias MyList t =
    { someType : t
    }


intConvert : MyList Int -> Bool
`, `
MyList : [alias MyList RecordType([t])]
func(MyList<Int> -> Bool) : [func  MyList<Int> [primitive Bool]]
`)
}

func TestAppliedAnnotation2(t *testing.T) {
	testDecorateWithoutDefault(t,

		`
type alias MyList t =
    { someType : t
    }


intConvert : MyList Int -> Bool
intConvert intList =
    True
`, `
{someType:t} : {someType:t}
MyList : [alias MyList {someType:t}]
func(MyList<Int> -> Bool) : [func  MyList<Int> [primitive Bool]]

intConvert = [functionvalue ([[arg $intList = MyList<Int>]]) -> [bool true]]
`)
}

func TestBasicFunctionTypeParameters(t *testing.T) {
	testDecorate(t,

		`
simple : List a -> a
simple lst =
    2


main : Bool -> Int
main ignored =
    simple [ 1, 2, 3 ]
`, `
func(List<a> -> a) : [func  List<a> [localtype a]]
func(Bool -> Int) : [func  [primitive Bool] [primitive Int]]

main = [functionvalue ([[arg $ignored = [primitive Bool]]]) -> [fcall [getvar $simple [primitive Int]] [[ListLiteral List<Int> [[integer 1] [integer 2] [integer 3]]]]]]
simple = [functionvalue ([[arg $lst = List<a>]]) -> [integer 2]]

`)
}

func TestFunctionTypeParameters(t *testing.T) {
	testDecorate(t,

		`
intConvert : List Int -> Bool
intConvert lst =
    True


simple : (List a -> Bool) -> List a -> Bool
simple fn lst =
    True


main : Bool -> Bool
main ignored =
    simple intConvert [ 1, 2, 3 ]
`, `
func(List<Int> -> Bool) : [func  List<Int> [primitive Bool]]
func(List<localtype> -> Bool) : [func  List<localtype> [primitive Bool]]
func(func(List<localtype> -> Bool) -> List<localtype> -> Bool) : [func  [func  List<localtype> [primitive Bool]] List<localtype> [primitive Bool]]
func(Bool -> Bool) : [func  [primitive Bool] [primitive Bool]]

intConvert = [functionvalue ([[arg $lst = List<Int>]]) -> [bool true]]
main = [functionvalue ([[arg $ignored = [primitive Bool]]]) -> [fcall [getvar $simple [primitive Bool]] [[getvar $intConvert [functype [List<Int> [primitive Bool]]]] [ListLiteral List<Int> [[integer 1] [integer 2] [integer 3]]]]]]
simple = [functionvalue ([[arg $fn = [functype [List<localtype> [primitive Bool]]]] [arg $lst = List<localtype>]]) -> [bool true]]
`)
}

func TestUpdateComplex(t *testing.T) {
	testDecorate(t,

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
{scaleX:Int;scaleY:Int} : [record-type  [record-field scaleX [primitive Int]] [record-field scaleY [primitive Int]]]
Scale2 : [alias Scale2 {scaleX:Int;scaleY:Int}]
{dummy:Int;scale:Scale2} : [record-type  [record-field dummy [primitive Int]] [record-field scale [alias Scale2 {scaleX:Int;scaleY:Int}]]]
Sprite : [alias Sprite {dummy:Int;scale:Scale2}]
func(Sprite -> Int -> Sprite) : [func  [alias Sprite {dummy:Int;scale:Scale2}] [primitive Int] [alias Sprite {dummy:Int;scale:Scale2}]]

updateSprite = [functionvalue ([[arg $inSprite = [alias Sprite record-type [[record-type-field dummy [primitive Int] (0)] [record-type-field scale [alias Scale2 record-type [[record-type-field scaleX [primitive Int] (0)] [record-type-field scaleY [primitive Int] (1)]]]] (1)]]]]] [arg $newScale = [primitive Int]]]) -> [record-literal record-type [[record-type-field dummy [primitive Int] (0)] [record-type-field scale [alias Scale2 record-type [[record-type-field scaleX [primitive Int] (0)] [record-type-field scaleY [primitive Int] (1)]]]] (1)]]] [1 = [record-literal record-type [[record-type-field scaleX [primitive Int] (0)] [record-type-field scaleY [primitive Int] (1)]]] [0 = [getvar $newScale [primitive Int]] 1 = [getvar $newScale [primitive Int]]]]]]]
`)
}

func TestCaseAlignmentSurprisedByIndentation(t *testing.T) {
	testDecorate(t,
		`
checkMaybeInt : Maybe Int -> Maybe Int -> Int
checkMaybeInt a b =
    case a of
        Just firstInt -> case b of
            Nothing -> 0

            Just secondInt -> secondInt

        Nothing -> 0
      `,
		`
func(Maybe<Int> -> Maybe<Int> -> Int) : [func  Maybe<Int> Maybe<Int> [primitive Int]]

checkMaybeInt = [functionvalue ([[arg $a = Maybe<Int>] [arg $b = Maybe<Int>]]) -> [dcase: [getvar $a Maybe<Int>] of [dcasecons $Just ([[dcaseparm $firstInt type:[primitive Int]]]) => [dcase: [getvar $b Maybe<Int>] of [dcasecons $Nothing ([]) => [integer 0]];[dcasecons $Just ([[dcaseparm $secondInt type:[primitive Int]]]) => [getvar $secondInt [primitive Int]]]]];[dcasecons $Nothing ([]) => [integer 0]]]]
`)
}

func TestCaseNotCoveredByAllConsequences(t *testing.T) {
	testDecorateFail(t, //    -- Must include default (_) or Nothing here
		`
checkMaybeInt : Maybe Int -> Int
checkMaybeInt a =
    case a of
        Just firstInt -> firstInt
`,
		&decorated.UnhandledCustomTypeVariants{})
}

func TestCaseCoveredMultipleTimes(t *testing.T) {
	testDecorateFail(t,
		`
checkMaybeInt : Maybe Int -> Int
checkMaybeInt a =
    case a of
        Just firstInt -> firstInt

        Nothing -> 0

        Nothing -> 0
`,
		&decorated.AlreadyHandledCustomTypeVariant{})
}

func TestFunctionAnnotationWithoutParameters(t *testing.T) {
	testDecorate(t,
		`
single : Int
single =
    2
`,
		`
func(Int) : [func  [primitive Int]]

single = [functionvalue ([]) -> [integer 2]]
    `)
}

func TestMinimalMaybeAndNothing(t *testing.T) {
	testDecorate(t,
		`
type alias Thing =
    { something : String
    }


main : Bool -> Maybe Thing
main ignore =
    Nothing
`, `
{something:String} : [record-type  [record-field something [primitive String]]]
Thing : [alias Thing {something:String}]
func(Bool -> Maybe<Thing>) : [func  [primitive Bool] Maybe<Thing>]

main = [functionvalue ([[arg $ignore = [primitive Bool]]]) -> [variant-constructor [variant $Nothing []] []]]
`)
}

func TestConsequenceCheck(t *testing.T) {
	testDecorateFail(t, `
type Direction =
    NotMoving
    | Right
    | Down
    | Left
    | Up


-- returns true if direction is vertical
isVertical : Direction -> Bool
isVertical direction =
    case direction of
        NotMoving ->
            False

        Right ->
            False

        Down ->
            True

        Left ->
            2

        Up ->
            True
`, &decorated.UnMatchingTypesExpression{})
}

func TestConsequenceCheck3(t *testing.T) {
	testDecorateFail(t, `
type Direction =
    NotMoving
    | Right
    | Down
    | Left
    | Up


-- returns true if direction is vertical
isVertical : Direction -> Bool
isVertical direction =
    case direction of
        NotMoving ->
            False

        Right ->
            False

        Down ->
            True

        Left ->
            False

        Up ->
            True

        _ ->
            2
`, &decorated.UnMatchingTypesExpression{})
}

func TestConsequenceCheck2(t *testing.T) {
	testDecorateFail(t, `
type Direction =
    NotMoving
    | Right
    | Down
    | Left
    | Up


isVertical : Direction -> Int
isVertical direction =
    case direction of
        NotMoving ->
            False

        Right ->
            False

        Down ->
            True

        Left ->
            False

        Up ->
            True

`, &decorated.UnMatchingFunctionReturnTypesInFunctionValue{})
}

func TestMinimalMaybeAndJust(t *testing.T) {
	testDecorate(t,
		`
type alias Thing =
    { something : String
    }


main : Bool -> Maybe Thing
main ignore =
    Just { something = "hello" }
`, `
{something:String} : [record-type  [record-field something [primitive String]]]
Thing : [alias Thing {something:String}]
func(Bool -> Maybe<Thing>) : [func  [primitive Bool] Maybe<Thing>]

main = [functionvalue ([[arg $ignore = [primitive Bool]]]) -> [variant-constructor [variant $Just [record-type [[record-type-field something [primitive String] (0)]]]]] [[record-literal record-type [[record-type-field something [primitive String] (0)]]] [0 = [str hello]]]]]]
`)
}

func TestJustAndNothing(t *testing.T) {
	testDecorate(t,
		`
type alias Thing =
    { something : String
    }


main : Bool -> Maybe Thing
main ignore =
    if True then
        Nothing
    else
        Just { something = "hello" }
`, `
{something:String} : [record-type  [record-field something [primitive String]]]
Thing : [alias Thing {something:String}]
func(Bool -> Maybe<Thing>) : [func  [primitive Bool] Maybe<Thing>]

main = [functionvalue ([[arg $ignore = [primitive Bool]]]) -> [if [bool true] then [variant-constructor [variant $Nothing []] []] else [variant-constructor [variant $Just [record-type [[record-type-field something [primitive String] (0)]]]]] [[record-literal record-type [[record-type-field something [primitive String] (0)]]] [0 = [str hello]]]]]]]
`)
}

func TestArrayGetWithMaybe(t *testing.T) {
	testDecorate(t,
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
None : [variantconstr [variant $None []]]
Jump : [variantconstr [variant $Jump []]]
PlayerAction : [custom-type  [variant $None] [variant $Jump]]
{a:Int} : [record-type  [record-field a [primitive Int]]]
Gamepad : [alias Gamepad {a:Int}]
{gamepads:Array<Gamepad>} : [record-type  [record-field gamepads Array<Gamepad>]]
UserInput : [alias UserInput {gamepads:Array<Gamepad>}]
{inputs:List<PlayerAction>} : [record-type  [record-field inputs List<PlayerAction>]]
PlayerInputs : [alias PlayerInputs {inputs:List<PlayerAction>}]
func(Maybe<Gamepad> -> PlayerAction) : [func  Maybe<Gamepad> [custom-type  [variant $None] [variant $Jump]]]
func(UserInput -> PlayerInputs) : [func  [alias UserInput {gamepads:Array<Gamepad>}] [alias PlayerInputs {inputs:List<PlayerAction>}]]

checkMaybeGamepad = [functionvalue ([[arg $a = Maybe<Gamepad>]]) -> [variant-constructor [variant $None []] []]]
main = [functionvalue ([[arg $oldUserInputs = [alias UserInput record-type [[record-type-field gamepads Array<Gamepad> (0)]]]]]]) -> [let [[letassign $gamepads = [lookups [lookupvar $oldUserInputs ([alias UserInput record-type [[record-type-field gamepads Array<Gamepad> (0)]]]])] [[lookup [record-type-field gamepads Array<Gamepad> (0)]]]]] [letassign $maybeOld = [fcall [getvar Array.$get [custom-type [[variant $Nothing []] [variant $Just [[alias Gamepad record-type [[record-type-field a [primitive Int] (0)]]]]]]]]] [[integer 0] [getvar $gamepads Array<Gamepad>]]]] [letassign $playerAction = [fcall [getvar $checkMaybeGamepad [custom-type [[variant $None []] [variant $Jump []]]]] [[getvar $maybeOld [custom-type [[variant $Nothing []] [variant $Just [[alias Gamepad record-type [[record-type-field a [primitive Int] (0)]]]]]]]]]]]]] in [record-literal record-type [[record-type-field inputs List<PlayerAction> (0)]]] [0 = [ListLiteral List<PlayerAction> [[getvar $playerAction [custom-type [[variant $None []] [variant $Jump []]]]]]]]]]]
`)
}

func TestArrayFromNothing(t *testing.T) {
	testDecorate(t,
		`
needInputs : Array Int -> Int
needInputs ignore =
    2


a : Bool -> Int
a ignore =
    needInputs (Array.fromList [])
`,
		`
func(Array<Int> -> Int) : [func  Array<Int> [primitive Int]]
func(Bool -> Int) : [func  [primitive Bool] [primitive Int]]

a = [functionvalue ([[arg $ignore = [primitive Bool]]]) -> [fcall [getvar $needInputs [primitive Int]] [[fcall [getvar Array.$fromList [concrcolltype [Array [a]] [[any]]]] [[ListLiteral List<any> []]]]]]]
needInputs = [functionvalue ([[arg $ignore = Array<Int>]]) -> [integer 2]]
`)
}

func TestReturnFunction(t *testing.T) {
	testDecorate(t,
		`
needInputs : Array Int -> Int
needInputs ignore =
    2


returnSomeFunc : Bool -> (Array Int -> Int)
returnSomeFunc ignore =
    needInputs


a : Bool -> Int
a ignore =
    let
        fn = returnSomeFunc ignore
    in
    fn (Array.fromList [])
`,
		`
func(Array<Int> -> Int) : [func  Array<Int> [primitive Int]]
func(Bool -> func(Array<Int> -> Int)) : [func  [primitive Bool] [func  Array<Int> [primitive Int]]]
func(Bool -> Int) : [func  [primitive Bool] [primitive Int]]

a = [functionvalue ([[arg $ignore = [primitive Bool]]]) -> [let [[letassign $fn = [fcall [getvar $returnSomeFunc [functype [Array<Int> [primitive Int]]]] [[getvar $ignore [primitive Bool]]]]]] in [fcall [getvar $fn [primitive Int]] [[fcall [getvar Array.$fromList [concrcolltype [Array [a]] [[any]]]] [[ListLiteral List<any> []]]]]]]]
needInputs = [functionvalue ([[arg $ignore = Array<Int>]]) -> [integer 2]]
returnSomeFunc = [functionvalue ([[arg $ignore = [primitive Bool]]]) -> [getvar $needInputs [functype [Array<Int> [primitive Int]]]]]
`)
}

func TestUpdateRecordDescendingOrder(t *testing.T) {
	testDecorate(t,
		`
type alias Something =
    { time : Int
    , playerX : Int
    }


fn : Something -> Something
fn something =
    { something | playerX = 3 }
`,
		`
{playerX:Int;time:Int} : [record-type  [record-field playerX [primitive Int]] [record-field time [primitive Int]]]
Something : [alias Something {playerX:Int;time:Int}]
func(Something -> Something) : [func  [alias Something {playerX:Int;time:Int}] [alias Something {playerX:Int;time:Int}]]

fn = [functionvalue ([[arg $something = [alias Something record-type [[record-type-field playerX [primitive Int] (0)] [record-type-field time [primitive Int] (1)]]]]]]) -> [record-literal record-type [[record-type-field playerX [primitive Int] (0)] [record-type-field time [primitive Int] (1)]]] [0 = [integer 3]]]]
`)
}

func TestPipeRight2(t *testing.T) {
	testDecorate(t, `
first : Int -> Int
first a =
    2


second : String -> Int -> Int
second a b =
    2


third : Int -> Bool
third a =
    a > 45


tester : String -> Bool
tester b =
    first 2 + 2
        |> second b
        |> third
`, `
func(Int -> Int) : [func  [primitive Int] [primitive Int]]
func(String -> Int -> Int) : [func  [primitive String] [primitive Int] [primitive Int]]
func(Int -> Bool) : [func  [primitive Int] [primitive Bool]]
func(String -> Bool) : [func  [primitive String] [primitive Bool]]

first = [functionvalue ([[arg $a = [primitive Int]]]) -> [integer 2]]
second = [functionvalue ([[arg $a = [primitive String]] [arg $b = [primitive Int]]]) -> [integer 2]]
tester = [functionvalue ([[arg $b = [primitive String]]]) -> [fcall [getvar $third [primitive Bool]] [[fcall [getvar $second [primitive Int]] [[getvar $b [primitive String]] (arithmetic [fcall [getvar $first [primitive Int]] [[integer 2]]] PLUS [integer 2])]]]]]
third = [functionvalue ([[arg $a = [primitive Int]]]) -> (boolop [getvar $a [primitive Int]] GR [integer 45])]
`)
}
