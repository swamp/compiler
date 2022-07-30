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
[ModuleDef $first = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
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
[ModuleDef $fn = [Constant [String Hello]]]
[ModuleDef $another = [FunctionValue ([]) -> [ConstantRef [NamedDefinitionReference /fn]]]]
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
Something : [Alias Something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [Primitive Int]

[ModuleDef $another = [FunctionValue ([]) -> [Let [[LetAssign [[LetVar $b]] = [Cast [Integer 32] [AliasRefExpr [AliasRef [Alias Something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]] in [BoolOp [LetVarRef [LetVar $b]] GRE [Integer 32]]]]]
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
[ModuleDef $first = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $ResourceName]]]]]) -> [Integer 2]]]
[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /first]] [[ResourceName this/is/something.txt]]]]]
`)
}

func TestAnyMatchingType(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
fn : * -> Int
fn _ _ =
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
fn _ =
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
fn _ _ =
    23


main : Int
main =
    fn "hello" 42.0
`,
		`
[ModuleDef $fn = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]] [Arg $_ : [AnyMatching AnyMatchingType]]]) -> [Integer 23]]]
[ModuleDef $main = [FunctionValue ([]) -> [FnCall [FunctionRef [NamedDefinitionReference /fn]] [[String hello] [Fixed 42000]]]]]
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
[ModuleDef $fn = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]] [Arg $_ : [AnyMatching AnyMatchingType]] [Arg $_ : [GenericParam a]]]) -> [ListLiteral [[Fixed 23000]]]]]
[ModuleDef $main = [FunctionValue ([]) -> [FnCall [FunctionRef [NamedDefinitionReference /fn]] [[String hello] [Integer -23939] [Fixed 42000]]]]]
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
[ModuleDef $fn = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]] [Arg $_ : [FunctionTypeRef [FunctionType [[AnyMatching AnyMatchingType] [GenericParam a]]]]]) -> [ListLiteral [[Fixed 23000]]]]]
[ModuleDef $someOther = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Fixed 3500]]]
[ModuleDef $main = [FunctionValue ([]) -> [FnCall [FunctionRef [NamedDefinitionReference /fn]] [[String hello] [FunctionRef [NamedDefinitionReference /someOther]]]]]]

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
[ModuleDef $first = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]]]]) -> (Arithmetic [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]]]] FMULTIPLY [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]]]])]]
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
[ModuleDef $first = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Char]]]]]) -> [Integer 2]]]
[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /first]] [[Char 99]]]]]
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
[ModuleDef $first = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference Int/toFixed]] [[FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]]
[ModuleDef $another = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /first]] [[Integer 2]]]]]
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
[ModuleDef $first = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference Int/round]] [[FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]]]]]]]]
[ModuleDef $another = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /first]] [[Fixed 2300]]]]]
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
[ModuleDef $first = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Fixed 300]]]
[ModuleDef $another = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Let [[LetAssign [[LetVar $x]] = [FnCall [FunctionRef [NamedDefinitionReference /first]] [[Integer 2]]]]] in [LetVarRef [LetVar $x]]]]]
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
State : [Alias State [RecordType [[RecordTypeField $playerX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]] => [RecordType [[RecordTypeField $playerX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]

[ModuleDef $another = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Let [[LetAssign [[LetVar $state]] = [RecordLiteral [RecordType [[RecordTypeField $playerX [Primitive Int] (0)]][]] [0 = [Integer 0]]]]] in [RecordLiteral [RecordType [[RecordTypeField $playerX [Primitive Int] (0)]][]] [0 = [Integer 22]]]]]]
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
[ModuleDef $first = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
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
[ModuleDef $first = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
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
[ModuleDef $first = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
[ModuleDef $second = [FunctionValue ([[Arg $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [FunctionParamRef [Arg $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] PLUS [FunctionParamRef [Arg $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
[ModuleDef $third = [FunctionValue ([[Arg $c : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [FunctionParamRef [Arg $c : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MINUS [FunctionParamRef [Arg $c : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]

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
[ModuleDef $first = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
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
Status : [CustomType Status [[Variant $Unknown []] [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]
Unknown : [Variant $Unknown []]
Something : [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]

[ModuleDef $receiveStatus = [FunctionValue ([[Arg $status : [VariantRef [NamedDefTypeRef :[TypeReference $Status]]]]]) -> [FunctionParamRef [Arg $status : [VariantRef [NamedDefTypeRef :[TypeReference $Status]]]]]]]
[ModuleDef $someFunc = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /receiveStatus]] [[VariantConstructor [Variant $Unknown []] []]]]]]

`)
}

func TestCustomTypeVariantAtomLiteral(t *testing.T) {
	testDecorate(t,
		`
type Status =
    Unknown
    | Something Int


type Unrelated =
    ShouldNotMatch
    | SomethingElse Int


receiveStatus : Unknown -> Status
receiveStatus status =
    status


someFunc : String -> Status
someFunc _ =
    receiveStatus Unknown
`,
		`
Status : [CustomType Status [[Variant $Unknown []] [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]
Unknown : [Variant $Unknown []]
Something : [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]
Unrelated : [CustomType Unrelated [[Variant $ShouldNotMatch []] [Variant $SomethingElse [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]
ShouldNotMatch : [Variant $ShouldNotMatch []]
SomethingElse : [Variant $SomethingElse [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]

[ModuleDef $receiveStatus = [FunctionValue ([[Arg $status : [VariantRef [NamedDefTypeRef :[TypeReference $Unknown]] [Variant $Unknown []]]]]) -> [FunctionParamRef [Arg $status : [VariantRef [NamedDefTypeRef :[TypeReference $Unknown]] [Variant $Unknown []]]]]]]
[ModuleDef $someFunc = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /receiveStatus]] [[VariantConstructor [Variant $Unknown []] []]]]]]

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
[ModuleDef $someFunc = [FunctionValue ([[Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [BoolOp [FunctionParamRef [Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] EQ [String Something]]]]
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
someFunc _ =
    Something 42
`, `
Status : [CustomType Status [[Variant $Unknown []] [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]
Unknown : [Variant $Unknown []]
Something : [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]

[ModuleDef $someFunc = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [VariantConstructor [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [[Integer 42]]]]]

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
someFunc a =
    a.another
`, `
Lister : [Alias Lister [RecordType [[RecordTypeField $another [GenericParam u] (0)] [RecordTypeField $fake [GenericParam t] (1)]][[GenericParam t] [GenericParam u]]]] => [RecordType [[RecordTypeField $another [GenericParam u] (0)] [RecordTypeField $fake [GenericParam t] (1)]][[GenericParam t] [GenericParam u]]]
Sprite : [Alias Sprite [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]] => [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]

[ModuleDef $someFunc = [FunctionValue ([[Arg $a : [AliasRef [Alias Lister [RecordType [[RecordTypeField $another [GenericParam u] (0)] [RecordTypeField $fake [GenericParam t] (1)]][[GenericParam t] [GenericParam u]]]]]<[AliasRef [Alias Sprite [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]]],[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]>]]) -> [lookups [FunctionParamRef [Arg $a : [AliasRef [Alias Lister [RecordType [[RecordTypeField $another [GenericParam u] (0)] [RecordTypeField $fake [GenericParam t] (1)]][[GenericParam t] [GenericParam u]]]]]<[AliasRef [Alias Sprite [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]]],[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]>]] [[lookup [RecordTypeField $another [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]]]
`)
}

func TestPipeRight(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
first : Int -> Int
first a =
    a * a


second : String -> Int -> Bool
second _ a =
    a > 25


third : Bool -> Bool
third a =
    a


tester : String -> Bool
tester b =
    first (2 + 2) |> second b |> third
`, `
[ModuleDef $first = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
[ModuleDef $second = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]] [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [BoolOp [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] GR [Integer 25]]]]
[ModuleDef $third = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]]]
[ModuleDef $tester = [FunctionValue ([[Arg $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /first]] [(Arithmetic [Integer 2] PLUS [Integer 2])]] |> [FnCall [FnCall [FunctionRef [NamedDefinitionReference /second]] [[FunctionParamRef [Arg $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] [FnCall [FunctionRef [NamedDefinitionReference /first]] [(Arithmetic [Integer 2] PLUS [Integer 2])]]]] [[FunctionParamRef [Arg $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]]] |> [FnCall [FnCall [FunctionRef [NamedDefinitionReference /third]] [[FnCall [FunctionRef [NamedDefinitionReference /first]] [(Arithmetic [Integer 2] PLUS [Integer 2])]] |> [FnCall [FnCall [FunctionRef [NamedDefinitionReference /second]] [[FunctionParamRef [Arg $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] [FnCall [FunctionRef [NamedDefinitionReference /first]] [(Arithmetic [Integer 2] PLUS [Integer 2])]]]] [[FunctionParamRef [Arg $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]]]]] []]]]
`)
}

func TestBasicStringInterpolation(t *testing.T) {
	testDecorate(t,
		`
sample : Int -> String
sample a =
    $"hello {a}"
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
tester _ =
    first (2 + 2) |> second
`, `
[ModuleDef $first = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
[ModuleDef $second = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [BoolOp [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] GR [Integer 25]]]]
[ModuleDef $tester = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /first]] [(Arithmetic [Integer 2] PLUS [Integer 2])]] |> [FnCall [FnCall [FunctionRef [NamedDefinitionReference /second]] [[FnCall [FunctionRef [NamedDefinitionReference /first]] [(Arithmetic [Integer 2] PLUS [Integer 2])]]]] []]]]
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
tester _ =
    second <| first (2 + 2)
`, `
[ModuleDef $first = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
[ModuleDef $second = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [BoolOp [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] GR [Integer 25]]]]
[ModuleDef $tester = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /second]] []] <| [FnCall [FunctionRef [NamedDefinitionReference /first]] [(Arithmetic [Integer 2] PLUS [Integer 2])]]]]
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
tester _ =
    first (2 + 2)
        |> second
`, `
[ModuleDef $first = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
[ModuleDef $second = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [BoolOp [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] GR [Integer 25]]]]
[ModuleDef $tester = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /first]] [(Arithmetic [Integer 2] PLUS [Integer 2])]] |> [FnCall [FnCall [FunctionRef [NamedDefinitionReference /second]] [[FnCall [FunctionRef [NamedDefinitionReference /first]] [(Arithmetic [Integer 2] PLUS [Integer 2])]]]] []]]]

`)
}

func TestBasicCall(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
move : Int -> Int -> Int
move pos delta =
    pos + delta
`, `
[ModuleDef $move = [FunctionValue ([[Arg $pos : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]] [Arg $delta : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [FunctionParamRef [Arg $pos : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] PLUS [FunctionParamRef [Arg $delta : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
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
create _ =
    { a = 2, b = True }
`, `
Constructor : [Alias Constructor [RecordType [[RecordTypeField $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $b [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]][]]] => [RecordType [[RecordTypeField $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $b [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]][]]

[ModuleDef $create = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [RecordLiteral [RecordType [[RecordTypeField $a [Primitive Int] (0)] [RecordTypeField $b [Primitive Bool] (1)]][]] [0 = [Integer 2] 1 = [Bool true]]]]]
`)
}

func TestConstant2(t *testing.T) {
	testDecorate(t,
		`
tileHeight : Int
tileHeight =
    2


create : Int -> Int
create _ =
    tileHeight
`, `
[ModuleDef $tileHeight = [Constant [Integer 2]]]
[ModuleDef $create = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [ConstantRef [NamedDefinitionReference /tileHeight]]]]
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
a _ =
    First "Hello"
`, `
SomeCustomType : [CustomType SomeCustomType [[Variant $First [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] [Variant $Anon []] [Variant $Second [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]
First : [Variant $First [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]
Anon : [Variant $Anon []]
Second : [Variant $Second [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]

[ModuleDef $a = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [VariantConstructor [Variant $First [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] [[String Hello]]]]]
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
[ModuleDef $move = [FunctionValue ([[Arg $pos : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]] [Arg $delta : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Let [[LetAssign [[LetVar $tenMore]] = (Arithmetic [FunctionParamRef [Arg $pos : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] PLUS [Integer 10])] [LetAssign [[LetVar $tenLess]] = (Arithmetic [FunctionParamRef [Arg $delta : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MINUS [Integer 10])]] in (Arithmetic [LetVarRef [LetVar $tenMore]] PLUS [LetVarRef [LetVar $tenLess]])]]]
`)
}

func TestBasicCallFromVariable(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
move : Int -> Int -> Int
move pos delta =
    pos + delta


main : Bool -> Int
main _ =
    let
        fn = move
    in
    fn 2 3
`, `
[ModuleDef $move = [FunctionValue ([[Arg $pos : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]] [Arg $delta : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [FunctionParamRef [Arg $pos : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] PLUS [FunctionParamRef [Arg $delta : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [Let [[LetAssign [[LetVar $fn]] = [FunctionRef [NamedDefinitionReference /move]]]] in [FnCall [LetVarRef [LetVar $fn]] [[Integer 2] [Integer 3]]]]]]
`)
}

func TestBasicAppend(t *testing.T) {
	testDecorate(t,
		`
a : Int -> List Int
a _ =
    [ 1, 3, 4 ] ++ [ 5, 6, 7, 8 ]
`, `
[ModuleDef $a = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [ListLiteral [[Integer 1] [Integer 3] [Integer 4]]] APPEND [ListLiteral [[Integer 5] [Integer 6] [Integer 7] [Integer 8]]])]]
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
Position : [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]] => [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]
Sprite : [Alias Sprite [RecordType [[RecordTypeField $pos [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (0)]][]]] => [RecordType [[RecordTypeField $pos [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (0)]][]]

[ModuleDef $getx = [FunctionValue ([[Arg $sprite : [AliasRef [Alias Sprite [RecordType [[RecordTypeField $pos [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (0)]][]]]]]]) -> [lookups [FunctionParamRef [Arg $sprite : [AliasRef [Alias Sprite [RecordType [[RecordTypeField $pos [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (0)]][]]]]]] [[lookup [RecordTypeField $pos [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (0)]] [lookup [RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]]]
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
[ModuleDef $isBestAge = [FunctionValue ([[Arg $age : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [If [Logical [BoolOp [FunctionParamRef [Arg $age : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] EQ [Integer 50]] or [BoolOp [FunctionParamRef [Arg $age : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] GRE [Integer 65]]] then [Integer 100] else [Integer 0]]]]
`)
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
[ModuleDef $isBeautiful = [FunctionValue ([[Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [If [BoolOp [FunctionParamRef [Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] EQ [String Rebecca]] then [Bool true] else [Bool false]]]]
`)
}

func TestGuard(t *testing.T) {
	testDecorate(t,
		`
isBeautiful : String -> Bool
isBeautiful name =
    | name == "Rebecca" -> True
    | _ -> False
`, `
[ModuleDef $isBeautiful = [FunctionValue ([[Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [DGuard: [DGuardItem [BoolOp [FunctionParamRef [Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] EQ [String Rebecca]] [Bool true]] default: [Bool false]]]]
`)
}

func TestBoolPerson(t *testing.T) {
	testDecorate(t,
		`
isBeautiful : String -> Bool
isBeautiful name =
    name == "Rebecca"
`, `
[ModuleDef $isBeautiful = [FunctionValue ([[Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [BoolOp [FunctionParamRef [Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] EQ [String Rebecca]]]]
`)
}

func TestBoolPersonCall(t *testing.T) {
	testDecorate(t,
		`
isLoveOfMyLife : String -> Bool
isLoveOfMyLife name =
    name == "Rebecca"


main : String -> Bool
main _ =
    isLoveOfMyLife "Lisa"
`, `
[ModuleDef $isLoveOfMyLife = [FunctionValue ([[Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [BoolOp [FunctionParamRef [Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] EQ [String Rebecca]]]]
[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /isLoveOfMyLife]] [[String Lisa]]]]]
`)
}

func TestBoolPersonCallAgain(t *testing.T) {
	testDecorate(t,
		`
isLoveOfMyLife : String -> Int -> Bool
isLoveOfMyLife name _ =
    name == "Rebecca"


main : String -> Bool
main _ =
    isLoveOfMyLife "Lisa" 2
`, `
[ModuleDef $isLoveOfMyLife = [FunctionValue ([[Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]] [Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [BoolOp [FunctionParamRef [Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] EQ [String Rebecca]]]]
[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /isLoveOfMyLife]] [[String Lisa] [Integer 2]]]]]
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
Chore : [CustomType Chore [[Variant $Meeting [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] [Variant $Running [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [Variant $Unknown []]]]
Meeting : [Variant $Meeting [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]
Running : [Variant $Running [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]
Unknown : [Variant $Unknown []]
`)
}

func TestCustomTypeGenerics(t *testing.T) {
	testDecorate(t,
		`
type Perhaps a =
    None
    | Actual a
`, `
Perhaps : [CustomType Perhaps<a> [[Variant $None []] [Variant $Actual [[GenericParam a]]]]]
None : [Variant $None []]
Actual : [Variant $Actual [[GenericParam a]]]
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
Tinkering : [Alias Tinkering [RecordType [[RecordTypeField $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (0)]][]]] => [RecordType [[RecordTypeField $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (0)]][]]
Studying : [Alias Studying [RecordType [[RecordTypeField $subject [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]][]]] => [RecordType [[RecordTypeField $subject [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]][]]
Work : [Alias Work [RecordType [[RecordTypeField $web [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]] => [RecordType [[RecordTypeField $web [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]
Child : [CustomType Child [[Variant $Aron [[AliasRef [Alias Tinkering [RecordType [[RecordTypeField $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (0)]][]]]]]] [Variant $Alexandra []] [Variant $Alma [[AliasRef [Alias Studying [RecordType [[RecordTypeField $subject [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]][]]]]]] [Variant $Isabelle [[AliasRef [Alias Work [RecordType [[RecordTypeField $web [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]]]]]]]
Aron : [Variant $Aron [[AliasRef [Alias Tinkering [RecordType [[RecordTypeField $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (0)]][]]]]]]
Alexandra : [Variant $Alexandra []]
Alma : [Variant $Alma [[AliasRef [Alias Studying [RecordType [[RecordTypeField $subject [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]][]]]]]]
Isabelle : [Variant $Isabelle [[AliasRef [Alias Work [RecordType [[RecordTypeField $web [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]]]]]
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
[ModuleDef $another = [FunctionValue ([[Arg $score : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Let [[LetAssign [[LetVar $af]] = [Curry [FunctionRef [NamedDefinitionReference /f]] [[String Peter]]]]] in [FnCall [LetVarRef [LetVar $af]] [[FunctionParamRef [Arg $score : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]]]
[ModuleDef $f = [FunctionValue ([[Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]] [Arg $score : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [If [BoolOp [FunctionParamRef [Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] EQ [String Peter]] then [BoolOp (arithmetic [FunctionParamRef [Arg $score : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [Integer 2]) GR [Integer 100]] else [BoolOp [FunctionParamRef [Arg $score : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] GR [Integer 100]]]]]
`)
}

func TestBlob(t *testing.T) {
	testDecorate(t,
		`
a : Blob -> List Int
a _ =
    [ 10, 20, 99 ]
`, `
[ModuleDef $a = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Blob]]]]]) -> [ListLiteral [[Integer 10] [Integer 20] [Integer 99]]]]]
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
a _ =
    (2, "Hello")
`, `
[ModuleDef $a = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [TupleLiteral [TupleType [[Primitive Int] [Primitive String]]] [TupleLiteral [#2 'Hello']] [[Integer 2] [String Hello]]]]]
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

func TestTupleGenerics(t *testing.T) {
	testDecorate(t,
		`
createTuple : a -> b -> (a, b)
createTuple first second =
    (first, second)


a : Bool -> (Int, String)
a _ =
    createTuple 2 "Hello"
`, `
[ModuleDef $createTuple = [FunctionValue ([[Arg $first : [GenericParam a]] [Arg $second : [GenericParam b]]]) -> [TupleLiteral [TupleType [[GenericParam a] [GenericParam b]]] [TupleLiteral [$first $second]] [[FunctionParamRef [Arg $first : [GenericParam a]]] [FunctionParamRef [Arg $second : [GenericParam b]]]]]]]
[ModuleDef $a = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /createTuple]] [[Integer 2] [String Hello]]]]]
`)
}

func TestTupleGenericsFail(t *testing.T) {
	testDecorateFail(t,
		`
createTuple : a -> b -> (a, b)
createTuple first second =
    (first, second)


a : Bool -> (Int, String)
a _ =
    createTuple "2" "Hello"
`, &decorated.UnMatchingFunctionReturnTypesInFunctionValue{})
}

func TestArrayLiteral2(t *testing.T) {
	testDecorate(t,
		`
a : Bool -> Array Int
a _ =
    [| 10, 20, 99 |]
`, `
[ModuleDef $a = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [ArrayLiteral Array [[Integer 10] [Integer 20] [Integer 99]]]]]
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
a _ =
    [ Cool { name = 95 } ]
`, `
Cool : [Alias Cool [RecordType [[RecordTypeField $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]] => [RecordType [[RecordTypeField $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]

[ModuleDef $a = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [ListLiteral [[RecordConstructorRecord [CCall [TypeReference $Cool] [[RecordLiteral [[$name = #95]]]]] [RecordLiteral [RecordType [[RecordTypeField $name [Primitive Int] (0)]][]] [0 = [Integer 95]]]]]]]]
`)
}

func TestRecordConstructorWithoutSpace(t *testing.T) {
	testDecorate(t,
		`
type alias Cool =
    { name : Int
    }


a : Bool -> List Cool
a _ =
    [ Cool{ name = 95 } ]
`, `
Cool : [Alias Cool [RecordType [[RecordTypeField $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]] => [RecordType [[RecordTypeField $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]

[ModuleDef $a = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [ListLiteral [[RecordConstructorRecord [CCall [TypeReference $Cool] [[RecordLiteral [[$name = #95]]]]] [RecordLiteral [RecordType [[RecordTypeField $name [Primitive Int] (0)]][]] [0 = [Integer 95]]]]]]]]
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

func TestRecordGenerics(t *testing.T) {
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
Tinkering : [Alias Tinkering [RecordType [[RecordTypeField $secret [GenericParam t] (0)] [RecordTypeField $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]][[GenericParam t]]]] => [RecordType [[RecordTypeField $secret [GenericParam t] (0)] [RecordTypeField $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]][[GenericParam t]]]

[ModuleDef $f = [FunctionValue ([[Arg $tinkering : [AliasRef [Alias Tinkering [RecordType [[RecordTypeField $secret [GenericParam t] (0)] [RecordTypeField $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]][[GenericParam t]]]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]>]]) -> [lookups [FunctionParamRef [Arg $tinkering : [AliasRef [Alias Tinkering [RecordType [[RecordTypeField $secret [GenericParam t] (0)] [RecordTypeField $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]][[GenericParam t]]]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]>]] [[lookup [RecordTypeField $secret [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]]]
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
Tinkering : [Alias Tinkering [RecordType [[RecordTypeField $secret [GenericParam t] (0)] [RecordTypeField $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]][[GenericParam t]]]] => [RecordType [[RecordTypeField $secret [GenericParam t] (0)] [RecordTypeField $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]][[GenericParam t]]]
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
drawSprite _ =
    True


drawSprites : List Sprite -> List Bool
drawSprites sprites =
    List.map drawSprite sprites


drawWorld : World -> List (List Bool)
drawWorld world =
    List.map drawSprites world.drawTasks


main : Bool -> List (List Bool)
main _ =
    drawWorld { drawTasks = [ [ { x = 10, y = 20 }, { x = 44, y = 98 } ], [ { x = 99, y = 98 } ] ] }
`, `
Sprite : [Alias Sprite [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]] => [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]
World : [Alias World [RecordType [[RecordTypeField $drawTasks [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[AliasRef [Alias Sprite [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]>> (0)]][]]] => [RecordType [[RecordTypeField $drawTasks [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[AliasRef [Alias Sprite [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]>> (0)]][]]

[ModuleDef $drawSprite = [FunctionValue ([[Arg $_ : [AliasRef [Alias Sprite [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]]]) -> [Bool true]]]
[ModuleDef $drawSprites = [FunctionValue ([[Arg $sprites : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[AliasRef [Alias Sprite [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]>]]) -> [FnCall [FunctionRef [NamedDefinitionReference List/map]] [[FunctionRef [NamedDefinitionReference /drawSprite]] [FunctionParamRef [Arg $sprites : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[AliasRef [Alias Sprite [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]>]]]]]]
[ModuleDef $drawWorld = [FunctionValue ([[Arg $world : [AliasRef [Alias World [RecordType [[RecordTypeField $drawTasks [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[AliasRef [Alias Sprite [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]>> (0)]][]]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference List/map]] [[FunctionRef [NamedDefinitionReference /drawSprites]] [lookups [FunctionParamRef [Arg $world : [AliasRef [Alias World [RecordType [[RecordTypeField $drawTasks [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[AliasRef [Alias Sprite [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]>> (0)]][]]]]]] [[lookup [RecordTypeField $drawTasks [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[AliasRef [Alias Sprite [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]>> (0)]]]]]]]]
[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /drawWorld]] [[record-literal [RecordType [[RecordTypeField $drawTasks [Primitive List<[Primitive List<[RecordType [[RecordTypeField $x [Primitive Int] (0)] [RecordTypeField $y [Primitive Int] (1)]][]]>]>] (0)]][]] [0 = [ListLiteral [[ListLiteral [[record-literal [RecordType [[RecordTypeField $x [Primitive Int] (0)] [RecordTypeField $y [Primitive Int] (1)]][]] [0 = [Integer 10] 1 = [Integer 20]]] [record-literal [RecordType [[RecordTypeField $x [Primitive Int] (0)] [RecordTypeField $y [Primitive Int] (1)]][]] [0 = [Integer 44] 1 = [Integer 98]]]]] [ListLiteral [[record-literal [RecordType [[RecordTypeField $x [Primitive Int] (0)] [RecordTypeField $y [Primitive Int] (1)]][]] [0 = [Integer 99] 1 = [Integer 98]]]]]]]]]]]]]
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
MyList : [Alias MyList [RecordType [[RecordTypeField $someType [GenericParam t] (0)]][[GenericParam t]]]] => [RecordType [[RecordTypeField $someType [GenericParam t] (0)]][[GenericParam t]]]
`)
}

func TestAppliedAnnotation2(t *testing.T) {
	testDecorateWithoutDefault(t,

		`
type alias MyList t =
    { someType : t
    }


intConvert : MyList Int -> Bool
intConvert _ =
    True
`, `
MyList : [Alias MyList [RecordType [[RecordTypeField $someType [GenericParam t] (0)]][[GenericParam t]]]] => [RecordType [[RecordTypeField $someType [GenericParam t] (0)]][[GenericParam t]]]

[ModuleDef $intConvert = [FunctionValue ([[Arg $_ : [AliasRef [Alias MyList [RecordType [[RecordTypeField $someType [GenericParam t] (0)]][[GenericParam t]]]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]>]]) -> [Bool true]]]
`)
}

func TestBasicFunctionTypeParameters(t *testing.T) {
	testDecorate(t,

		`
simple : List a -> a
simple _ =
    2


main : Bool -> Int
main _ =
    simple [ 1, 2, 3 ]
`, `
[ModuleDef $simple = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[GenericParam a]>]]) -> [Integer 2]]]
[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /simple]] [[ListLiteral [[Integer 1] [Integer 2] [Integer 3]]]]]]]
`)
}

func TestFunctionTypeParameters(t *testing.T) {
	testDecorate(t,

		`
intConvert : List Int -> Bool
intConvert _ =
    True


simple : (List a -> Bool) -> List a -> Bool
simple _ _ =
    True


main : Bool -> Bool
main _ =
    simple intConvert [ 1, 2, 3 ]
`, `
[ModuleDef $intConvert = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]>]]) -> [Bool true]]]
[ModuleDef $simple = [FunctionValue ([[Arg $_ : [FunctionTypeRef [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[GenericParam a]> [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]] [Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[GenericParam a]>]]) -> [Bool true]]]
[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /simple]] [[FunctionRef [NamedDefinitionReference /intConvert]] [ListLiteral [[Integer 1] [Integer 2] [Integer 3]]]]]]]
`)
}

func TestFunctionTypeGenerics(t *testing.T) {
	testDecorate(t,

		`
intConvert : Int -> Int
intConvert _ =
    42


simple : a -> (a -> a)
simple _ =
    intConvert


main : Bool -> Int
main _ =
    let
        fn = simple 33
    in
    intConvert (fn 22)
`, `
[ModuleDef $intConvert = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Integer 42]]]
[ModuleDef $simple = [FunctionValue ([[Arg $_ : [GenericParam a]]]) -> [FunctionRef [NamedDefinitionReference /intConvert]]]]
[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [Let [[LetAssign [[LetVar $fn]] = [FnCall [FunctionRef [NamedDefinitionReference /simple]] [[Integer 33]]]]] in [FnCall [FunctionRef [NamedDefinitionReference /intConvert]] [[FnCall [LetVarRef [LetVar $fn]] [[Integer 22]]]]]]]]
`)
}

func TestFunctionTypeGenericsFail(t *testing.T) {
	testDecorateFail(t,

		`
intConvert : Int -> Int
intConvert _ =
    42


simple : a -> (a -> a)
simple _ =
    intConvert


main : Bool -> Int
main _ =
    let
        fn = simple "hello"
    in
    -- This should fail, since fn takes a String instead of an Int
    intConvert (fn 22)
`, &decorated.FunctionCallTypeMismatch{})
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
Scale2 : [Alias Scale2 [RecordType [[RecordTypeField $scaleX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $scaleY [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]] => [RecordType [[RecordTypeField $scaleX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $scaleY [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]
Sprite : [Alias Sprite [RecordType [[RecordTypeField $dummy [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $scale [AliasRef [Alias Scale2 [RecordType [[RecordTypeField $scaleX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $scaleY [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (1)]][]]] => [RecordType [[RecordTypeField $dummy [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $scale [AliasRef [Alias Scale2 [RecordType [[RecordTypeField $scaleX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $scaleY [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (1)]][]]

[ModuleDef $updateSprite = [FunctionValue ([[Arg $inSprite : [AliasRef [Alias Sprite [RecordType [[RecordTypeField $dummy [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $scale [AliasRef [Alias Scale2 [RecordType [[RecordTypeField $scaleX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $scaleY [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (1)]][]]]]] [Arg $newScale : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [RecordLiteral [RecordType [[RecordTypeField $dummy [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $scale [AliasRef [Alias Scale2 [RecordType [[RecordTypeField $scaleX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $scaleY [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (1)]][]] [1 = [RecordLiteral [RecordType [[RecordTypeField $scaleX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $scaleY [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]] [0 = [FunctionParamRef [Arg $newScale : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] 1 = [FunctionParamRef [Arg $newScale : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]]]]
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
[ModuleDef $checkMaybeInt = [FunctionValue ([[Arg $a : [VariantRef [NamedDefTypeRef :[TypeReference $Maybe]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]>] [Arg $b : [VariantRef [NamedDefTypeRef :[TypeReference $Maybe]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]>]]) -> [dcase: [FunctionParamRef [Arg $a : [VariantRef [NamedDefTypeRef :[TypeReference $Maybe]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]>]] of [dcasecons [VariantRef [NamedDefTypeRef :[TypeReference $Just]] [Variant $Just [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]] ([[dcaseparm $firstInt type:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) => [dcase: [FunctionParamRef [Arg $b : [VariantRef [NamedDefTypeRef :[TypeReference $Maybe]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]>]] of [dcasecons [VariantRef [NamedDefTypeRef :[TypeReference $Nothing]] [Variant $Nothing []]] ([]) => [Integer 0]];[dcasecons [VariantRef [NamedDefTypeRef :[TypeReference $Just]] [Variant $Just [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]] ([[dcaseparm $secondInt type:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) => [functionparamref $secondInt [dcaseparm $secondInt type:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]];[dcasecons [VariantRef [NamedDefTypeRef :[TypeReference $Nothing]] [Variant $Nothing []]] ([]) => [Integer 0]]]]]
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
[ModuleDef $single = [Constant [Integer 2]]]
    `)
}

func TestMinimalMaybeAndNothing(t *testing.T) {
	testDecorate(t,
		`
type alias Thing =
    { something : String
    }


main : Bool -> Maybe Thing
main _ =
    Nothing
`, `
Thing : [Alias Thing [RecordType [[RecordTypeField $something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]][]]] => [RecordType [[RecordTypeField $something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]][]]

[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [VariantConstructor [Variant $Nothing []] []]]]
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
main _ =
    Just { something = "hello" }
`, `
Thing : [Alias Thing [RecordType [[RecordTypeField $something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]][]]] => [RecordType [[RecordTypeField $something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]][]]

[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [VariantConstructor [Variant $Just [[GenericParam a]]] [[RecordLiteral [RecordType [[RecordTypeField $something [Primitive String] (0)]][]] [0 = [String hello]]]]]]]
`)
}

func TestJustAndNothing(t *testing.T) {
	testDecorate(t,
		`
type alias Thing =
    { something : String
    }


main : Bool -> Maybe Thing
main _ =
    if True then
        Nothing
    else
        Just { something = "hello" }
`, `
Thing : [Alias Thing [RecordType [[RecordTypeField $something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]][]]] => [RecordType [[RecordTypeField $something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]][]]

[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [If [Bool true] then [VariantConstructor [Variant $Nothing []] []] else [VariantConstructor [Variant $Just [[GenericParam a]]] [[RecordLiteral [RecordType [[RecordTypeField $something [Primitive String] (0)]][]] [0 = [String hello]]]]]]]]
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
checkMaybeGamepad _ =
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
needInputs _ =
    2


a : Bool -> Int
a _ =
    needInputs (Array.fromList [])
`,
		`
[ModuleDef $needInputs = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Array]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]>]]) -> [Integer 2]]]
[ModuleDef $a = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /needInputs]] [[FnCall [FunctionRef [NamedDefinitionReference Array/fromList]] [[ListLiteral []]]]]]]]
`)
}

func TestReturnFunction(t *testing.T) {
	testDecorate(t,
		`
needInputs : Array Int -> Int
needInputs _ =
    2


returnSomeFunc : Bool -> (Array Int -> Int)
returnSomeFunc _ =
    needInputs


a : Bool -> Int
a ignore =
    let
        fn = returnSomeFunc ignore
    in
    fn (Array.fromList [])
`,
		`
[ModuleDef $needInputs = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Array]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]>]]) -> [Integer 2]]]
[ModuleDef $returnSomeFunc = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FunctionRef [NamedDefinitionReference /needInputs]]]]
[ModuleDef $a = [FunctionValue ([[Arg $ignore : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [Let [[LetAssign [[LetVar $fn]] = [FnCall [FunctionRef [NamedDefinitionReference /returnSomeFunc]] [[FunctionParamRef [Arg $ignore : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]]]]] in [FnCall [LetVarRef [LetVar $fn]] [[FnCall [FunctionRef [NamedDefinitionReference Array/fromList]] [[ListLiteral []]]]]]]]]
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
Something : [Alias Something [RecordType [[RecordTypeField $playerX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $time [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]] => [RecordType [[RecordTypeField $playerX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $time [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]

[ModuleDef $fn = [FunctionValue ([[Arg $something : [AliasRef [Alias Something [RecordType [[RecordTypeField $playerX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $time [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]]]) -> [RecordLiteral [RecordType [[RecordTypeField $playerX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $time [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]] [0 = [Integer 3]]]]]
`)
}

func TestPipeRight2(t *testing.T) {
	testDecorate(t, `
first : Int -> Int
first _ =
    2


second : String -> Int -> Int
second _ _ =
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
[ModuleDef $first = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Integer 2]]]
[ModuleDef $second = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]] [Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Integer 2]]]
[ModuleDef $third = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [BoolOp [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] GR [Integer 45]]]]
[ModuleDef $tester = [FunctionValue ([[Arg $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> (Arithmetic [FnCall [FunctionRef [NamedDefinitionReference /first]] [[Integer 2]]] PLUS [Integer 2]) |> [FnCall [FnCall [FunctionRef [NamedDefinitionReference /second]] [[FunctionParamRef [Arg $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] (Arithmetic [FnCall [FunctionRef [NamedDefinitionReference /first]] [[Integer 2]]] PLUS [Integer 2])]] [[FunctionParamRef [Arg $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]]] |> [FnCall [FnCall [FunctionRef [NamedDefinitionReference /third]] [(Arithmetic [FnCall [FunctionRef [NamedDefinitionReference /first]] [[Integer 2]]] PLUS [Integer 2]) |> [FnCall [FnCall [FunctionRef [NamedDefinitionReference /second]] [[FunctionParamRef [Arg $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] (Arithmetic [FnCall [FunctionRef [NamedDefinitionReference /first]] [[Integer 2]]] PLUS [Integer 2])]] [[FunctionParamRef [Arg $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]]]]] []]]]
`)
}
