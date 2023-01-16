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
first : (a: Int) -> Int =
    a * a
`,
		`
[ModuleDef $first = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
`)
}

func TestNewFunction(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
--| ignore this
first : (startNumber: Int) -> Int =
    startNumber * startNumber
`,
		`
[ModuleDef $first = [FunctionValue ([[Arg $startNumber : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [FunctionParamRef [Arg $startNumber : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [FunctionParamRef [Arg $startNumber : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
`)
}

func TestConstant(t *testing.T) {
	testDecorateWithoutDefault(t, `
fn =
    "Hello"

another : String =
    fn


`, `
[ModuleDef $fn = [Constant [String Hello]]]
[ModuleDef $another = [FunctionValue ([]) -> [ConstantRef [NamedDefinitionReference /fn]]]]
`)
}

func TestBooleanLookup(t *testing.T) {
	testDecorateWithoutDefault(t, `
another : Bool =
    let
        a = True
        b = False
    in
    a && b
`, `
[ModuleDef $another = [FunctionValue ([]) -> [Let [[LetAssign [[LetVar $a]] = [Bool true]] [LetAssign [[LetVar $b]] = [Bool false]]] in [Logical [LetVarRef [LetVar $a]] and [LetVarRef [LetVar $b]]]]]]`)
}

func TestCast(t *testing.T) {
	testDecorateWithoutDefault(t, `
type alias Something = Int


another : Bool =
    let
        b = 32 : Something
    in
    b >= 32
`, `
Something : [Alias Something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [Primitive Int]

[ModuleDef $another = [FunctionValue ([]) -> [Let [[LetAssign [[LetVar $b]] = [Cast [Integer 32] [AliasRefExpr [AliasRef [Alias Something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]] in [BoolOp [LetVarRef [LetVar $b]] GRE [Integer 32]]]]]
`)
}

func TestFunctionTypeWithoutParameters(t *testing.T) {
	testDecorateWithoutDefault(t, `__externalvarfn another : Int -> String -> Int

`, `
Something : [Alias Something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [Primitive Int]

[ModuleDef $another = [FunctionValue ([]) -> [Let [[LetAssign [[LetVar $b]] = [Cast [Integer 32] [AliasRefExpr [AliasRef [Alias Something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]] in [BoolOp [LetVarRef [LetVar $b]] GRE [Integer 32]]]]]
`)
}

func TestResourceName(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
--| ignore this
first : (ResourceName) -> Int =
    2


main : (Bool) -> Int =
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
fn : (*) -> Int =
    23
`,
		`
[ModuleDef $fn = [FunctionValue ([[Arg $_ : [AnyMatching AnyMatchingType]]]) -> [Integer 23]]]
`)
}

func TestAnyMatchingTypeCall(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
fn : (*) -> Int =
    23


main : Int =
    fn 23.0 "hello" 42
`,
		`
[ModuleDef $fn = [FunctionValue ([[Arg $_ : [AnyMatching AnyMatchingType]]]) -> [Integer 23]]]
[ModuleDef $main = [FunctionValue ([]) -> [FnCall [FunctionRef [NamedDefinitionReference /fn]] [[Fixed 23000] [String hello] [Integer 42]]]]]
`)
}

func TestAnyMatchingTypeCallMiddle(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
fn : (String, *) -> Int =
    23


main : Int =
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
fn : (String, *, a) -> List a =
    [ 23.0 ]


main : List Fixed =
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
fn : (String, (* -> a)) -> List a =
    [ 23.0 ]


someOther : Int -> Fixed =
    3.5


main : List Fixed =
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
first : (a: Fixed) -> Fixed =
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
first : (Char) -> Int =
    2


main : (Bool) -> Int =
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
first : (a: Int) -> Fixed =
    Int.toFixed a


another : (Int) -> Fixed =
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
first : (a: Fixed) -> Int =
    Int.round a


another : (_: Int) -> Int =
    first 2.3
`,
		`
[ModuleDef $first = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference Int/round]] [[FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]]]]]]]]
[ModuleDef $another = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /first]] [[Fixed 2300]]]]]
`)
}

func TestFixedConvertLet(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
--| ignore this
first : (Int) -> Fixed =
    0.3


another : (Int) -> Fixed =
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
	testDecorateWithoutDefault(t,
		`
type alias State =
    { playerX : Int
    }


another : (Int) -> State =
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
first : (a: Int) -> Int =
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
first : (a: Int) -> Int =
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
first : (a: Int) -> Int =
    a * a


{-
   multiline

-}
second : (b: Int) -> Int =
    b + b


-- single line
third : (c: Int) -> Int =
    c - c
`,
		`
[ModuleDef $first = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
[ModuleDef $second = [FunctionValue ([[Arg $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [FunctionParamRef [Arg $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] PLUS [FunctionParamRef [Arg $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
[ModuleDef $third = [FunctionValue ([[Arg $c : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [FunctionParamRef [Arg $c : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MINUS [FunctionParamRef [Arg $c : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]

`)
}

func TestWrongSpacingMultiDocFail(t *testing.T) {
	testDecorateFail(t,
		`
{-|
   ignore this
      for sure

   -}
first : (a: Int) -> Int =
    a * a
{-
   something else

   -}
second : (b: Int) -> Int =
    b + b
`, parerr.ExpectedTwoLinesAfterStatement{})
}

func TestSimpleCallParameterErrFail(t *testing.T) {
	testDecorateWithoutDefaultFail(t,
		`
someFunc : (Int) -> String =
    "2"


another : (Int) -> String =
    someFunc "2"
`,
		&decorated.CouldNotSmashFunctions{})
}

func TestSomething2(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
first : (a: Int) -> Int =
    a * a
`,
		`
[ModuleDef $first = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
`)
}

func TestCustomTypeVariantLiteral(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
type Status =
    Unknown
    | Something Int


receiveStatus : (status: Status) -> Status =
    status


someFunc : (String) -> Status =
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
	testDecorateWithoutDefault(t,
		`
type Status =
    Unknown
    | Something Int


type Unrelated =
    ShouldNotMatch
    | SomethingElse Int


receiveStatus : (status: Unknown) -> Status =
    status


someFunc : (String) -> Status =
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
	testDecorateWithoutDefault(t,
		`
someFunc : (name: String) -> Bool =
    name == "Something"
`,
		`
[ModuleDef $someFunc = [FunctionValue ([[Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [BoolOp [FunctionParamRef [Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] EQ [String Something]]]]
`)
}

func TestUnknownAnnotationType(t *testing.T) {
	testDecorateWithoutDefaultFail(t, `
    someFunc : Position2
    `, &decorated.UnknownAnnotationTypeReference{})
}

func TestCustomTypeVariantLiteralWithParameters(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
type Status =
    Unknown
    | Something Int


someFunc : (String) -> Status =
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


someFunc : (a: Lister Sprite Int) -> Int =
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
first : (a: Int) -> Int =
    a * a


second : (String, a: Int) -> Bool =
    a > 25


third : (a: Bool) -> Bool =
    a


tester : (b: String) -> Bool =
    first (2 + 2) |> second b |> third
`, `
[ModuleDef $first = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
[ModuleDef $second = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]] [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [BoolOp [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] GR [Integer 25]]]]
[ModuleDef $third = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]]]
[ModuleDef $tester = [FunctionValue ([[Arg $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /first]] [(Arithmetic [Integer 2] PLUS [Integer 2])]] |> [FnCall [FnCall [FunctionRef [NamedDefinitionReference /second]] [[FunctionParamRef [Arg $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] [FnCall [FunctionRef [NamedDefinitionReference /first]] [(Arithmetic [Integer 2] PLUS [Integer 2])]]]] [[FunctionParamRef [Arg $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]]] |> [FnCall [FnCall [FunctionRef [NamedDefinitionReference /third]] [[FnCall [FunctionRef [NamedDefinitionReference /first]] [(Arithmetic [Integer 2] PLUS [Integer 2])]] |> [FnCall [FnCall [FunctionRef [NamedDefinitionReference /second]] [[FunctionParamRef [Arg $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] [FnCall [FunctionRef [NamedDefinitionReference /first]] [(Arithmetic [Integer 2] PLUS [Integer 2])]]]] [[FunctionParamRef [Arg $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]]]]] []]]]
`)
}

func xTestBasicStringInterpolation(t *testing.T) {
	testDecorate(t,
		`
sample : (a: Int) -> String =
    $"hello {a}"
`, `
func(Int -> String) : [func  [primitive Int] [primitive String]]

sample = [functionvalue ([[arg $a = [primitive Int]]]) -> (arithmetic [str hello ] APPEND [fcall [getvar $Debug.toString [primitive String]] [[getvar $a [primitive Int]]]])]
`)
}

func TestArrayVsListFail(t *testing.T) {
	testDecorateFail(t,
		`
updater : Int -> String -> Int =
    42


sample : Int -> List Int =
    let
        intArray = Array.fromList [ 0, 1, 2 ]

        arraySlice = Array.slice 0 2 intArray
    in
    List.map2 updater arraySlice [ "hello", "world", "fail" ]
`, &decorated.CouldNotSmashFunctions{})
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


move : (pos: Position, delta: Position) -> Position =
    let
        newX = pos.x + delta.x

        newY = pos.y + delta.y
    in
    { x = newX
    , y = newY
    }


moveSprite : (sprite: Sprite, delta: Position) -> Sprite =
    { rootPosition = move sprite.rootPosition delta }
`, `
Position : [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]] => [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]
Sprite : [Alias Sprite [RecordType [[RecordTypeField $rootPosition [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (0)]][]]] => [RecordType [[RecordTypeField $rootPosition [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (0)]][]]

[ModuleDef $move = [FunctionValue ([[Arg $pos : [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]] [Arg $delta : [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]]]) -> [Let [[LetAssign [[LetVar $newX]] = (Arithmetic [lookups [FunctionParamRef [Arg $pos : [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]]] [[lookup [RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]] PLUS [lookups [FunctionParamRef [Arg $delta : [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]]] [[lookup [RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]])] [LetAssign [[LetVar $newY]] = (Arithmetic [lookups [FunctionParamRef [Arg $pos : [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]]] [[lookup [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]] PLUS [lookups [FunctionParamRef [Arg $delta : [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]]] [[lookup [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]])]] in [RecordLiteral [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]] [0 = [LetVarRef [LetVar $newX]] 1 = [LetVarRef [LetVar $newY]]]]]]]
[ModuleDef $moveSprite = [FunctionValue ([[Arg $sprite : [AliasRef [Alias Sprite [RecordType [[RecordTypeField $rootPosition [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (0)]][]]]]] [Arg $delta : [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]]]) -> [RecordLiteral [RecordType [[RecordTypeField $rootPosition [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (0)]][]] [0 = [FnCall [FunctionRef [NamedDefinitionReference /move]] [[lookups [FunctionParamRef [Arg $sprite : [AliasRef [Alias Sprite [RecordType [[RecordTypeField $rootPosition [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (0)]][]]]]]] [[lookup [RecordTypeField $rootPosition [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (0)]]]] [FunctionParamRef [Arg $delta : [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]]]]]]]]]
`)
}

func TestOperatorPipeRight(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
first : (a: Int) -> Int =
    a * a


second : (a: Int) -> Bool =
    a > 25


tester : (String) -> Bool =
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
first : (a: Int) -> Int =
    a * a


second : (a: Int) -> Bool =
    a > 25


tester : (String) -> Bool =
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
first : (a: Int) -> Int =
    a * a


second : (a: Int) -> Bool =
    a > 25


tester : (String) -> Bool =
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
move : (pos: Int, delta: Int) -> Int =
    pos + delta
`, `
[ModuleDef $move = [FunctionValue ([[Arg $pos : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]] [Arg $delta : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [FunctionParamRef [Arg $pos : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] PLUS [FunctionParamRef [Arg $delta : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
`)
}

func TestBasicConstructor(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
type alias Constructor =
    { a : Int
    , b : Bool
    }


create : (Int) -> Constructor =
    { a = 2, b = True }
`, `
Constructor : [Alias Constructor [RecordType [[RecordTypeField $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $b [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]][]]] => [RecordType [[RecordTypeField $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $b [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]][]]

[ModuleDef $create = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [RecordLiteral [RecordType [[RecordTypeField $a [Primitive Int] (0)] [RecordTypeField $b [Primitive Bool] (1)]][]] [0 = [Integer 2] 1 = [Bool true]]]]]
`)
}

func TestConstant2(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
tileHeight : Int =
    2


create : (Int) -> Int =
    tileHeight
`, `
[ModuleDef $tileHeight = [Constant [Integer 2]]]
[ModuleDef $create = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [ConstantRef [NamedDefinitionReference /tileHeight]]]]
`)
}

func TestCustomTypeConstructor(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
type SomeCustomType =
    First String
    | Anon
    | Second Int


a : (Bool) -> SomeCustomType =
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
	testDecorateWithoutDefaultFail(t,
		`
move : (pos: Int, delta: String) -> Int =
    pos + delta
`, &decorated.UnMatchingBinaryOperatorTypes{})
}

func TestBasicCallLet(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
move : (pos: Int, delta: Int) -> Int =
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
move : (pos: Int, delta: Int) -> Int =
    pos + delta


main : (Bool) -> Int =
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
	testDecorateWithoutDefault(t,
		`
a : (Int) -> List Int =
    [ 1, 3, 4 ] ++ [ 5, 6, 7, 8 ]
`, `
[ModuleDef $a = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [ListLiteral [[Integer 1] [Integer 3] [Integer 4]]] APPEND [ListLiteral [[Integer 5] [Integer 6] [Integer 7] [Integer 8]]])]]
`)
}

func TestBasicCallLetFail(t *testing.T) {
	testDecorateWithoutDefaultFail(t,
		`
move : (pos: Int, delta: Int) -> String =
    let
        tenMore = pos + 10

        tenLess = delta - 10
    in
    tenMore + tenLess
`, &decorated.UnMatchingFunctionReturnTypesInFunctionValue{})
}

func TestBasicLookup(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
type alias Position =
    { x : Int
    , y : Int
    }


getx : (pos: Position) -> Int =
    pos.x + pos.y
`, `
Position : [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]] => [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]

[ModuleDef $getx = [FunctionValue ([[Arg $pos : [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]]]) -> (Arithmetic [lookups [FunctionParamRef [Arg $pos : [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]]] [[lookup [RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]] PLUS [lookups [FunctionParamRef [Arg $pos : [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]]] [[lookup [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]])]]
`)
}

func TestSubLookup(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
type alias Position =
    { x : Int
    , y : Int
    }


type alias Sprite =
    { pos : Position
    }


getx : (sprite: Sprite) -> Int =
    sprite.pos.x
`, `
Position : [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]] => [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]
Sprite : [Alias Sprite [RecordType [[RecordTypeField $pos [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (0)]][]]] => [RecordType [[RecordTypeField $pos [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (0)]][]]

[ModuleDef $getx = [FunctionValue ([[Arg $sprite : [AliasRef [Alias Sprite [RecordType [[RecordTypeField $pos [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (0)]][]]]]]]) -> [lookups [FunctionParamRef [Arg $sprite : [AliasRef [Alias Sprite [RecordType [[RecordTypeField $pos [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (0)]][]]]]]] [[lookup [RecordTypeField $pos [AliasRef [Alias Position [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (0)]] [lookup [RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]]]
`)
}

func TestIf(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
isBestAge : (age: Int) -> Int =
    if age == 50 || age >= 65 then
        100
    else
        0
`, `
[ModuleDef $isBestAge = [FunctionValue ([[Arg $age : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [If [Logical [BoolOp [FunctionParamRef [Arg $age : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] EQ [Integer 50]] or [BoolOp [FunctionParamRef [Arg $age : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] GRE [Integer 65]]] then [Integer 100] else [Integer 0]]]]
`)
}

func TestIfFail(t *testing.T) {
	testDecorateWithoutDefaultFail(t,
		`
isBestAge : (age: Int) -> Bool =
    if age == 50 then
        100
    else
        0
`, &decorated.UnMatchingFunctionReturnTypesInFunctionValue{})
}

func TestIfPerson(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
isBeautiful : (name: String) -> Bool =
    if name == "Rebecca" then
        True
    else
        False
`, `
[ModuleDef $isBeautiful = [FunctionValue ([[Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [If [BoolOp [FunctionParamRef [Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] EQ [String Rebecca]] then [Bool true] else [Bool false]]]]
`)
}

func TestGuard(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
isBeautiful : (name: String) -> Bool =
    | name == "Rebecca" -> True
    | _ -> False
`, `
[ModuleDef $isBeautiful = [FunctionValue ([[Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [DGuard: [DGuardItem [BoolOp [FunctionParamRef [Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] EQ [String Rebecca]] [Bool true]] default: [Bool false]]]]
`)
}

func TestBoolPerson(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
isBeautiful : (name: String) -> Bool =
    name == "Rebecca"
`, `
[ModuleDef $isBeautiful = [FunctionValue ([[Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [BoolOp [FunctionParamRef [Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] EQ [String Rebecca]]]]
`)
}

func TestBoolPersonCall(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
isLoveOfMyLife : (name: String) -> Bool =
    name == "Rebecca"


main : (String) -> Bool =
    isLoveOfMyLife "Lisa"
`, `
[ModuleDef $isLoveOfMyLife = [FunctionValue ([[Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [BoolOp [FunctionParamRef [Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] EQ [String Rebecca]]]]
[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /isLoveOfMyLife]] [[String Lisa]]]]]
`)
}

func TestBoolPersonCallAgain(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
isLoveOfMyLife : (name: String, Int) -> Bool =
    name == "Rebecca"


main : (String) -> Bool =
    isLoveOfMyLife "Lisa" 2
`, `
[ModuleDef $isLoveOfMyLife = [FunctionValue ([[Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]] [Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [BoolOp [FunctionParamRef [Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] EQ [String Rebecca]]]]
[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /isLoveOfMyLife]] [[String Lisa] [Integer 2]]]]]
`)
}

func TestCustomType(t *testing.T) {
	testDecorateWithoutDefault(t,
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
	testDecorateWithoutDefault(t,
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
	testDecorateWithoutDefault(t,
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
	testDecorateWithoutDefault(t,
		`
f : (name: String, score: Int) -> Bool =
    if name == "Peter" then
        score * 2 > 100
    else
        score > 100


another : (score: Int) -> Bool =
    let
        af = f "Peter"
    in
    af score
`, `
[ModuleDef $f = [FunctionValue ([[Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]] [Arg $score : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [If [BoolOp [FunctionParamRef [Arg $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] EQ [String Peter]] then [BoolOp (Arithmetic [FunctionParamRef [Arg $score : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [Integer 2]) GR [Integer 100]] else [BoolOp [FunctionParamRef [Arg $score : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] GR [Integer 100]]]]]
[ModuleDef $another = [FunctionValue ([[Arg $score : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Let [[LetAssign [[LetVar $af]] = [Curry [FunctionRef [NamedDefinitionReference /f]] [[String Peter]]]]] in [FnCall [LetVarRef [LetVar $af]] [[FunctionParamRef [Arg $score : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]]]
`)
}

func TestBlob(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
a : (Blob) -> List Int =
    [ 10, 20, 99 ]
`, `
[ModuleDef $a = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Blob]]]]]) -> [ListLiteral [[Integer 10] [Integer 20] [Integer 99]]]]]
`)
}

func TestListLiteral2(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
a : (Bool) -> List Int =
    [ 10, 20, 99 ]
`, `
[ModuleDef $a = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [ListLiteral [[Integer 10] [Integer 20] [Integer 99]]]]]
`)
}

func TestTuple(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
a : (Bool) -> (Int, String) =
    (2, "Hello")
`, `
[ModuleDef $a = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [TupleLiteral [TupleType [[Primitive Int] [Primitive String]]] [TupleLiteral [#2 'Hello']] [[Integer 2] [String Hello]]]]]
`)
}

func xTestTupleSecond(t *testing.T) {
	testDecorate(t,
		`
a : Bool -> String
a x =
    Tuple.second (2, "Hello")
`, `
a = [functionvalue ([[arg $x = typeref $Bool [primitive Bool]]]) -> [fcall [functionref named definition reference [functionvalue ([[arg $tuple = [tupletype [[localtype a] [localtype b]]]]]) -> [asm callexternal 00 coreTupleSecond 01]]] [[TupleLiteral [tupletype [[primitive Int] [primitive String]]] [tuple-literal: [#2 'Hello']] [[integer 2] [str Hello]]]]]]
`)
}

func xTestTupleFirst(t *testing.T) {
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
	testDecorateWithoutDefault(t,
		`
createTuple : (first: a, second: b) -> (a, b) =
    (first, second)


a : (Bool) -> (Int, String) =
    createTuple 2 "Hello"
`, `
[ModuleDef $createTuple = [FunctionValue ([[Arg $first : [GenericParam a]] [Arg $second : [GenericParam b]]]) -> [TupleLiteral [TupleType [[GenericParam a] [GenericParam b]]] [TupleLiteral [$first $second]] [[FunctionParamRef [Arg $first : [GenericParam a]]] [FunctionParamRef [Arg $second : [GenericParam b]]]]]]]
[ModuleDef $a = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /createTuple]] [[Integer 2] [String Hello]]]]]
`)
}

func TestTupleGenericsFail(t *testing.T) {
	testDecorateWithoutDefaultFail(t,
		`
createTuple : (first: a, second: b) -> (a, b) =
    (first, second)


a : Bool -> (Int, String) =
    createTuple "2" "Hello"
`, &decorated.UnMatchingFunctionReturnTypesInFunctionValue{})
}

func TestArrayLiteral2(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
a : (Bool) -> Array Int =
    [| 10, 20, 99 |]
`, `
[ModuleDef $a = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [ArrayLiteral Array [[Integer 10] [Integer 20] [Integer 99]]]]]
`)
}

func TestListLiteral3(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
type alias Cool =
    { name : String
    }


a : (Bool) -> List Cool =
    [ { name = "hi" }, { name = "another" }, { name = "tjoho" } ]
`, `
Cool : [Alias Cool [RecordType [[RecordTypeField $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]][]]] => [RecordType [[RecordTypeField $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]][]]

[ModuleDef $a = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [ListLiteral [[RecordLiteral [RecordType [[RecordTypeField $name [Primitive String] (0)]][]] [0 = [String hi]]] [RecordLiteral [RecordType [[RecordTypeField $name [Primitive String] (0)]][]] [0 = [String another]]] [RecordLiteral [RecordType [[RecordTypeField $name [Primitive String] (0)]][]] [0 = [String tjoho]]]]]]]
`)
}

func TestListLiteral4Fail(t *testing.T) {
	testDecorateWithoutDefaultFail(t,
		`
type alias Cool =
    { name : Int
    }


a : (Bool) -> List Cool =
    [ { name = "hi" }, { name = "another" }, { name = "tjoho" } ]
`, &decorated.UnMatchingFunctionReturnTypesInFunctionValue{})
}

func TestRecordConstructor(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
type alias Cool =
    { name : Int
    }


a : (Bool) -> List Cool =
    [ Cool { name = 95 } ]
`, `
Cool : [Alias Cool [RecordType [[RecordTypeField $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]] => [RecordType [[RecordTypeField $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]

[ModuleDef $a = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [ListLiteral [[RecordConstructorRecord [CCall [TypeReference $Cool] [[RecordLiteral [[$name = #95]]]]] [RecordLiteral [RecordType [[RecordTypeField $name [Primitive Int] (0)]][]] [0 = [Integer 95]]]]]]]]
`)
}

func TestRecordConstructorWithoutSpace(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
type alias Cool =
    { name : Int
    }


a : (Bool) -> List Cool =
    [ Cool{ name = 95 } ]
`, `
Cool : [Alias Cool [RecordType [[RecordTypeField $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]] => [RecordType [[RecordTypeField $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]

[ModuleDef $a = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [ListLiteral [[RecordConstructorRecord [CCall [TypeReference $Cool] [[RecordLiteral [[$name = #95]]]]] [RecordLiteral [RecordType [[RecordTypeField $name [Primitive Int] (0)]][]] [0 = [Integer 95]]]]]]]]
`)
}

func TestRecordConstructorValues(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
type alias Cool =
    { name : Int
    }


a : (Bool) -> List Cool =
    [ Cool 2 ]
`, `
Cool : [Alias Cool [RecordType [[RecordTypeField $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]] => [RecordType [[RecordTypeField $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]

[ModuleDef $a = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [ListLiteral [[record-constructor [AliasRef [Alias Cool [RecordType [[RecordTypeField $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]]] [0 = [Integer 2]]]]]]]
`)
}

func TestRecordConstructorValuesFail(t *testing.T) {
	testDecorateWithoutDefaultFail(t,
		`
type alias Cool =
    { name : Int
    }


a : (Bool) -> List Cool =
    [ Cool "2" ]
`, &decorated.WrongTypeForRecordConstructorField{})
}

func TestCaseDefault(t *testing.T) {
	testDecorateWithoutDefault(t, // -- just a comment
		`
type CustomType =
    First
    | Second


some : (a: CustomType) -> String =
    case a of
        _ -> ""
`, `
CustomType : [CustomType CustomType [[Variant $First []] [Variant $Second []]]]
First : [Variant $First []]
Second : [Variant $Second []]

[ModuleDef $some = [FunctionValue ([[Arg $a : [VariantRef [NamedDefTypeRef :[TypeReference $CustomType]]]]]) -> [PMCase: [FunctionParamRef [Arg $a : [VariantRef [NamedDefTypeRef :[TypeReference $CustomType]]]]] of  default: [String ]]]]
`)
}

func TestCaseCustomTypeWithStructs(t *testing.T) {
	testDecorateWithoutDefault(t, // -- just a comment
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


some : (child: Child) -> String =
    case child of
        Aron x -> "Aron"

        Alexandra -> "Alexandris"

        _ -> ""
`, `
Tinkering : [Alias Tinkering [RecordType [[RecordTypeField $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (0)]][]]] => [RecordType [[RecordTypeField $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (0)]][]]
Studying : [Alias Studying [RecordType [[RecordTypeField $subject [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]][]]] => [RecordType [[RecordTypeField $subject [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]][]]
Work : [Alias Work [RecordType [[RecordTypeField $web [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]] => [RecordType [[RecordTypeField $web [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]
Child : [CustomType Child [[Variant $Aron [[AliasRef [Alias Tinkering [RecordType [[RecordTypeField $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (0)]][]]]]]] [Variant $Alexandra []] [Variant $Alma [[AliasRef [Alias Studying [RecordType [[RecordTypeField $subject [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]][]]]]]] [Variant $Isabelle [[AliasRef [Alias Work [RecordType [[RecordTypeField $web [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]]]]]]]
Aron : [Variant $Aron [[AliasRef [Alias Tinkering [RecordType [[RecordTypeField $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (0)]][]]]]]]
Alexandra : [Variant $Alexandra []]
Alma : [Variant $Alma [[AliasRef [Alias Studying [RecordType [[RecordTypeField $subject [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]][]]]]]]
Isabelle : [Variant $Isabelle [[AliasRef [Alias Work [RecordType [[RecordTypeField $web [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]]]]]

[ModuleDef $some = [FunctionValue ([[Arg $child : [VariantRef [NamedDefTypeRef :[TypeReference $Child]]]]]) -> [dcase: [FunctionParamRef [Arg $child : [VariantRef [NamedDefTypeRef :[TypeReference $Child]]]]] of [dcasecons [VariantRef [NamedDefTypeRef :[TypeReference $Aron]] [Variant $Aron [[AliasRef [Alias Tinkering [RecordType [[RecordTypeField $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (0)]][]]]]]]] ([[dcaseparm $x type:[AliasRef [Alias Tinkering [RecordType [[RecordTypeField $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (0)]][]]]]]]) => [String Aron]];[dcasecons [VariantRef [NamedDefTypeRef :[TypeReference $Alexandra]] [Variant $Alexandra []]] ([]) => [String Alexandris]] default: [String ]]]]
`)
}

func TestCaseStringAndDefault(t *testing.T) {
	testDecorateWithoutDefault(t, // -- just a comment
		`
some : (a: String) -> Int =
    case a of
        "hello" -> 0

        _ -> -1
`, `
[ModuleDef $some = [FunctionValue ([[Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [PMCase: [FunctionParamRef [Arg $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] of [PMCaseCons [String hello] => [Integer 0]] default: [Integer -1]]]]
`)
}

func TestRecordGenerics(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
type alias Tinkering t =
    { solder : Bool
    , secret : t
    }


f : (tinkering: Tinkering Int) -> Int =
    tinkering.secret
`, `
Tinkering : [Alias Tinkering [RecordType [[RecordTypeField $secret [GenericParam t] (0)] [RecordTypeField $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]][[GenericParam t]]]] => [RecordType [[RecordTypeField $secret [GenericParam t] (0)] [RecordTypeField $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]][[GenericParam t]]]

[ModuleDef $f = [FunctionValue ([[Arg $tinkering : [AliasRef [Alias Tinkering [RecordType [[RecordTypeField $secret [GenericParam t] (0)] [RecordTypeField $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]][[GenericParam t]]]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]>]]) -> [lookups [FunctionParamRef [Arg $tinkering : [AliasRef [Alias Tinkering [RecordType [[RecordTypeField $secret [GenericParam t] (0)] [RecordTypeField $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]][[GenericParam t]]]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]>]] [[lookup [RecordTypeField $secret [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]]]
`)
}

func TestGenericsFail(t *testing.T) {
	testDecorateWithoutDefaultFail(t,
		`
type alias Tinkering t =
    { solder : Bool
    , secret : t
    }


f : (tinkering: Tinkering String) -> Int =
    tinkering.secret
`, &decorated.UnMatchingFunctionReturnTypesInFunctionValue{})
}

func TestGenericsStructInstantiate(t *testing.T) {
	testDecorateWithoutDefault(t,
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
	testDecorateWithoutDefault(t, //-- just a comment
		`
type alias Enemy =
    { values : List Int
    }


type alias World =
    { enemies : List Enemy
    }


updateEnemy : (List Enemy) -> Bool =
    True


updateWorld : (w: World) -> Bool =
    updateEnemy w.enemies


main : (Bool) -> Bool =
    updateWorld { enemies = [ { values = [ 1, 3 ] } ] }
`, `
Enemy : [Alias Enemy [RecordType [[RecordTypeField $values [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]> (0)]][]]] => [RecordType [[RecordTypeField $values [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]> (0)]][]]
World : [Alias World [RecordType [[RecordTypeField $enemies [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[AliasRef [Alias Enemy [RecordType [[RecordTypeField $values [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]> (0)]][]]]]> (0)]][]]] => [RecordType [[RecordTypeField $enemies [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[AliasRef [Alias Enemy [RecordType [[RecordTypeField $values [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]> (0)]][]]]]> (0)]][]]

[ModuleDef $updateEnemy = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[AliasRef [Alias Enemy [RecordType [[RecordTypeField $values [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]> (0)]][]]]]>]]) -> [Bool true]]]
[ModuleDef $updateWorld = [FunctionValue ([[Arg $w : [AliasRef [Alias World [RecordType [[RecordTypeField $enemies [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[AliasRef [Alias Enemy [RecordType [[RecordTypeField $values [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]> (0)]][]]]]> (0)]][]]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /updateEnemy]] [[lookups [FunctionParamRef [Arg $w : [AliasRef [Alias World [RecordType [[RecordTypeField $enemies [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[AliasRef [Alias Enemy [RecordType [[RecordTypeField $values [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]> (0)]][]]]]> (0)]][]]]]]] [[lookup [RecordTypeField $enemies [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[AliasRef [Alias Enemy [RecordType [[RecordTypeField $values [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]> (0)]][]]]]> (0)]]]]]]]]
[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /updateWorld]] [[RecordLiteral [RecordType [[RecordTypeField $enemies [Primitive List<[RecordType [[RecordTypeField $values [Primitive List<[Primitive Int]>] (0)]][]]>] (0)]][]] [0 = [ListLiteral [[RecordLiteral [RecordType [[RecordTypeField $values [Primitive List<[Primitive Int]>] (0)]][]] [0 = [ListLiteral [[Integer 1] [Integer 3]]]]]]]]]]]]]
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


drawSprite : (Sprite) -> Bool =
    True


drawSprites : (sprites: List Sprite) -> List Bool =
    List.map drawSprite sprites


drawWorld : (world: World) -> List (List Bool) =
    List.map drawSprites world.drawTasks


main : (Bool) -> List (List Bool) =
    drawWorld { drawTasks = [ [ { x = 10, y = 20 }, { x = 44, y = 98 } ], [ { x = 99, y = 98 } ] ] }
`, `
Sprite : [Alias Sprite [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]] => [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]
World : [Alias World [RecordType [[RecordTypeField $drawTasks [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[AliasRef [Alias Sprite [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]>> (0)]][]]] => [RecordType [[RecordTypeField $drawTasks [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[AliasRef [Alias Sprite [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]>> (0)]][]]

[ModuleDef $drawSprite = [FunctionValue ([[Arg $_ : [AliasRef [Alias Sprite [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]]]) -> [Bool true]]]
[ModuleDef $drawSprites = [FunctionValue ([[Arg $sprites : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[AliasRef [Alias Sprite [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]>]]) -> [FnCall [FunctionRef [NamedDefinitionReference List/map]] [[FunctionRef [NamedDefinitionReference /drawSprite]] [FunctionParamRef [Arg $sprites : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[AliasRef [Alias Sprite [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]>]]]]]]
[ModuleDef $drawWorld = [FunctionValue ([[Arg $world : [AliasRef [Alias World [RecordType [[RecordTypeField $drawTasks [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[AliasRef [Alias Sprite [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]>> (0)]][]]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference List/map]] [[FunctionRef [NamedDefinitionReference /drawSprites]] [lookups [FunctionParamRef [Arg $world : [AliasRef [Alias World [RecordType [[RecordTypeField $drawTasks [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[AliasRef [Alias Sprite [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]>> (0)]][]]]]]] [[lookup [RecordTypeField $drawTasks [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[AliasRef [Alias Sprite [RecordType [[RecordTypeField $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]>> (0)]]]]]]]]
[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /drawWorld]] [[RecordLiteral [RecordType [[RecordTypeField $drawTasks [Primitive List<[Primitive List<[RecordType [[RecordTypeField $x [Primitive Int] (0)] [RecordTypeField $y [Primitive Int] (1)]][]]>]>] (0)]][]] [0 = [ListLiteral [[ListLiteral [[RecordLiteral [RecordType [[RecordTypeField $x [Primitive Int] (0)] [RecordTypeField $y [Primitive Int] (1)]][]] [0 = [Integer 10] 1 = [Integer 20]]] [RecordLiteral [RecordType [[RecordTypeField $x [Primitive Int] (0)] [RecordTypeField $y [Primitive Int] (1)]][]] [0 = [Integer 44] 1 = [Integer 98]]]]] [ListLiteral [[RecordLiteral [RecordType [[RecordTypeField $x [Primitive Int] (0)] [RecordTypeField $y [Primitive Int] (1)]][]] [0 = [Integer 99] 1 = [Integer 98]]]]]]]]]]]]]

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


intConvert : (MyList Int) -> Bool =
    True
`, `
MyList : [Alias MyList [RecordType [[RecordTypeField $someType [GenericParam t] (0)]][[GenericParam t]]]] => [RecordType [[RecordTypeField $someType [GenericParam t] (0)]][[GenericParam t]]]

[ModuleDef $intConvert = [FunctionValue ([[Arg $_ : [AliasRef [Alias MyList [RecordType [[RecordTypeField $someType [GenericParam t] (0)]][[GenericParam t]]]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]>]]) -> [Bool true]]]
`)
}

func TestBasicFunctionTypeParameters(t *testing.T) {
	testDecorateWithoutDefault(t,

		`
simple : (List a) -> a =
    2


main : (Bool) -> Int =
    simple [ 1, 2, 3 ]
`, `
[ModuleDef $simple = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[GenericParam a]>]]) -> [Integer 2]]]
[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /simple]] [[ListLiteral [[Integer 1] [Integer 2] [Integer 3]]]]]]]
`)
}

func TestFunctionTypeParameters(t *testing.T) {
	testDecorateWithoutDefault(t,

		`
intConvert : (List Int) -> Bool =
    True


simple : ((List a -> Bool), List a) -> Bool =
    True


main : (Bool) -> Bool =
    simple intConvert [ 1, 2, 3 ]
`, `
[ModuleDef $intConvert = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]>]]) -> [Bool true]]]
[ModuleDef $simple = [FunctionValue ([[Arg $_ : [FunctionTypeRef [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[GenericParam a]> [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]] [Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[GenericParam a]>]]) -> [Bool true]]]
[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FnCall [FunctionRef [NamedDefinitionReference /simple]] [[FunctionRef [NamedDefinitionReference /intConvert]] [ListLiteral [[Integer 1] [Integer 2] [Integer 3]]]]]]]
`)
}

func TestFunctionTypeGenerics(t *testing.T) {
	testDecorateWithoutDefault(t,

		`
intConvert : (Int) -> Int =
    42


simple : (a) -> (a -> a) =
    intConvert


main : (Bool) -> Int =
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
	testDecorateWithoutDefaultFail(t,

		`
intConvert : (Int) -> Int =
    42


simple : (a) -> (a -> a) =
    intConvert


main : (Bool) -> Int =
    let
        fn = simple "hello"
    in
    -- This should fail, since fn takes a String instead of an Int
    intConvert (fn 22)
`, &decorated.FunctionCallTypeMismatch{})
}

func TestUpdateComplex(t *testing.T) {
	testDecorateWithoutDefault(t,

		`
type alias Scale2 =
    { scaleX : Int
    , scaleY : Int
    }


type alias Sprite =
    { dummy : Int
    , scale : Scale2
    }


updateSprite : (inSprite: Sprite, newScale: Int) -> Sprite =
    { inSprite | scale = { scaleX = newScale, scaleY = newScale } }
`, `
Scale2 : [Alias Scale2 [RecordType [[RecordTypeField $scaleX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $scaleY [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]] => [RecordType [[RecordTypeField $scaleX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $scaleY [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]
Sprite : [Alias Sprite [RecordType [[RecordTypeField $dummy [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $scale [AliasRef [Alias Scale2 [RecordType [[RecordTypeField $scaleX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $scaleY [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (1)]][]]] => [RecordType [[RecordTypeField $dummy [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $scale [AliasRef [Alias Scale2 [RecordType [[RecordTypeField $scaleX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $scaleY [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (1)]][]]

[ModuleDef $updateSprite = [FunctionValue ([[Arg $inSprite : [AliasRef [Alias Sprite [RecordType [[RecordTypeField $dummy [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $scale [AliasRef [Alias Scale2 [RecordType [[RecordTypeField $scaleX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $scaleY [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (1)]][]]]]] [Arg $newScale : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [RecordLiteral [RecordType [[RecordTypeField $dummy [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $scale [AliasRef [Alias Scale2 [RecordType [[RecordTypeField $scaleX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $scaleY [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]] (1)]][]] [1 = [RecordLiteral [RecordType [[RecordTypeField $scaleX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $scaleY [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]] [0 = [FunctionParamRef [Arg $newScale : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] 1 = [FunctionParamRef [Arg $newScale : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]]]]
`)
}

func TestCaseAlignmentSurprisedByIndentation(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
checkMaybeInt : (a: Maybe Int, b: Maybe Int) -> Int =
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

func TestCaseNotCoveredByAllConsequencesFail(t *testing.T) {
	testDecorateWithoutDefaultFail(t, //    -- Must include default (_) or Nothing here
		`
checkMaybeInt : (a: Maybe Int) -> Int =
    case a of
        Just firstInt -> firstInt
`,
		&decorated.UnhandledCustomTypeVariants{})
}

func TestCaseCoveredMultipleTimesFail(t *testing.T) {
	testDecorateWithoutDefaultFail(t,
		`
checkMaybeInt : (a: Maybe Int) -> Int =
    case a of
        Just firstInt -> firstInt

        Nothing -> 0

        Nothing -> 0
`,
		&decorated.AlreadyHandledCustomTypeVariant{})
}

func TestFunctionAnnotationWithoutParameters(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
single : Int =
    2
`,
		`
[ModuleDef $single = [Constant [Integer 2]]]
    `)
}

func TestMinimalMaybeAndNothing(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
type alias Thing =
    { something : String
    }


main : (Bool) -> Maybe Thing =
    Nothing
`, `
Thing : [Alias Thing [RecordType [[RecordTypeField $something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]][]]] => [RecordType [[RecordTypeField $something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]][]]

[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [VariantConstructor [Variant $Nothing []] []]]]
`)
}

func TestConsequenceCheckFail(t *testing.T) {
	testDecorateWithoutDefaultFail(t, `
type Direction =
    NotMoving
    | Right
    | Down
    | Left
    | Up


-- returns true if direction is vertical
isVertical : (direction: Direction) -> Bool =
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

func TestConsequenceCheck3Fail(t *testing.T) {
	testDecorateWithoutDefaultFail(t, `
type Direction =
    NotMoving
    | Right
    | Down
    | Left
    | Up


-- returns true if direction is vertical
isVertical : (direction: Direction) -> Bool =
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

func TestConsequenceCheck2Fail(t *testing.T) {
	testDecorateWithoutDefaultFail(t, `
type Direction =
    NotMoving
    | Right
    | Down
    | Left
    | Up


isVertical : (direction: Direction) -> Int =
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
	testDecorateWithoutDefault(t,
		`
type alias Thing =
    { something : String
    }


main : (Bool) -> Maybe Thing =
    Just { something = "hello" }
`, `
Thing : [Alias Thing [RecordType [[RecordTypeField $something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]][]]] => [RecordType [[RecordTypeField $something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]][]]

[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [VariantConstructor [Variant $Just [[GenericParam a]]] [[RecordLiteral [RecordType [[RecordTypeField $something [Primitive String] (0)]][]] [0 = [String hello]]]]]]]
`)
}

func TestJustAndNothing(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
type alias Thing =
    { something : String
    }


main : (Bool) -> Maybe Thing =
    if True then
        Nothing
    else
        Just { something = "hello" }
`, `
Thing : [Alias Thing [RecordType [[RecordTypeField $something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]][]]] => [RecordType [[RecordTypeField $something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]][]]

[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [If [Bool true] then [VariantConstructor [Variant $Nothing []] []] else [VariantConstructor [Variant $Just [[GenericParam a]]] [[RecordLiteral [RecordType [[RecordTypeField $something [Primitive String] (0)]][]] [0 = [String hello]]]]]]]]
`)
}

func TestOwnGenericGet(t *testing.T) {
	testDecorateWithoutDefault(t,
		`

type alias Thing a =
    { something : a
    }


get : (Thing a) -> Maybe a =
    Nothing


main : (Bool) -> Maybe (Thing String) =
    let
        v = { something = "Something" }
        _ = get v
    in
    Just v
        
`, `
Thing : [Alias Thing [RecordType [[RecordTypeField $something [GenericParam a] (0)]][[GenericParam a]]]] => [RecordType [[RecordTypeField $something [GenericParam a] (0)]][[GenericParam a]]]

[ModuleDef $get = [FunctionValue ([[Arg $_ : [AliasRef [Alias Thing [RecordType [[RecordTypeField $something [GenericParam a] (0)]][[GenericParam a]]]]]<[GenericParam a]>]]) -> [VariantConstructor [Variant $Nothing []] []]]]
[ModuleDef $main = [FunctionValue ([[Arg $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [Let [[LetAssign [[LetVar $v]] = [RecordLiteral [RecordType [[RecordTypeField $something [Primitive String] (0)]][]] [0 = [String Something]]]] [LetAssign [[LetVar $_]] = [FnCall [FunctionRef [NamedDefinitionReference /get]] [[LetVarRef [LetVar $v]]]]]] in [VariantConstructor [Variant $Nothing []] []]]]]
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


checkMaybeGamepad : (Maybe Gamepad) -> PlayerAction =
    None


main : (oldUserInputs: UserInput) -> PlayerInputs =
    let
        gamepads = oldUserInputs.gamepads

        maybeOld = Array.get 0 gamepads

        playerAction = checkMaybeGamepad maybeOld
    in
    { inputs = [ playerAction ] }
`,
		`
PlayerAction : [CustomType PlayerAction [[Variant $None []] [Variant $Jump []]]]
None : [Variant $None []]
Jump : [Variant $Jump []]
Gamepad : [Alias Gamepad [RecordType [[RecordTypeField $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]] => [RecordType [[RecordTypeField $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]
UserInput : [Alias UserInput [RecordType [[RecordTypeField $gamepads [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Array]]]<[AliasRef [Alias Gamepad [RecordType [[RecordTypeField $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]]]> (0)]][]]] => [RecordType [[RecordTypeField $gamepads [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Array]]]<[AliasRef [Alias Gamepad [RecordType [[RecordTypeField $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]]]> (0)]][]]
PlayerInputs : [Alias PlayerInputs [RecordType [[RecordTypeField $inputs [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[VariantRef [NamedDefTypeRef :[TypeReference $PlayerAction]]]> (0)]][]]] => [RecordType [[RecordTypeField $inputs [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]]<[VariantRef [NamedDefTypeRef :[TypeReference $PlayerAction]]]> (0)]][]]

[ModuleDef $checkMaybeGamepad = [FunctionValue ([[Arg $_ : [VariantRef [NamedDefTypeRef :[TypeReference $Maybe]]]<[AliasRef [Alias Gamepad [RecordType [[RecordTypeField $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]]]>]]) -> [VariantConstructor [Variant $None []] []]]]
[ModuleDef $main = [FunctionValue ([[Arg $oldUserInputs : [AliasRef [Alias UserInput [RecordType [[RecordTypeField $gamepads [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Array]]]<[AliasRef [Alias Gamepad [RecordType [[RecordTypeField $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]]]> (0)]][]]]]]]) -> [Let [[LetAssign [[LetVar $gamepads]] = [lookups [FunctionParamRef [Arg $oldUserInputs : [AliasRef [Alias UserInput [RecordType [[RecordTypeField $gamepads [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Array]]]<[AliasRef [Alias Gamepad [RecordType [[RecordTypeField $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]]]> (0)]][]]]]]] [[lookup [RecordTypeField $gamepads [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Array]]]<[AliasRef [Alias Gamepad [RecordType [[RecordTypeField $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]][]]]]> (0)]]]]] [LetAssign [[LetVar $maybeOld]] = [FnCall [FunctionRef [NamedDefinitionReference Array/get]] [[Integer 0] [LetVarRef [LetVar $gamepads]]]]] [LetAssign [[LetVar $playerAction]] = [FnCall [FunctionRef [NamedDefinitionReference /checkMaybeGamepad]] [[LetVarRef [LetVar $maybeOld]]]]]] in [RecordLiteral [RecordType [[RecordTypeField $inputs [Primitive List<[VariantRef [NamedDefTypeRef :[TypeReference $PlayerAction]]]>] (0)]][]] [0 = [ListLiteral [[LetVarRef [LetVar $playerAction]]]]]]]]]
`)
}

func TestArrayIntGetWithMaybe(t *testing.T) {
	testDecorate(t,
		`

getInt : (Maybe Int) -> Int =
    22


main : (ints: Array Int) -> Int =
    let
        maybeInt = Array.get 0 ints

        extractedInt = getInt maybeInt
    in
    extractedInt
`,
		`
[ModuleDef $getInt = [FunctionValue ([[Arg $_ : [VariantRef [NamedDefTypeRef :[TypeReference $Maybe]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]>]]) -> [Integer 22]]]
[ModuleDef $main = [FunctionValue ([[Arg $ints : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Array]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]>]]) -> [Let [[LetAssign [[LetVar $maybeInt]] = [FnCall [FunctionRef [NamedDefinitionReference Array/get]] [[Integer 0] [FunctionParamRef [Arg $ints : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Array]]]<[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]>]]]]] [LetAssign [[LetVar $extractedInt]] = [FnCall [FunctionRef [NamedDefinitionReference /getInt]] [[LetVarRef [LetVar $maybeInt]]]]]] in [LetVarRef [LetVar $extractedInt]]]]]
`)
}

func TestArrayIntGetWithMaybeFail(t *testing.T) {
	testDecorateFail(t,
		`

getInt : (Maybe String) -> Int =
    22


main : (ints: Array Int) -> Int =
    let
        maybeInt = Array.get 0 ints

        extractedInt = getInt maybeInt
    in
    extractedInt
`,
		&decorated.CouldNotSmashFunctions{})
}

func TestArrayFromNothing(t *testing.T) {
	testDecorate(t,
		`
needInputs : (Array Int) -> Int =
    2


a : (Bool) -> Int =
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
needInputs : (Array Int) -> Int
    2


returnSomeFunc : (Bool) -> (Array Int -> Int) =
    needInputs


a : (Bool) -> Int =
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
	testDecorateWithoutDefault(t,
		`
type alias Something =
    { time : Int
    , playerX : Int
    }


fn : (something: Something) -> Something =
    { something | playerX = 3 }
`,
		`
Something : [Alias Something [RecordType [[RecordTypeField $playerX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $time [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]] => [RecordType [[RecordTypeField $playerX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $time [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]

[ModuleDef $fn = [FunctionValue ([[Arg $something : [AliasRef [Alias Something [RecordType [[RecordTypeField $playerX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $time [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]]]]]]) -> [RecordLiteral [RecordType [[RecordTypeField $playerX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [RecordTypeField $time [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]][]] [0 = [Integer 3]]]]]
`)
}

func TestPipeRight2(t *testing.T) {
	testDecorateWithoutDefault(t, `
first : (Int) -> Int =
    2


second : (String, Int) -> Int =
    2


third : (a: Int) -> Bool =
    a > 45


tester : (b: String) -> Bool =
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
