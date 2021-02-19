/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser_test

import (
	"strings"
	"testing"

	"github.com/swamp/compiler/src/coloring"
	deccy "github.com/swamp/compiler/src/decorated"
	"github.com/swamp/compiler/src/parser"
)

func testColor(t *testing.T, code string) {
	code = strings.TrimSpace(code)
	const useCores = true
	module, compileErr := deccy.CompileToModuleOnceForTest(code, useCores, false)
	if compileErr != nil {
		t.Fatal(compileErr)
	}

	colorer := coloring.NewColorerWithColor()
	defs := module.Definitions()
	for _, expr := range defs.Definitions() {
		parser.ColorType(expr.Expression().Type(), 0, false, colorer)
	}

	for _, definedType := range module.TypeRepo().AllLocalTypes() {
		parser.ColorType(definedType, 0, false, colorer)
	}
}

func TestColor(t *testing.T) {
	testColor(t, `
type alias Cell =
    { x : Int
    , name : String
    , z : Maybe Int
    }


test : Array Cell -> Cell
test cells =
    let
        count = Array.length cells
    in
    { x = 10, name = "hello", z = Nothing }
`)
}
