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
[ModuleDef $first = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
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
[ModuleDef $first = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $startNumber : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [ParamRef [Param $startNumber : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [ParamRef [Param $startNumber : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
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
[ModuleDef $another = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] ([]) -> [ConstantRef [NamedDefinitionReference /fn]]]]
`)
}

func TestBooleanLookup(t *testing.T) {
	testDecorateWithoutDefault(t, `
another : Bool =
    let
        a = true
        b = false
    in
    a && b
`, `
[ModuleDef $another = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([]) -> [Let [[LetAssign [[LetVar $a]] = [Bool true]] [LetAssign [[LetVar $b]] = [Bool false]]] in [Logical [LetVarRef $a] and [LetVarRef $b]]]]]
`)
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
Something : [Alias Something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]

[ModuleDef $another = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([]) -> [Let [[LetAssign [[LetVar $b]] = [Cast [Integer 32] [AliasRefExpr [AliasRef [Alias Something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]] in [BoolOp [LetVarRef $b] GRE [Integer 32]]]]]
`)
}

func TestFunctionTypeWithoutParameters(t *testing.T) {
	testDecorateWithoutDefault(t, `__externalvarfn another : (Int, String) -> Int

`, `
[ModuleDef $another = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]] [Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [ExternalFunctionDeclarationExpression [FnDeclExpr 2] [Primitive Any]]]]`)
}

func TestTypeId(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
--| ignore this
type alias Something = Int


main : TypeId Something =
    $Something
`,
		`
Something : [Alias Something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]

[ModuleDef $main = [Constant [TypeIdLiteral [ResolvedContext a:[AliasRef [Alias Something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $TypeRef]] [LocalTypeNameOnlyContext a => [Primitive TypeRef<[LocalTypeNameRef a]>]]]]]]]
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
[ModuleDef $first = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $ResourceName]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $ResourceName]]]]]) -> [Integer 2]]]
[ModuleDef $main = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $ResourceName]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [FunctionRef [NamedDefinitionReference /first]] [[ResourceName this/is/something.txt]]]]]
`)
}

func TestAnyMatchingType(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
fn : (*) -> Int =
    23
`,
		`
[ModuleDef $fn = [FunctionValue [FunctionType [[AnyMatching AnyMatchingType] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [AnyMatching AnyMatchingType]]]) -> [Integer 23]]]
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
[ModuleDef $fn = [FunctionValue [FunctionType [[AnyMatching AnyMatchingType] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [AnyMatching AnyMatchingType]]]) -> [Integer 23]]]
[ModuleDef $main = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([]) -> [FnCall [FunctionType [[AnyMatching AnyMatchingType] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [FunctionRef [NamedDefinitionReference /fn]] [[Fixed 23000] [String hello] [Integer 42]]]]]
`)
}

func TestAnyMatchingTypeCallMiddle(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
fn : (String, *) -> Int =
    23


main : Int =
    fn "hello" 42.0 33
`,
		`
[ModuleDef $fn = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [AnyMatching AnyMatchingType] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]] [Param $_ : [AnyMatching AnyMatchingType]]]) -> [Integer 23]]]
[ModuleDef $main = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([]) -> [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [AnyMatching AnyMatchingType] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [FunctionRef [NamedDefinitionReference /fn]] [[String hello] [Fixed 42000] [Integer 33]]]]]
`)
}

func xTestAnyMatchingTypeCallMiddleLocalType(t *testing.T) {
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

func xTestAnyMatchingTypeCallMiddleLocalTypeFn(t *testing.T) {
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
[ModuleDef $first = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]]]] ([[Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]]]]) -> (Arithmetic [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]]]] FMULTIPLY [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]]]])]]
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
[ModuleDef $first = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Char]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Char]]]]]) -> [Integer 2]]]
[ModuleDef $main = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Char]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [FunctionRef [NamedDefinitionReference /first]] [[Char 99]]]]]
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
[ModuleDef $first = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]]]] ([[Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]]]] [FunctionRef [NamedDefinitionReference Int/toFixed]] [[ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]]
[ModuleDef $another = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]]]] [FunctionRef [NamedDefinitionReference /first]] [[Integer 2]]]]]
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
[ModuleDef $first = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]]]]) -> [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [FunctionRef [NamedDefinitionReference Int/round]] [[ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]]]]]]]]
[ModuleDef $another = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [FunctionRef [NamedDefinitionReference /first]] [[Fixed 2300]]]]]
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
[ModuleDef $first = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Fixed 300]]]
[ModuleDef $another = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Let [[LetAssign [[LetVar $x]] = [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Fixed]]]]] [FunctionRef [NamedDefinitionReference /first]] [[Integer 2]]]]] in [LetVarRef $x]]]]
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
State : [Alias State [RecordType [[Field $playerX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]

[ModuleDef $another = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [AliasRef [Alias State [RecordType [[Field $playerX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Let [[LetAssign [[LetVar $state]] = [RecordLiteral [RecordType [[Field $playerX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]] [0 = [Integer 0]]]]] in [RecordLiteral [RecordType [[Field $playerX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]] [0 = [Integer 22]]]]]]
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
[ModuleDef $first = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
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
[ModuleDef $first = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
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
[ModuleDef $first = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
[ModuleDef $second = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [ParamRef [Param $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] PLUS [ParamRef [Param $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
[ModuleDef $third = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $c : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [ParamRef [Param $c : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MINUS [ParamRef [Param $c : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
`)
}

func xTestWrongSpacingMultiDocFail(t *testing.T) {
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
		&decorated.FunctionArgumentTypeMismatch{})
}

func TestSomething2(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
first : (a: Int) -> Int =
    a * a
`,
		`
[ModuleDef $first = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
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
Status : [CustomType Status [[Variant $Unknown] [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]
Unknown : [Variant $Unknown]
Something : [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]

[ModuleDef $receiveStatus = [FunctionValue [FunctionType [[CustomTypeRef [NamedDefTypeRef :[TypeReference $Status]] [CustomType Status [[Variant $Unknown] [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]] [CustomTypeRef [NamedDefTypeRef :[TypeReference $Status]] [CustomType Status [[Variant $Unknown] [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]]] ([[Param $status : [CustomTypeRef [NamedDefTypeRef :[TypeReference $Status]] [CustomType Status [[Variant $Unknown] [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]]]) -> [ParamRef [Param $status : [CustomTypeRef [NamedDefTypeRef :[TypeReference $Status]] [CustomType Status [[Variant $Unknown] [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]]]]]
[ModuleDef $someFunc = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [CustomTypeRef [NamedDefTypeRef :[TypeReference $Status]] [CustomType Status [[Variant $Unknown] [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [FnCall [FunctionType [[CustomTypeRef [NamedDefTypeRef :[TypeReference $Status]] [CustomType Status [[Variant $Unknown] [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]] [CustomTypeRef [NamedDefTypeRef :[TypeReference $Status]] [CustomType Status [[Variant $Unknown] [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]]] [FunctionRef [NamedDefinitionReference /receiveStatus]] [[VariantConstructor [Variant $Unknown] []]]]]]
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
Status : [CustomType Status [[Variant $Unknown] [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]
Unknown : [Variant $Unknown]
Something : [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]
Unrelated : [CustomType Unrelated [[Variant $ShouldNotMatch] [Variant $SomethingElse [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]
ShouldNotMatch : [Variant $ShouldNotMatch]
SomethingElse : [Variant $SomethingElse [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]

[ModuleDef $receiveStatus = [FunctionValue [FunctionType [[VariantRef [NamedDefTypeRef :[TypeReference $Unknown]] [Variant $Unknown]] [CustomTypeRef [NamedDefTypeRef :[TypeReference $Status]] [CustomType Status [[Variant $Unknown] [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]]] ([[Param $status : [VariantRef [NamedDefTypeRef :[TypeReference $Unknown]] [Variant $Unknown]]]]) -> [ParamRef [Param $status : [VariantRef [NamedDefTypeRef :[TypeReference $Unknown]] [Variant $Unknown]]]]]]
[ModuleDef $someFunc = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [CustomTypeRef [NamedDefTypeRef :[TypeReference $Status]] [CustomType Status [[Variant $Unknown] [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [FnCall [FunctionType [[VariantRef [NamedDefTypeRef :[TypeReference $Unknown]] [Variant $Unknown]] [CustomTypeRef [NamedDefTypeRef :[TypeReference $Status]] [CustomType Status [[Variant $Unknown] [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]]] [FunctionRef [NamedDefinitionReference /receiveStatus]] [[VariantConstructor [Variant $Unknown] []]]]]]
`)
}

func TestCustomTypeVariantEqual(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
someFunc : (name: String) -> Bool =
    name == "Something"
`,
		`
[ModuleDef $someFunc = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [BoolOp [ParamRef [Param $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] EQ [String Something]]]]
`)
}

func TestUnknownAnnotationTypeFail(t *testing.T) {
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
Status : [CustomType Status [[Variant $Unknown] [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]
Unknown : [Variant $Unknown]
Something : [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]

[ModuleDef $someFunc = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [CustomTypeRef [NamedDefTypeRef :[TypeReference $Status]] [CustomType Status [[Variant $Unknown] [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [VariantConstructor [Variant $Something [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [[Integer 42]]]]]
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
Lister : [Alias Lister [LocalTypeNameOnlyContext t, u => [RecordType [[Field $another [LocalTypeNameRef u] (0)] [Field $fake [LocalTypeNameRef t] (1)]]]]]
Sprite : [Alias Sprite [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]

[ModuleDef $someFunc = [FunctionValue [FunctionType [[ResolvedContext t:[AliasRef [Alias Sprite [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]], u:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference <nil> [LocalTypeNameOnlyContext t, u => [RecordType [[Field $another [LocalTypeNameRef u] (0)] [Field $fake [LocalTypeNameRef t] (1)]]]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $a : [ResolvedContext t:[AliasRef [Alias Sprite [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]], u:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference <nil> [LocalTypeNameOnlyContext t, u => [RecordType [[Field $another [LocalTypeNameRef u] (0)] [Field $fake [LocalTypeNameRef t] (1)]]]]]]]]) -> [lookups [ParamRef [Param $a : [ResolvedContext t:[AliasRef [Alias Sprite [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]], u:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference <nil> [LocalTypeNameOnlyContext t, u => [RecordType [[Field $another [LocalTypeNameRef u] (0)] [Field $fake [LocalTypeNameRef t] (1)]]]]]]]] [[lookup [Field $another [ConcreteGenericRef [GenericParam u] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]] (0)]]]]]]
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
[ModuleDef $first = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
[ModuleDef $second = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]] [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [BoolOp [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] GR [Integer 25]]]]
[ModuleDef $third = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]]]
[ModuleDef $tester = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [FunctionRef [NamedDefinitionReference /first]] [(Arithmetic [Integer 2] PLUS [Integer 2])]] |> [Curry [FunctionRef [NamedDefinitionReference /second]] [[ParamRef [Param $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]]] |> [FunctionRef [NamedDefinitionReference /third]]]]
 `)
}

func TestBasicStringInterpolation(t *testing.T) {
	testDecorate(t,
		`
sample : (a: Int) -> String =
    $"hello {a}"
`, `
[ModuleDef $sample = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] ([[Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [StringInterpolation '[str:hello {a}]']]]
`)
}

func xTestArrayVsListFail(t *testing.T) {
	testDecorateFail(t,
		`
updater : (Int, String) -> Int =
    42


sample : (Int) -> List Int =
    let
        intArray = Array.fromList [ 0, 1, 2 ]

        arraySlice = Array.slice 0 2 intArray
    in
    List.map2 updater arraySlice [ "hello", "world", "fail" ]
`, &decorated.FunctionArgumentTypeMismatch{})
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
Position : [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]
Sprite : [Alias Sprite [RecordType [[Field $rootPosition [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] (0)]]]]

[ModuleDef $move = [FunctionValue [FunctionType [[AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]]]] ([[Param $pos : [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]]] [Param $delta : [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]]]]) -> [Let [[LetAssign [[LetVar $newX]] = (Arithmetic [lookups [ParamRef [Param $pos : [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]]]] [[lookup [Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]] PLUS [lookups [ParamRef [Param $delta : [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]]]] [[lookup [Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]])] [LetAssign [[LetVar $newY]] = (Arithmetic [lookups [ParamRef [Param $pos : [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]]]] [[lookup [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]] PLUS [lookups [ParamRef [Param $delta : [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]]]] [[lookup [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]])]] in [RecordLiteral [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]] [0 = [LetVarRef $newX] 1 = [LetVarRef $newY]]]]]]
[ModuleDef $moveSprite = [FunctionValue [FunctionType [[AliasRef [Alias Sprite [RecordType [[Field $rootPosition [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] (0)]]]]] [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] [AliasRef [Alias Sprite [RecordType [[Field $rootPosition [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] (0)]]]]]]] ([[Param $sprite : [AliasRef [Alias Sprite [RecordType [[Field $rootPosition [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] (0)]]]]]] [Param $delta : [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]]]]) -> [RecordLiteral [RecordType [[Field $rootPosition [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] (0)]]] [0 = [FnCall [FunctionType [[AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]]]] [FunctionRef [NamedDefinitionReference /move]] [[lookups [ParamRef [Param $sprite : [AliasRef [Alias Sprite [RecordType [[Field $rootPosition [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] (0)]]]]]]] [[lookup [Field $rootPosition [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] (0)]]]] [ParamRef [Param $delta : [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]]]]]]]]]]
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
[ModuleDef $first = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
[ModuleDef $second = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [BoolOp [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] GR [Integer 25]]]]
[ModuleDef $tester = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [FunctionRef [NamedDefinitionReference /first]] [(Arithmetic [Integer 2] PLUS [Integer 2])]] |> [FunctionRef [NamedDefinitionReference /second]]]]
  `)
}

func xTestOperatorPipeLeft(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
initial : (a: Int) -> String =
    "hej"

first : (a: Int, String) -> Int =
    a * a


second : (a: Int) -> Bool =
    a > 25


tester : (String) -> Bool =
    second <| first (2 + 2) <| initial 2
`, `
[ModuleDef $first = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
[ModuleDef $second = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [BoolOp [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] GR [Integer 25]]]]
[ModuleDef $tester = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [FunctionRef [NamedDefinitionReference /second]] <| [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [FunctionRef [NamedDefinitionReference /first]] [(Arithmetic [Integer 2] PLUS [Integer 2])]]]]
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
[ModuleDef $first = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
[ModuleDef $second = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [BoolOp [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] GR [Integer 25]]]]
[ModuleDef $tester = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [FunctionRef [NamedDefinitionReference /first]] [(Arithmetic [Integer 2] PLUS [Integer 2])]] |> [FunctionRef [NamedDefinitionReference /second]]]]
`)
}

func TestBasicCall(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
move : (pos: Int, delta: Int) -> Int =
    pos + delta
`, `
[ModuleDef $move = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $pos : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]] [Param $delta : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [ParamRef [Param $pos : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] PLUS [ParamRef [Param $delta : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
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
    { a = 2, b = true }
`, `
Constructor : [Alias Constructor [RecordType [[Field $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $b [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]]]]

[ModuleDef $create = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [AliasRef [Alias Constructor [RecordType [[Field $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $b [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [RecordLiteral [RecordType [[Field $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $b [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]]] [0 = [Integer 2] 1 = [Bool true]]]]]
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
[ModuleDef $create = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [ConstantRef [NamedDefinitionReference /tileHeight]]]]
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
SomeCustomType : [CustomType SomeCustomType [[Variant $First [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] [Variant $Anon] [Variant $Second [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]
First : [Variant $First [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]
Anon : [Variant $Anon]
Second : [Variant $Second [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]

[ModuleDef $a = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [CustomTypeRef [NamedDefTypeRef :[TypeReference $SomeCustomType]] [CustomType SomeCustomType [[Variant $First [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] [Variant $Anon] [Variant $Second [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [VariantConstructor [Variant $First [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] [[String Hello]]]]]
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
[ModuleDef $move = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $pos : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]] [Param $delta : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Let [[LetAssign [[LetVar $tenMore]] = (Arithmetic [ParamRef [Param $pos : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] PLUS [Integer 10])] [LetAssign [[LetVar $tenLess]] = (Arithmetic [ParamRef [Param $delta : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MINUS [Integer 10])]] in (Arithmetic [LetVarRef $tenMore] PLUS [LetVarRef $tenLess])]]]
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
[ModuleDef $move = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $pos : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]] [Param $delta : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [ParamRef [Param $pos : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] PLUS [ParamRef [Param $delta : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]])]]
[ModuleDef $main = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [Let [[LetAssign [[LetVar $fn]] = [FunctionRef [NamedDefinitionReference /move]]]] in [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [LetVarRef $fn] [[Integer 2] [Integer 3]]]]]]
`)
}

func TestBasicAppend(t *testing.T) {
	testDecorate(t,
		`
a : (Int) -> List Int =
    [ 1, 3, 4 ] ++ [ 5, 6, 7, 8 ]
`, `
[ModuleDef $a = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> (Arithmetic [ListLiteral [[Integer 1] [Integer 3] [Integer 4]]] APPEND [ListLiteral [[Integer 5] [Integer 6] [Integer 7] [Integer 8]]])]]
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
Position : [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]

[ModuleDef $getx = [FunctionValue [FunctionType [[AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $pos : [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]]]]) -> (Arithmetic [lookups [ParamRef [Param $pos : [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]]]] [[lookup [Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]] PLUS [lookups [ParamRef [Param $pos : [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]]]] [[lookup [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]])]]
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
Position : [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]
Sprite : [Alias Sprite [RecordType [[Field $pos [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] (0)]]]]

[ModuleDef $getx = [FunctionValue [FunctionType [[AliasRef [Alias Sprite [RecordType [[Field $pos [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] (0)]]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $sprite : [AliasRef [Alias Sprite [RecordType [[Field $pos [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] (0)]]]]]]]) -> [lookups [ParamRef [Param $sprite : [AliasRef [Alias Sprite [RecordType [[Field $pos [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] (0)]]]]]]] [[lookup [Field $pos [AliasRef [Alias Position [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] (0)]] [lookup [Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]]]
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
[ModuleDef $isBestAge = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $age : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [If [Logical [BoolOp [ParamRef [Param $age : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] EQ [Integer 50]] or [BoolOp [ParamRef [Param $age : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] GRE [Integer 65]]] then [Integer 100] else [Integer 0]]]]
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
        true
    else
        false
`, `
[ModuleDef $isBeautiful = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [If [BoolOp [ParamRef [Param $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] EQ [String Rebecca]] then [Bool true] else [Bool false]]]]
`)
}

func TestGuard(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
isBeautiful : (name: String) -> Bool =
    | name == "Rebecca" -> true
    | _ -> false
`, `
[ModuleDef $isBeautiful = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [DGuard: [DGuardItem [BoolOp [ParamRef [Param $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] EQ [String Rebecca]] [Bool true]] default: [Bool false]]]]
`)
}

func TestBoolPerson(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
isBeautiful : (name: String) -> Bool =
    name == "Rebecca"
`, `
[ModuleDef $isBeautiful = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [BoolOp [ParamRef [Param $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] EQ [String Rebecca]]]]
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
[ModuleDef $isLoveOfMyLife = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [BoolOp [ParamRef [Param $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] EQ [String Rebecca]]]]
[ModuleDef $main = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] [FunctionRef [NamedDefinitionReference /isLoveOfMyLife]] [[String Lisa]]]]]
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
[ModuleDef $isLoveOfMyLife = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]] [Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [BoolOp [ParamRef [Param $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] EQ [String Rebecca]]]]
[ModuleDef $main = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] [FunctionRef [NamedDefinitionReference /isLoveOfMyLife]] [[String Lisa] [Integer 2]]]]]
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
Chore : [CustomType Chore [[Variant $Meeting [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] [Variant $Running [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [Variant $Unknown]]]
Meeting : [Variant $Meeting [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]
Running : [Variant $Running [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]
Unknown : [Variant $Unknown]
`)
}

func TestCustomTypeGenerics(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
type Perhaps a =
    None
    | Actual a
`, `
Perhaps : [LocalTypeNameOnlyContext a => [CustomType Perhaps [[Variant $None] [Variant $Actual [[LocalTypeNameRef a]]]]]]
None : [Variant $None]
Actual : [Variant $Actual [[LocalTypeNameRef a]]]
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
Tinkering : [Alias Tinkering [RecordType [[Field $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (0)]]]]
Studying : [Alias Studying [RecordType [[Field $subject [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]]]
Work : [Alias Work [RecordType [[Field $web [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]
Child : [CustomType Child [[Variant $Aron [[AliasRef [Alias Tinkering [RecordType [[Field $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (0)]]]]]]] [Variant $Alexandra] [Variant $Alma [[AliasRef [Alias Studying [RecordType [[Field $subject [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]]]]]] [Variant $Isabelle [[AliasRef [Alias Work [RecordType [[Field $web [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]]]]]]
Aron : [Variant $Aron [[AliasRef [Alias Tinkering [RecordType [[Field $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (0)]]]]]]]
Alexandra : [Variant $Alexandra]
Alma : [Variant $Alma [[AliasRef [Alias Studying [RecordType [[Field $subject [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]]]]]]
Isabelle : [Variant $Isabelle [[AliasRef [Alias Work [RecordType [[Field $web [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]]]]
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
[ModuleDef $f = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]] [Param $score : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [If [BoolOp [ParamRef [Param $name : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] EQ [String Peter]] then [BoolOp (Arithmetic [ParamRef [Param $score : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] MULTIPLY [Integer 2]) GR [Integer 100]] else [BoolOp [ParamRef [Param $score : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] GR [Integer 100]]]]]
[ModuleDef $another = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $score : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Let [[LetAssign [[LetVar $af]] = [Curry [FunctionRef [NamedDefinitionReference /f]] [[String Peter]]]]] in [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] [LetVarRef $af] [[ParamRef [Param $score : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]]]
`)
}

func TestBlob(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
a : (Blob) -> List Int =
    [ 10, 20, 99 ]
`, `
[ModuleDef $a = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Blob]]] [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Blob]]]]]) -> [ListLiteral [[Integer 10] [Integer 20] [Integer 99]]]]]
`)
}

func TestListLiteral2(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
a : (Bool) -> List Int =
    [ 10, 20, 99 ]
`, `
[ModuleDef $a = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [ListLiteral [[Integer 10] [Integer 20] [Integer 99]]]]]
`)
}

func TestTuple(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
a : (Bool) -> (Int, String) =
    (2, "Hello")
`, `
[ModuleDef $a = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [TupleType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [TupleLiteral [TupleType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] [TupleLiteral [#2 'Hello']] [[Integer 2] [String Hello]]]]]
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
[ModuleDef $createTuple = [FunctionValue [LocalTypeNameOnlyContext a, b => [FunctionType [[LocalTypeNameRef a] [LocalTypeNameRef b] [TupleType [[LocalTypeNameRef a] [LocalTypeNameRef b]]]]]] ([[Param $first : [LocalTypeNameRef a]] [Param $second : [LocalTypeNameRef b]]]) -> [TupleLiteral [TupleType [[LocalTypeNameRef a] [LocalTypeNameRef b]]] [TupleLiteral [$first $second]] [[ParamRef [Param $first : [LocalTypeNameRef a]]] [ParamRef [Param $second : [LocalTypeNameRef b]]]]]]]
[ModuleDef $a = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [TupleType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FnCall [FunctionType [[Primitive Int] [Primitive String] [TupleType [[ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]] [ConcreteGenericRef [GenericParam b] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]]]] [FunctionRef [NamedDefinitionReference /createTuple]] [[Integer 2] [String Hello]]]]]
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

/*
func TestTypeRef(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
haveTypeRefAsParameter : (TypeRef a) -> Bool

main : Bool =
    haveTypeRefAsParameter $String
`, `
[ModuleDef $haveTypeRefAsParameter = [FunctionValue [LocalTypeNameOnlyContext a => [FunctionType [[LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $TypeRef]] [LocalTypeNameOnlyContext a => [Primitive TypeRef<[LocalTypeNameRef a]>]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]] ([[Param $_ : [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $TypeRef]] [LocalTypeNameOnlyContext a => [Primitive TypeRef<[LocalTypeNameRef a]>]]]]]) -> [ExternalFunctionDeclarationExpression [FnDeclExpr 0] [Primitive Any]]]]
[ModuleDef $main = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([]) -> [FnCall [FunctionType [[Primitive TypeRef<[ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]>] [Primitive Bool]]] [FunctionRef [NamedDefinitionReference /haveTypeRefAsParameter]] [[String hi]]]]]
  `)
}
}
*/

func TestTypeRefPipeRight(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
haveStringAsParameter : (String) -> Bool

another : (Bool) -> Bool =
    true

main : Bool =
    haveStringAsParameter "hi" |> another
`, `
[ModuleDef $haveStringAsParameter = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [ExternalFunctionDeclarationExpression [FnDeclExpr 0] [Primitive Any]]]]
[ModuleDef $another = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [Bool true]]]
[ModuleDef $main = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([]) -> [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] [FunctionRef [NamedDefinitionReference /haveStringAsParameter]] [[String hi]]] |> [FunctionRef [NamedDefinitionReference /another]]]]
`)
}

func TestTypeRefPipeRight2(t *testing.T) {
	testDecorate(t,
		`
type alias GoldPrice = Int


haveTypeRefAsParameter : (TypeRef a) -> Maybe a


main : Bool =
    haveTypeRefAsParameter $GoldPrice |> Maybe.withDefault false
`, `
[ModuleDef $haveTypeRefAsParameter = [FunctionValue [LocalTypeNameContext a = [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $TypeRef]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $TypeRef]]]]]) -> [ExternalFunctionDeclarationExpression [FnDeclExpr 0] [Primitive Any]]]]
[ModuleDef $another = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [Bool true]]]
[ModuleDef $main = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([]) -> [FnCall [FunctionType [[ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]] [Primitive Bool]]] [FunctionRef [NamedDefinitionReference /haveTypeRefAsParameter]] [[String hi]]] |> [FunctionRef [NamedDefinitionReference /another]]]]
`)
}

func TestArrayLiteral2(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
a : (Bool) -> Array Int =
    [| 10, 20, 99 |]
`, `
[ModuleDef $a = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $Array]] [LocalTypeNameOnlyContext a => [Primitive Array<[LocalTypeNameRef a]>]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [ArrayLiteral resolved context [[Integer 10] [Integer 20] [Integer 99]]]]]
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
Cool : [Alias Cool [RecordType [[Field $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]]]

[ModuleDef $a = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [ResolvedContext a:[AliasRef [Alias Cool [RecordType [[Field $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [ListLiteral [[RecordLiteral [RecordType [[Field $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]] [0 = [String hi]]] [RecordLiteral [RecordType [[Field $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]] [0 = [String another]]] [RecordLiteral [RecordType [[Field $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]] [0 = [String tjoho]]]]]]]
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
Cool : [Alias Cool [RecordType [[Field $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]

[ModuleDef $a = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [ResolvedContext a:[AliasRef [Alias Cool [RecordType [[Field $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [ListLiteral [[RecordConstructorRecord [CCall [TypeReference $Cool] [[RecordLiteral [[$name = #95]]]]] [RecordLiteral [RecordType [[Field $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]] [0 = [Integer 95]]]]]]]]
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
Cool : [Alias Cool [RecordType [[Field $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]

[ModuleDef $a = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [ResolvedContext a:[AliasRef [Alias Cool [RecordType [[Field $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [ListLiteral [[RecordConstructorRecord [CCall [TypeReference $Cool] [[RecordLiteral [[$name = #95]]]]] [RecordLiteral [RecordType [[Field $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]] [0 = [Integer 95]]]]]]]]
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
Cool : [Alias Cool [RecordType [[Field $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]

[ModuleDef $a = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [ResolvedContext a:[AliasRef [Alias Cool [RecordType [[Field $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [ListLiteral [[record-constructor [AliasRef [Alias Cool [RecordType [[Field $name [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]] [0 = [Integer 2]]]]]]]
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
CustomType : [CustomType CustomType [[Variant $First] [Variant $Second]]]
First : [Variant $First]
Second : [Variant $Second]

[ModuleDef $some = [FunctionValue [FunctionType [[CustomTypeRef [NamedDefTypeRef :[TypeReference $CustomType]] [CustomType CustomType [[Variant $First] [Variant $Second]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] ([[Param $a : [CustomTypeRef [NamedDefTypeRef :[TypeReference $CustomType]] [CustomType CustomType [[Variant $First] [Variant $Second]]]]]]) -> [PMCase: [ParamRef [Param $a : [CustomTypeRef [NamedDefTypeRef :[TypeReference $CustomType]] [CustomType CustomType [[Variant $First] [Variant $Second]]]]]] of  default: [String ]]]]
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
Tinkering : [Alias Tinkering [RecordType [[Field $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (0)]]]]
Studying : [Alias Studying [RecordType [[Field $subject [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]]]
Work : [Alias Work [RecordType [[Field $web [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]
Child : [CustomType Child [[Variant $Aron [[AliasRef [Alias Tinkering [RecordType [[Field $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (0)]]]]]]] [Variant $Alexandra] [Variant $Alma [[AliasRef [Alias Studying [RecordType [[Field $subject [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]]]]]] [Variant $Isabelle [[AliasRef [Alias Work [RecordType [[Field $web [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]]]]]]
Aron : [Variant $Aron [[AliasRef [Alias Tinkering [RecordType [[Field $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (0)]]]]]]]
Alexandra : [Variant $Alexandra]
Alma : [Variant $Alma [[AliasRef [Alias Studying [RecordType [[Field $subject [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]]]]]]
Isabelle : [Variant $Isabelle [[AliasRef [Alias Work [RecordType [[Field $web [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]]]]

[ModuleDef $some = [FunctionValue [FunctionType [[CustomTypeRef [NamedDefTypeRef :[TypeReference $Child]] [CustomType Child [[Variant $Aron [[AliasRef [Alias Tinkering [RecordType [[Field $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (0)]]]]]]] [Variant $Alexandra] [Variant $Alma [[AliasRef [Alias Studying [RecordType [[Field $subject [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]]]]]] [Variant $Isabelle [[AliasRef [Alias Work [RecordType [[Field $web [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]]]]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] ([[Param $child : [CustomTypeRef [NamedDefTypeRef :[TypeReference $Child]] [CustomType Child [[Variant $Aron [[AliasRef [Alias Tinkering [RecordType [[Field $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (0)]]]]]]] [Variant $Alexandra] [Variant $Alma [[AliasRef [Alias Studying [RecordType [[Field $subject [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]]]]]] [Variant $Isabelle [[AliasRef [Alias Work [RecordType [[Field $web [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]]]]]]]]]) -> [dcase: [ParamRef [Param $child : [CustomTypeRef [NamedDefTypeRef :[TypeReference $Child]] [CustomType Child [[Variant $Aron [[AliasRef [Alias Tinkering [RecordType [[Field $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (0)]]]]]]] [Variant $Alexandra] [Variant $Alma [[AliasRef [Alias Studying [RecordType [[Field $subject [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]]]]]] [Variant $Isabelle [[AliasRef [Alias Work [RecordType [[Field $web [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]]]]]]]]] of [dcasecons [VariantRef [NamedDefTypeRef :[TypeReference $Aron]] [Variant $Aron [[AliasRef [Alias Tinkering [RecordType [[Field $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (0)]]]]]]]] ([[dcaseparm $x type:[AliasRef [Alias Tinkering [RecordType [[Field $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (0)]]]]]]]) => [String Aron]];[dcasecons [VariantRef [NamedDefTypeRef :[TypeReference $Alexandra]] [Variant $Alexandra]] ([]) => [String Alexandris]] default: [String ]]]]
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
[ModuleDef $some = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [PMCase: [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]] of [PMCaseCons [String hello] => [Integer 0]] default: [Integer -1]]]]
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
Tinkering : [Alias Tinkering [LocalTypeNameOnlyContext t => [RecordType [[Field $secret [LocalTypeNameRef t] (0)] [Field $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]]]]]

[ModuleDef $f = [FunctionValue [FunctionType [[ResolvedContext t:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference <nil> [LocalTypeNameOnlyContext t => [RecordType [[Field $secret [LocalTypeNameRef t] (0)] [Field $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]]]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $tinkering : [ResolvedContext t:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference <nil> [LocalTypeNameOnlyContext t => [RecordType [[Field $secret [LocalTypeNameRef t] (0)] [Field $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]]]]]]]]) -> [lookups [ParamRef [Param $tinkering : [ResolvedContext t:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference <nil> [LocalTypeNameOnlyContext t => [RecordType [[Field $secret [LocalTypeNameRef t] (0)] [Field $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]]]]]]]] [[lookup [Field $secret [ConcreteGenericRef [GenericParam t] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]] (0)]]]]]]
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
Tinkering : [Alias Tinkering [LocalTypeNameOnlyContext t => [RecordType [[Field $secret [LocalTypeNameRef t] (0)] [Field $solder [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]]]]]
`)
}

func TestUnmanaged(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
type alias CAbstraction =
    { used : Bool
    , native : Unmanaged<CStruct>
    }
`, `
CAbstraction : [Alias CAbstraction [RecordType [[Field $native [unmanaged [Unmanaged $CStruct]] (0)] [Field $used [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] (1)]]]]
`)
}

// can not add a child that is not within the range of the parent [[2:5](23) to [3:5](167) (145)]  [[1:1](0) to [1:1](0) (1)]  (*dectype.RecordAtom and *dectype.PrimitiveAtom)
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
    true


updateWorld : (w: World) -> Bool =
    updateEnemy w.enemies


main : (Bool) -> Bool =
    updateWorld { enemies = [ { values = [ 1, 3 ] } ] }
`, `
Enemy : [Alias Enemy [RecordType [[Field $values [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]]]
World : [Alias World [RecordType [[Field $enemies [ResolvedContext a:[AliasRef [Alias Enemy [RecordType [[Field $values [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]]]

[ModuleDef $updateEnemy = [FunctionValue [FunctionType [[ResolvedContext a:[AliasRef [Alias Enemy [RecordType [[Field $values [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $_ : [ResolvedContext a:[AliasRef [Alias Enemy [RecordType [[Field $values [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]]) -> [Bool true]]]
[ModuleDef $updateWorld = [FunctionValue [FunctionType [[AliasRef [Alias World [RecordType [[Field $enemies [ResolvedContext a:[AliasRef [Alias Enemy [RecordType [[Field $values [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $w : [AliasRef [Alias World [RecordType [[Field $enemies [ResolvedContext a:[AliasRef [Alias Enemy [RecordType [[Field $values [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]]]]]]) -> [FnCall [FunctionType [[ResolvedContext a:[AliasRef [Alias Enemy [RecordType [[Field $values [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] [FunctionRef [NamedDefinitionReference /updateEnemy]] [[lookups [ParamRef [Param $w : [AliasRef [Alias World [RecordType [[Field $enemies [ResolvedContext a:[AliasRef [Alias Enemy [RecordType [[Field $values [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]]]]]] [[lookup [Field $enemies [ResolvedContext a:[AliasRef [Alias Enemy [RecordType [[Field $values [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]]]]]]]
[ModuleDef $main = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FnCall [FunctionType [[AliasRef [Alias World [RecordType [[Field $enemies [ResolvedContext a:[AliasRef [Alias Enemy [RecordType [[Field $values [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] [FunctionRef [NamedDefinitionReference /updateWorld]] [[RecordLiteral [RecordType [[Field $enemies [ResolvedContext a:[RecordType [[Field $values [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]] [0 = [ListLiteral [[RecordLiteral [RecordType [[Field $values [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]] [0 = [ListLiteral [[Integer 1] [Integer 3]]]]]]]]]]]]]
`)
}

func TestRecordListInList2(t *testing.T) {
	testDecorate(t,

		`
drawSprite : (Int) -> Bool =
    true


drawSprites : (sprites: List Int) -> List Bool =
    List.map drawSprite sprites
`, `
[ModuleDef $drawSprite = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Bool true]]]
[ModuleDef $drawSprites = [FunctionValue [FunctionType [[ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]] ([[Param $sprites : [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]]) -> [FnCall [FunctionType [[FunctionType [[Primitive Int] [Primitive Bool]]] [Primitive List<[ConcreteGenericRef [GenericParam a] => [ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]>] [Primitive List<[ConcreteGenericRef [GenericParam a] => [ConcreteGenericRef [GenericParam b] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]>]]] [FunctionRef [NamedDefinitionReference List/map]] [[FunctionRef [NamedDefinitionReference /drawSprite]] [ParamRef [Param $sprites : [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]]]]]]
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
    true


drawSprites : (sprites: List Sprite) -> List Bool =
    List.map drawSprite sprites


drawWorld : (world: World) -> List (List Bool) =
    List.map drawSprites world.drawTasks


main : (Bool) -> List (List Bool) =
    drawWorld { drawTasks = [ [ { x = 10, y = 20 }, { x = 44, y = 98 } ], [ { x = 99, y = 98 } ] ] }
`, `
Sprite : [Alias Sprite [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]
World : [Alias World [RecordType [[Field $drawTasks [ResolvedContext a:[ResolvedContext a:[AliasRef [Alias Sprite [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]]]

[ModuleDef $drawSprite = [FunctionValue [FunctionType [[AliasRef [Alias Sprite [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $_ : [AliasRef [Alias Sprite [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]]]]) -> [Bool true]]]
[ModuleDef $drawSprites = [FunctionValue [FunctionType [[ResolvedContext a:[AliasRef [Alias Sprite [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]] ([[Param $sprites : [ResolvedContext a:[AliasRef [Alias Sprite [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]]) -> [FnCall [FunctionType [[FunctionType [[RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]] [Primitive Bool]]] [Primitive List<[ConcreteGenericRef [GenericParam a] => [ConcreteGenericRef [GenericParam a] => [AliasRef [Alias Sprite [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]]]]>] [Primitive List<[ConcreteGenericRef [GenericParam a] => [ConcreteGenericRef [GenericParam b] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]>]]] [FunctionRef [NamedDefinitionReference List/map]] [[FunctionRef [NamedDefinitionReference /drawSprite]] [ParamRef [Param $sprites : [ResolvedContext a:[AliasRef [Alias Sprite [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]]]]]]
[ModuleDef $drawWorld = [FunctionValue [FunctionType [[AliasRef [Alias World [RecordType [[Field $drawTasks [ResolvedContext a:[ResolvedContext a:[AliasRef [Alias Sprite [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]]]] [ResolvedContext a:[ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]] ([[Param $world : [AliasRef [Alias World [RecordType [[Field $drawTasks [ResolvedContext a:[ResolvedContext a:[AliasRef [Alias Sprite [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]]]]]]) -> [FnCall [FunctionType [[FunctionType [[Primitive List<[ConcreteGenericRef [GenericParam a] => [AliasRef [Alias Sprite [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]]]>] [Primitive List<[ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]>]]] [Primitive List<[ConcreteGenericRef [GenericParam a] => [ConcreteGenericRef [GenericParam a] => [ResolvedContext a:[AliasRef [Alias Sprite [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]]>] [Primitive List<[ConcreteGenericRef [GenericParam a] => [ConcreteGenericRef [GenericParam b] => [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]]>]]] [FunctionRef [NamedDefinitionReference List/map]] [[FunctionRef [NamedDefinitionReference /drawSprites]] [lookups [ParamRef [Param $world : [AliasRef [Alias World [RecordType [[Field $drawTasks [ResolvedContext a:[ResolvedContext a:[AliasRef [Alias Sprite [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]]]]]] [[lookup [Field $drawTasks [ResolvedContext a:[ResolvedContext a:[AliasRef [Alias Sprite [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]]]]]]]
[ModuleDef $main = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [ResolvedContext a:[ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FnCall [FunctionType [[AliasRef [Alias World [RecordType [[Field $drawTasks [ResolvedContext a:[ResolvedContext a:[AliasRef [Alias Sprite [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]]]] [ResolvedContext a:[ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]] [FunctionRef [NamedDefinitionReference /drawWorld]] [[RecordLiteral [RecordType [[Field $drawTasks [ResolvedContext a:[ResolvedContext a:[RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] (0)]]] [0 = [ListLiteral [[ListLiteral [[RecordLiteral [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]] [0 = [Integer 10] 1 = [Integer 20]]] [RecordLiteral [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]] [0 = [Integer 44] 1 = [Integer 98]]]]] [ListLiteral [[RecordLiteral [RecordType [[Field $x [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $y [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]] [0 = [Integer 99] 1 = [Integer 98]]]]]]]]]]]]]
`)
}

func TestAppliedAnnotation(t *testing.T) {
	testDecorateWithoutDefault(t,

		`
type alias MyList t =
    { someType : t
    }


intConvert : (MyList Int) -> Bool
`, `
MyList : [Alias MyList [LocalTypeNameOnlyContext t => [RecordType [[Field $someType [LocalTypeNameRef t] (0)]]]]]

[ModuleDef $intConvert = [FunctionValue [FunctionType [[ResolvedContext t:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference <nil> [LocalTypeNameOnlyContext t => [RecordType [[Field $someType [LocalTypeNameRef t] (0)]]]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $_ : [ResolvedContext t:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference <nil> [LocalTypeNameOnlyContext t => [RecordType [[Field $someType [LocalTypeNameRef t] (0)]]]]]]]]) -> [ExternalFunctionDeclarationExpression [FnDeclExpr 0] [Primitive Any]]]]
`)
}

func TestAppliedAnnotation2(t *testing.T) {
	testDecorateWithoutDefault(t,

		`
type alias MyList t =
    { someType : t
    }


intConvert : (MyList Int) -> Bool =
    true
`, `
MyList : [Alias MyList [LocalTypeNameOnlyContext t => [RecordType [[Field $someType [LocalTypeNameRef t] (0)]]]]]

[ModuleDef $intConvert = [FunctionValue [FunctionType [[ResolvedContext t:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference <nil> [LocalTypeNameOnlyContext t => [RecordType [[Field $someType [LocalTypeNameRef t] (0)]]]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $_ : [ResolvedContext t:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference <nil> [LocalTypeNameOnlyContext t => [RecordType [[Field $someType [LocalTypeNameRef t] (0)]]]]]]]]) -> [Bool true]]]
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
[ModuleDef $simple = [FunctionValue [LocalTypeNameOnlyContext a => [FunctionType [[ResolvedContext a:[LocalTypeNameRef a] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] [LocalTypeNameRef a]]]] ([[Param $_ : [ResolvedContext a:[LocalTypeNameRef a] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]]) -> [Integer 2]]]
[ModuleDef $main = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FnCall [FunctionType [[Primitive List<[ConcreteGenericRef [GenericParam a] => [ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]>] [Primitive Int]]] [FunctionRef [NamedDefinitionReference /simple]] [[ListLiteral [[Integer 1] [Integer 2] [Integer 3]]]]]]]
`)
}

func TestFunctionTypeParameters(t *testing.T) {
	testDecorateWithoutDefault(t,

		`
intConvert : (List Int) -> Bool =
    true


simple : ((List a -> Bool), List a) -> Bool =
    true


main : (Bool) -> Bool =
    simple intConvert [ 1, 2, 3 ]
`, `
[ModuleDef $intConvert = [FunctionValue [FunctionType [[ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $_ : [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]]) -> [Bool true]]]
[ModuleDef $simple = [FunctionValue [LocalTypeNameOnlyContext a => [FunctionType [[FunctionType [[ResolvedContext a:[LocalTypeNameRef a] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] [ResolvedContext a:[LocalTypeNameRef a] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]] ([[Param $_ : [FunctionType [[ResolvedContext a:[LocalTypeNameRef a] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]] [Param $_ : [ResolvedContext a:[LocalTypeNameRef a] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]]) -> [Bool true]]]
[ModuleDef $main = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FnCall [FunctionType [[FunctionType [[Primitive List<[ConcreteGenericRef [GenericParam a] => [ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]>] [Primitive Bool]]] [Primitive List<[ConcreteGenericRef [GenericParam a] => [ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]>] [Primitive Bool]]] [FunctionRef [NamedDefinitionReference /simple]] [[FunctionRef [NamedDefinitionReference /intConvert]] [ListLiteral [[Integer 1] [Integer 2] [Integer 3]]]]]]]
`)
}

func xTestFunctionTypeGenerics(t *testing.T) {
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
[ModuleDef $intConvert = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Integer 42]]]
[ModuleDef $simple = [FunctionValue [LocalTypeNameContext a = [FunctionType [[LocalTypeNameRef a] [FunctionType [[LocalTypeNameRef a] [LocalTypeNameRef a]]]]]] ([[Param $_ : [LocalTypeNameRef a]]]) -> [FunctionRef [NamedDefinitionReference /intConvert]]]]
[ModuleDef $main = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [Let [[LetAssign [[LetVar $fn]] = [FnCall [FunctionType [[ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]] [FunctionType [[ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]] [ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]] [FunctionRef [NamedDefinitionReference /simple]] [[Integer 33]]]]] in [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [FunctionRef [NamedDefinitionReference /intConvert]] [[FnCall [FunctionType [[ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]] [ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]] [LetVarRef $fn] [[Integer 22]]]]]]]]
`)
}

func xTestFunctionTypeGenericsFail(t *testing.T) {
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
`, &decorated.FunctionArgumentTypeMismatch{})
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
Scale2 : [Alias Scale2 [RecordType [[Field $scaleX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $scaleY [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]
Sprite : [Alias Sprite [RecordType [[Field $dummy [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $scale [AliasRef [Alias Scale2 [RecordType [[Field $scaleX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $scaleY [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] (1)]]]]

[ModuleDef $updateSprite = [FunctionValue [FunctionType [[AliasRef [Alias Sprite [RecordType [[Field $dummy [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $scale [AliasRef [Alias Scale2 [RecordType [[Field $scaleX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $scaleY [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] (1)]]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [AliasRef [Alias Sprite [RecordType [[Field $dummy [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $scale [AliasRef [Alias Scale2 [RecordType [[Field $scaleX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $scaleY [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] (1)]]]]]]] ([[Param $inSprite : [AliasRef [Alias Sprite [RecordType [[Field $dummy [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $scale [AliasRef [Alias Scale2 [RecordType [[Field $scaleX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $scaleY [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] (1)]]]]]] [Param $newScale : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [RecordLiteral [RecordType [[Field $dummy [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $scale [AliasRef [Alias Scale2 [RecordType [[Field $scaleX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $scaleY [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] (1)]]] [1 = [RecordLiteral [RecordType [[Field $scaleX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $scaleY [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]] [0 = [ParamRef [Param $newScale : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] 1 = [ParamRef [Param $newScale : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]]]]
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
[ModuleDef $checkMaybeInt = [FunctionValue [FunctionType [[ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $Maybe]] [LocalTypeNameOnlyContext a => [CustomType Maybe [[Variant $Nothing] [Variant $Just [[LocalTypeNameRef a]]]]]]]] [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $Maybe]] [LocalTypeNameOnlyContext a => [CustomType Maybe [[Variant $Nothing] [Variant $Just [[LocalTypeNameRef a]]]]]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $a : [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $Maybe]] [LocalTypeNameOnlyContext a => [CustomType Maybe [[Variant $Nothing] [Variant $Just [[LocalTypeNameRef a]]]]]]]]] [Param $b : [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $Maybe]] [LocalTypeNameOnlyContext a => [CustomType Maybe [[Variant $Nothing] [Variant $Just [[LocalTypeNameRef a]]]]]]]]]]) -> [dcase: [ParamRef [Param $a : [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $Maybe]] [LocalTypeNameOnlyContext a => [CustomType Maybe [[Variant $Nothing] [Variant $Just [[LocalTypeNameRef a]]]]]]]]]] of [dcasecons [VariantRef [NamedDefTypeRef :[TypeReference $Just]] [Variant $Just [[ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]] ([[dcaseparm $firstInt type:[ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]) => [dcase: [ParamRef [Param $b : [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $Maybe]] [LocalTypeNameOnlyContext a => [CustomType Maybe [[Variant $Nothing] [Variant $Just [[LocalTypeNameRef a]]]]]]]]]] of [dcasecons [VariantRef [NamedDefTypeRef :[TypeReference $Nothing]] [Variant $Nothing]] ([]) => [Integer 0]];[dcasecons [VariantRef [NamedDefTypeRef :[TypeReference $Just]] [Variant $Just [[ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]] ([[dcaseparm $secondInt type:[ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]) => [functionparamref $secondInt [dcaseparm $secondInt type:[ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]];[dcasecons [VariantRef [NamedDefTypeRef :[TypeReference $Nothing]] [Variant $Nothing]] ([]) => [Integer 0]]]]]
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
Thing : [Alias Thing [RecordType [[Field $something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]]]

[ModuleDef $main = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [ResolvedContext a:[AliasRef [Alias Thing [RecordType [[Field $something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $Maybe]] [LocalTypeNameOnlyContext a => [CustomType Maybe [[Variant $Nothing] [Variant $Just [[LocalTypeNameRef a]]]]]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [VariantConstructor [Variant $Nothing] []]]]
`)
}

func TestOwnMaybe(t *testing.T) {
	testDecorateWithoutDefault(t, `
type OwnMaybe a =
    None
    | Value a`, `
OwnMaybe : [LocalTypeNameOnlyContext a => [CustomType OwnMaybe [[Variant $None] [Variant $Value [[LocalTypeNameRef a]]]]]]
None : [Variant $None]
Value : [Variant $Value [[LocalTypeNameRef a]]]
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
            false

        Right ->
            false

        Down ->
            true

        Left ->
            2

        Up ->
            true
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
            false

        Right ->
            false

        Down ->
            true

        Left ->
            false

        Up ->
            true

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
            false

        Right ->
            false

        Down ->
            true

        Left ->
            false

        Up ->
            true

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
Thing : [Alias Thing [RecordType [[Field $something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]]]

[ModuleDef $main = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [ResolvedContext a:[AliasRef [Alias Thing [RecordType [[Field $something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $Maybe]] [LocalTypeNameOnlyContext a => [CustomType Maybe [[Variant $Nothing] [Variant $Just [[LocalTypeNameRef a]]]]]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [VariantConstructor [Variant $Just [[LocalTypeNameRef a]]] [[RecordLiteral [RecordType [[Field $something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]] [0 = [String hello]]]]]]]
`)
}

func TestJustAndNothing(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
type alias Thing =
    { something : String
    }


main : (Bool) -> Maybe Thing =
    if true then
        Nothing
    else
        Just { something = "hello" }
`, `
Thing : [Alias Thing [RecordType [[Field $something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]]]

[ModuleDef $main = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [ResolvedContext a:[AliasRef [Alias Thing [RecordType [[Field $something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $Maybe]] [LocalTypeNameOnlyContext a => [CustomType Maybe [[Variant $Nothing] [Variant $Just [[LocalTypeNameRef a]]]]]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [If [Bool true] then [VariantConstructor [Variant $Nothing] []] else [VariantConstructor [Variant $Just [[LocalTypeNameRef a]]] [[RecordLiteral [RecordType [[Field $something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]] [0 = [String hello]]]]]]]]
`)
}

func TestCustomTypeGenerics2(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
get : (String) -> Maybe Int

`, `
[ModuleDef $get = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $Maybe]] [LocalTypeNameOnlyContext a => [CustomType Maybe [[Variant $Nothing] [Variant $Just [[LocalTypeNameRef a]]]]]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [ExternalFunctionDeclarationExpression [FnDeclExpr 0] [Primitive Any]]]]
`)
}

func xTestOwnGenericGet(t *testing.T) {
	testDecorateWithoutDefault(t,
		`

type alias Thing a =
    { something : a
    }


get : (Thing a) -> Maybe a


main : (Bool) -> Maybe (Thing String) =
    let
        v = { something = "Something" }
        _ = get v
    in
    Just v
        
`, `
Thing : [Alias Thing [LocalTypeNameContext a = [RecordType [[Field $something [LocalTypeNameRef a] (0)]]]]]

[ModuleDef $get = [FunctionValue [LocalTypeNameContext a = [FunctionType [[AliasRef [Alias Thing [LocalTypeNameContext a = [RecordType [[Field $something [LocalTypeNameRef a] (0)]]]]]] [CustomTypeRef [NamedDefTypeRef :[TypeReference $Maybe]] [CustomType Maybe [[ConcreteGenericRef [GenericParam a] => [LocalTypeNameRef a]]] [[Variant $Nothing] [Variant $Just [[ConcreteGenericRef [GenericParam a] => [LocalTypeNameRef a]]]]]]]]]] ([[Param $_ : [AliasRef [Alias Thing [LocalTypeNameContext a = [RecordType [[Field $something [LocalTypeNameRef a] (0)]]]]]]]]) -> [ExternalFunctionDeclarationExpression [FnDeclExpr 0] [Primitive Any]]]]
[ModuleDef $main = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [CustomTypeRef [NamedDefTypeRef :[TypeReference $Maybe]] [CustomType Maybe [[ConcreteGenericRef [GenericParam a] => [AliasRef [Alias Thing [LocalTypeNameContext a = [RecordType [[Field $something [LocalTypeNameRef a] (0)]]]]]]]] [[Variant $Nothing] [Variant $Just [[ConcreteGenericRef [GenericParam a] => [AliasRef [Alias Thing [LocalTypeNameContext a = [RecordType [[Field $something [LocalTypeNameRef a] (0)]]]]]]]]]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [Let [[LetAssign [[LetVar $v]] = [RecordLiteral [RecordType [[Field $something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]] [0 = [String Something]]]] [LetAssign [[LetVar $_]] = [FnCall [FunctionType [[RecordType [[Field $something [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] (0)]]] [CustomType Maybe [[ConcreteGenericRef [GenericParam a] => [ConcreteGenericRef [GenericParam a] => [LocalTypeNameRef a]]]] [[Variant $Nothing] [Variant $Just [[ConcreteGenericRef [GenericParam a] => [ConcreteGenericRef [GenericParam a] => [LocalTypeNameRef a]]]]]]]]] [FunctionRef [NamedDefinitionReference /get]] [[LetVarRef $v]]]]] in [VariantConstructor [Variant $Just [[LocalTypeNameRef a]]] [[LetVarRef $v]]]]]]
`)
}

func xTestArrayGetWithMaybe(t *testing.T) {
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
PlayerAction : [CustomType PlayerAction [[Variant $None] [Variant $Jump]]]
None : [Variant $None]
Jump : [Variant $Jump]
Gamepad : [Alias Gamepad [RecordType [[Field $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]
UserInput : [Alias UserInput [RecordType [[Field $gamepads [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Array]]] (0)]]]]
PlayerInputs : [Alias PlayerInputs [RecordType [[Field $inputs [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]] (0)]]]]

[ModuleDef $checkMaybeGamepad = [FunctionValue [FunctionType [[CustomTypeRef [NamedDefTypeRef :[TypeReference $Maybe]] [CustomType Maybe [[ConcreteGenericRef [GenericParam a] => [AliasRef [Alias Gamepad [RecordType [[Field $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]]]] [[Variant $Nothing] [Variant $Just [[ConcreteGenericRef [GenericParam a] => [AliasRef [Alias Gamepad [RecordType [[Field $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]]]]]]]] [CustomTypeRef [NamedDefTypeRef :[TypeReference $PlayerAction]] [CustomType PlayerAction [[Variant $None] [Variant $Jump]]]]]] ([[Param $_ : [CustomTypeRef [NamedDefTypeRef :[TypeReference $Maybe]] [CustomType Maybe [[ConcreteGenericRef [GenericParam a] => [AliasRef [Alias Gamepad [RecordType [[Field $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]]]] [[Variant $Nothing] [Variant $Just [[ConcreteGenericRef [GenericParam a] => [AliasRef [Alias Gamepad [RecordType [[Field $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]]]]]]]]]]) -> [VariantConstructor [Variant $None] []]]]
[ModuleDef $main = [FunctionValue [FunctionType [[AliasRef [Alias UserInput [RecordType [[Field $gamepads [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Array]]] (0)]]]]] [AliasRef [Alias PlayerInputs [RecordType [[Field $inputs [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]] (0)]]]]]]] ([[Param $oldUserInputs : [AliasRef [Alias UserInput [RecordType [[Field $gamepads [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Array]]] (0)]]]]]]]) -> [Let [[LetAssign [[LetVar $gamepads]] = [lookups [ParamRef [Param $oldUserInputs : [AliasRef [Alias UserInput [RecordType [[Field $gamepads [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Array]]] (0)]]]]]]] [[lookup [Field $gamepads [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Array]]] (0)]]]]] [LetAssign [[LetVar $maybeOld]] = [FnCall [FunctionType [[Primitive Int] [Primitive Array<[ConcreteGenericRef [GenericParam a] => [AliasRef [Alias Gamepad [RecordType [[Field $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]]]>] [CustomType Maybe [[ConcreteGenericRef [GenericParam a] => [ConcreteGenericRef [GenericParam a] => [LocalTypeNameRef a]]]] [[Variant $Nothing] [Variant $Just [[ConcreteGenericRef [GenericParam a] => [ConcreteGenericRef [GenericParam a] => [LocalTypeNameRef a]]]]]]]]] [FunctionRef [NamedDefinitionReference Array/get]] [[Integer 0] [LetVarRef $gamepads]]]] [LetAssign [[LetVar $playerAction]] = [FnCall [FunctionType [[CustomTypeRef [NamedDefTypeRef :[TypeReference $Maybe]] [CustomType Maybe [[ConcreteGenericRef [GenericParam a] => [AliasRef [Alias Gamepad [RecordType [[Field $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]]]] [[Variant $Nothing] [Variant $Just [[ConcreteGenericRef [GenericParam a] => [AliasRef [Alias Gamepad [RecordType [[Field $a [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)]]]]]]]]]]] [CustomTypeRef [NamedDefTypeRef :[TypeReference $PlayerAction]] [CustomType PlayerAction [[Variant $None] [Variant $Jump]]]]]] [FunctionRef [NamedDefinitionReference /checkMaybeGamepad]] [[LetVarRef $maybeOld]]]]] in [RecordLiteral [RecordType [[Field $inputs [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $List]]] (0)]]] [0 = [ListLiteral [[LetVarRef $playerAction]]]]]]]]
`)
}

func xTestArrayIntGetWithMaybe(t *testing.T) {
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
[ModuleDef $getInt = [FunctionValue [FunctionType [[CustomTypeRef [NamedDefTypeRef :[TypeReference $Maybe]] [CustomType Maybe [[ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [[Variant $Nothing] [Variant $Just [[ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [CustomTypeRef [NamedDefTypeRef :[TypeReference $Maybe]] [CustomType Maybe [[ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [[Variant $Nothing] [Variant $Just [[ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]]]]) -> [Integer 22]]]
[ModuleDef $main = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Array]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $ints : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Array]]]]]) -> [Let [[LetAssign [[LetVar $maybeInt]] = [FnCall [FunctionType [[Primitive Int] [Primitive Array<[ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]>] [CustomType Maybe [[ConcreteGenericRef [GenericParam a] => [ConcreteGenericRef [GenericParam a] => [LocalTypeNameRef a]]]] [[Variant $Nothing] [Variant $Just [[ConcreteGenericRef [GenericParam a] => [ConcreteGenericRef [GenericParam a] => [LocalTypeNameRef a]]]]]]]]] [FunctionRef [NamedDefinitionReference Array/get]] [[Integer 0] [ParamRef [Param $ints : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Array]]]]]]]] [LetAssign [[LetVar $extractedInt]] = [FnCall [FunctionType [[CustomTypeRef [NamedDefTypeRef :[TypeReference $Maybe]] [CustomType Maybe [[ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [[Variant $Nothing] [Variant $Just [[ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [FunctionRef [NamedDefinitionReference /getInt]] [[LetVarRef $maybeInt]]]]] in [LetVarRef $extractedInt]]]]
`)
}

func xTestArrayIntGetWithMaybeFail(t *testing.T) {
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

func xTestArrayFromNothing(t *testing.T) {
	testDecorate(t,
		`
needInputs : (Array Int) -> Int =
    2


a : (Bool) -> Int =
    needInputs (Array.fromList [])
`,
		`
[ModuleDef $needInputs = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Array]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Array]]]]]) -> [Integer 2]]]
[ModuleDef $a = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Array]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [FunctionRef [NamedDefinitionReference /needInputs]] [[FnCall [FunctionType [[Primitive List<[Primitive Any]>] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Array]]]]] [FunctionRef [NamedDefinitionReference Array/fromList]] [[ListLiteral []]]]]]]]
`)
}

func TestReturnSimpleFunction(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
needInputs : (String) -> Int =
    2


returnSomeFunc : (Bool) -> (String -> Int) =
    needInputs


a : (Bool) -> Int =
    let
        fn = returnSomeFunc false
    in
    fn "hello world"
`,
		`
[ModuleDef $needInputs = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [Integer 2]]]
[ModuleDef $returnSomeFunc = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FunctionRef [NamedDefinitionReference /needInputs]]]]
[ModuleDef $a = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [Let [[LetAssign [[LetVar $fn]] = [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]] [FunctionRef [NamedDefinitionReference /returnSomeFunc]] [[Bool false]]]]] in [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [LetVarRef $fn] [[String hello world]]]]]]
`)
}

func TestReturnFunction(t *testing.T) {
	testDecorateWithoutDefault(t,
		`
needInputs : (List Int) -> Int =
    2


returnSomeFunc : (Bool) -> (List Int -> Int) =
    needInputs


a : (Bool) -> Int =
    let
        fn = returnSomeFunc false
    in
    fn []
`,
		`
[ModuleDef $needInputs = [FunctionValue [FunctionType [[ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]]]]) -> [Integer 2]]]
[ModuleDef $returnSomeFunc = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [FunctionType [[ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [FunctionRef [NamedDefinitionReference /needInputs]]]]
[ModuleDef $a = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]) -> [Let [[LetAssign [[LetVar $fn]] = [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] [FunctionType [[ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]]] [FunctionRef [NamedDefinitionReference /returnSomeFunc]] [[Bool false]]]]] in [FnCall [FunctionType [[ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $List]] [LocalTypeNameOnlyContext a => [Primitive List<[LocalTypeNameRef a]>]]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [LetVarRef $fn] [[ListLiteral []]]]]]]
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
Something : [Alias Something [RecordType [[Field $playerX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $time [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]

[ModuleDef $fn = [FunctionValue [FunctionType [[AliasRef [Alias Something [RecordType [[Field $playerX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $time [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]] [AliasRef [Alias Something [RecordType [[Field $playerX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $time [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]]]] ([[Param $something : [AliasRef [Alias Something [RecordType [[Field $playerX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $time [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]]]]]]) -> [RecordLiteral [RecordType [[Field $playerX [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (0)] [Field $time [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] (1)]]] [0 = [Integer 3]]]]]
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
    first (2 + 2)
        |> second b
        |> third
`, `
[ModuleDef $first = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Integer 2]]]
[ModuleDef $second = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([[Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]] [Param $_ : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [Integer 2]]]
[ModuleDef $third = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]) -> [BoolOp [ParamRef [Param $a : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] GR [Integer 45]]]]
[ModuleDef $tester = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([[Param $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]) -> [FnCall [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] [FunctionRef [NamedDefinitionReference /first]] [(Arithmetic [Integer 2] PLUS [Integer 2])]] |> [Curry [FunctionRef [NamedDefinitionReference /second]] [[ParamRef [Param $b : [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $String]]]]]]] |> [FunctionRef [NamedDefinitionReference /third]]]]
`)
}

func TestPipeRight2Fail(t *testing.T) {
	testDecorateWithoutDefaultFail(t, `
first : (Int) -> Int =
    2


second : (Int, Int) -> Int =
    2


third : (a: Int) -> Bool =
    a > 45


tester : (b: String) -> Bool =
    first (2 + 2)
        |> second b
        |> third
`, &decorated.FunctionArgumentTypeMismatch{})
}

func TestTypeRefGeneric(t *testing.T) {
	testDecorateWithoutDefault(t, `
getComponent : (TypeRef a) -> Maybe a

main : Bool =
    let
        maybeBool = getComponent $Bool
    in
    case maybeBool of
        Just a -> a

        Nothing -> false
`, `
[ModuleDef $getComponent = [FunctionValue [LocalTypeNameOnlyContext a => [FunctionType [[ResolvedContext a:[LocalTypeNameRef a] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $TypeRef]] [LocalTypeNameOnlyContext a => [Primitive TypeRef<[LocalTypeNameRef a]>]]]] [ResolvedContext a:[LocalTypeNameRef a] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $Maybe]] [LocalTypeNameOnlyContext a => [CustomType Maybe [[Variant $Nothing] [Variant $Just [[LocalTypeNameRef a]]]]]]]]]]] ([[Param $_ : [ResolvedContext a:[LocalTypeNameRef a] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $TypeRef]] [LocalTypeNameOnlyContext a => [Primitive TypeRef<[LocalTypeNameRef a]>]]]]]]) -> [ExternalFunctionDeclarationExpression [FnDeclExpr 0] [Primitive Any]]]]
[ModuleDef $main = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]] ([]) -> [Let [[LetAssign [[LetVar $maybeBool]] = [FnCall [FunctionType [[Primitive TypeRef<[ConcreteGenericRef [GenericParam a] => [ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]>] [CustomType Maybe [[Variant $Nothing] [Variant $Just [[ConcreteGenericRef [GenericParam a] => [ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]]]]]]] [FunctionRef [NamedDefinitionReference /getComponent]] [[TypeIdLiteral [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $TypeRef]] [LocalTypeNameOnlyContext a => [Primitive TypeRef<[LocalTypeNameRef a]>]]]]]]]]] in [dcase: [LetVarRef $maybeBool] of [dcasecons [VariantRef [NamedDefTypeRef :[TypeReference $Just]] [Variant $Just [[ConcreteGenericRef [GenericParam a] => [ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]]]] ([[dcaseparm $a type:[ConcreteGenericRef [GenericParam a] => [ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]]]) => [functionparamref $a [dcaseparm $a type:[ConcreteGenericRef [GenericParam a] => [ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Bool]]]]]]]];[dcasecons [VariantRef [NamedDefTypeRef :[TypeReference $Nothing]] [Variant $Nothing]] ([]) => [Bool false]]]]]]
`)
}

func TestTypeRefGeneric2(t *testing.T) {
	testDecorateWithoutDefault(t, `
getComponent : (TypeRef a) -> a


main : Int =
    getComponent $Int
`, `
[ModuleDef $getComponent = [FunctionValue [LocalTypeNameOnlyContext a => [FunctionType [[ResolvedContext a:[LocalTypeNameRef a] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $TypeRef]] [LocalTypeNameOnlyContext a => [Primitive TypeRef<[LocalTypeNameRef a]>]]]] [LocalTypeNameRef a]]]] ([[Param $_ : [ResolvedContext a:[LocalTypeNameRef a] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $TypeRef]] [LocalTypeNameOnlyContext a => [Primitive TypeRef<[LocalTypeNameRef a]>]]]]]]) -> [ExternalFunctionDeclarationExpression [FnDeclExpr 0] [Primitive Any]]]]
[ModuleDef $main = [FunctionValue [FunctionType [[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]] ([]) -> [FnCall [FunctionType [[Primitive TypeRef<[ConcreteGenericRef [GenericParam a] => [ConcreteGenericRef [GenericParam a] => [PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]]]]>] [Primitive Int]]] [FunctionRef [NamedDefinitionReference /getComponent]] [[TypeIdLiteral [ResolvedContext a:[PrimitiveTypeRef [NamedDefTypeRef :[TypeReference $Int]]] => [LocalTypeNameOnlyContextReference [NamedDefTypeRef :[TypeReference $TypeRef]] [LocalTypeNameOnlyContext a => [Primitive TypeRef<[LocalTypeNameRef a]>]]]]]]]]]
`)
}
