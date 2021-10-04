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
