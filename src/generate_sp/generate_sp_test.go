/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_sp

import (
	"testing"
)

func TestIntEqual(t *testing.T) {
	testGenerateWithoutCores(t,
		`
isCold : (temp: Int) -> Bool =
    temp == -1
`, `
func [constantfn DynPos 0008:104 func:isCold]
0000: ldi 8,-1
0009: cpeqi 0,4,8
0016: ret
`)
}

func TestBooleanOperator(t *testing.T) {
	testGenerateWithoutCores(t,
		`
main : (Int) -> Bool =
    let
        a = True
        b = False
    in
    a && b
`, `
func [constantfn DynPos 0008:104 func:main]
0000: ldb 8,true
0006: ldb 9,false
000c: cpy 0,(8:1)
0017: brfa 0 [label @0029]
001e: cpy 0,(9:1)
0029: ret
`)
}

func TestUnary2(t *testing.T) {
	testGenerateWithoutCores(t,
		`
someTest : (a: Bool, b: Bool) -> Bool =
    !a && b
`, `
func [constantfn DynPos 0010:104 func:someTest]
0000: not 0,1
0009: brfa 0 [label @001b]
0010: cpy 0,(2:1)
001b: ret
`)
}

func TestListLiteral(t *testing.T) {
	testGenerateWithoutCores(t,
		`
type alias Cool =
    { name : String
    }


a : (Bool) -> List Cool =
    [ { name = "hi" }, { name = "another" }, { name = "tjoho" } ]
`, `
[constantstring DynPos 0073:16 hi]
[constantstring DynPos 008B:16 another]
[constantstring DynPos 00A1:16 tjoho]
func [constantfn DynPos 0008:104 func:a]
0000: ldz 16,$0073
0009: ldz 24,$008B
0012: ldz 32,$00A1
001b: crl 0 [16 24 32] (8, 8)
0030: ret
`)
}

func TestArrayLiteral(t *testing.T) {
	testGenerateWithoutCores(t,
		`
type alias Cool =
    { name : String
    }


a : (Bool) -> Array Cool =
    [| { name = "hello" }, { name = "world" }, { name = "ossian" } |]
`, `
[constantstring DynPos 0076:16 hello]
[constantstring DynPos 008C:16 world]
[constantstring DynPos 00A3:16 ossian]
func [constantfn DynPos 0008:104 func:a]
0000: ldz 16,$0076
0009: ldz 24,$008C
0012: ldz 32,$00A3
001b: cra 0 [16 24 32] (8, 8)
0030: ret
`)
}

func TestCustomTypeConstructor(t *testing.T) {
	testGenerateWithoutCores(t,
		`
type SomeEnum =
    First String
    | Anon
    | Second Int


a : (Bool) -> SomeEnum =
    First "Hello"
`, `
[constantstring DynPos 0076:16 Hello]
func [constantfn DynPos 0008:104 func:a]
0000: lde 0,0 (16)
0008: ldz 8,$0076
0011: ret
`)
}

func TestCustomTypeVariantConstructorSecond(t *testing.T) {
	testGenerateWithoutCores(t,
		`
type SomeEnum =
    First String
    | Anon
    | Second Int


a : (Bool) -> SomeEnum =
    Second 1
`, `
func [constantfn DynPos 0008:104 func:a]
0000: lde 0,2 (16)
0008: ldi 4,1
0011: ret
`)
}

func TestMaybeInt(t *testing.T) {
	testGenerateWithoutCores(t,
		`
main : (Bool) -> Maybe Int =
    Just 3
`, `
func [constantfn DynPos 0008:104 func:main]
0000: lde 0,1 (8)
0008: ldi 4,3
0011: ret
`)
}

func TestIfStatement(t *testing.T) {
	testGenerateWithoutCores(t,
		`
main : (name: String) -> Bool =
    if name == "Rebecca" then
        let
            _ = "Rebecca"
        in
        True
    else
        let
            _ = "3"
        in
        False

`, `
[constantstring DynPos 0078:16 Rebecca]
[constantstring DynPos 008A:16 3]
func [constantfn DynPos 0008:104 func:main]
0000: ldz 24,$0078
0009: cpeqs 16,8,24
0016: brfa 16 [label @002f]
001d: ldz 32,$0078
0026: ldb 0,true
002c: jmp [label @003e]
002f: ldz 40,$008A
0038: ldb 0,false
003e: ret
`)
}

func TestCurry(t *testing.T) {
	testGenerateWithoutCores(t,
		`
isWinner : (name: String, score: Int) -> Bool =
    if name == "Ossian" then
        score * 2 > 100
    else
        score > 100


main : (score: Int) -> Bool =
    let
        checkScoreFn = isWinner "Ossian"
    in
    checkScoreFn score
`, `
[constantstring DynPos 00EF:16 Ossian]
func [constantfn DynPos 0010:104 func:isWinner]
0000: ldz 24,$00EF
0009: cpeqs 20,8,24
0016: brfa 20 [label @004c]
001d: ldi 36,2
0026: muli 32,16,36
0033: ldi 40,100
003c: cpgti 0,32,40
0049: jmp [label @0062]
004c: ldi 44,100
0055: cpgti 0,16,44
0062: ret

func [constantfn DynPos 0080:104 func:main]
0000: ldz 16,$0010
0009: ldz 32,$00EF
0012: curry 8,16,(32:8) (typeId:4, align:8)
0024: cpy 20,(4:4)
002f: call 16 8
0038: cpy 0,(16:1)
0043: ret
`)
}

func TestAppend(t *testing.T) {
	testGenerateWithoutCores(t,
		`
a : (Int) -> List Int =
    [ 1, 3, 4 ] ++ [ 5, 6, 7, 8 ] ++ [ 9 ]
`, `
func [constantfn DynPos 0008:104 func:a]
0000: ldi 32,1
0009: ldi 36,3
0012: ldi 40,4
001b: crl 24 [32 36 40] (4, 4)
0030: ldi 56,5
0039: ldi 60,6
0042: ldi 64,7
004b: ldi 68,8
0054: crl 48 [56 60 64 68] (4, 4)
006d: listappend 16,24,48
007a: ldi 80,9
0083: crl 72 [80] (4, 4)
0090: listappend 0,16,72
009d: ret
`)
}

func TestGuardLetInChar(t *testing.T) {
	testGenerateWithoutCores(t,
		`
tester : (Int) -> Char =
    let
        existingTile = 'a'
        isUpperLeft = False
    in
    | existingTile == '_' -> '@'
    | isUpperLeft -> '/'
    | _ -> '2'
`, `
func [constantfn DynPos 0008:104 func:tester]
0000: ldr 8,'a' (97)
0006: ldb 12,false
000c: ldr 16,'_' (95)
0012: cpeqi 13,8,16
001f: brfa 13 [label @002f]
0026: ldr 0,'@' (64)
002c: jmp [label @0045]
002f: brfa 12 [label @003f]
0036: ldr 0,'/' (47)
003c: jmp [label @0045]
003f: ldr 0,'2' (50)
0045: ret
`)
}

func TestCasePatternMatchingString(t *testing.T) {
	testGenerateWithoutCores(t,
		`
some : (a: Int) -> Int =
    case a of
        2 -> 0

        3 -> 1

        _ -> -1
	`, `
func [constantfn DynPos 0008:104 func:some]
0000: jmppmi 4 [[2 [label @0014]] [3 [label offset @0020]]] [label offset @002c]
0014: ldi 0,0
001d: jmp [label @0035]
0020: ldi 0,1
0029: jmp [label @0035]
002c: ldi 0,-1
0035: ret
`)
}
