/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_c

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

`)
}

func TestIfStatement(t *testing.T) {
	testGenerate(t,
		`
main : String -> Int
main name =
    if name == "Rebecca" then
        let
            y = let
                a = 3
            in
            a
            z = 99
        in
        y + z
    else
        let
            _ = "3"
        in
        33

`, `

`)
}

func TestGuardLetInChar(t *testing.T) {
	testGenerate(t,
		`
tester : Int -> Char
tester _ =
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
