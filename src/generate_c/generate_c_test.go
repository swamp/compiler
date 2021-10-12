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
main x =
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
            x = "Rebecca"
            y = let
                a = 3
            in
            a
            z = 99
        in
        y
    else
        let
            x = "3"
        in
        33

`, `

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
