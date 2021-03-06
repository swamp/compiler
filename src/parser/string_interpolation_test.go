/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"testing"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/token"
)

func TestInterpolation(t *testing.T) {
	shouldBeFalse, _ := isInterpolation("${x  test")
	if shouldBeFalse {
		t.Errorf("should not be true")
	}

	shouldBeTrue, _ := isInterpolation("something${dfko} x")
	if !shouldBeTrue {
		t.Errorf("should be true")
	}
}

func TestInterpolationRanges(t *testing.T) {
	ranges := parseInterpolationString("sdfj ${xid}")
	if len(ranges) != 1 {
		t.Errorf("wrong range count")
	}

	if ranges[0][0] != 5 || ranges[0][1] != 11 {
		t.Errorf("wrong range")
	}
}

func TestInterpolationSubstitution(t *testing.T) {
	replaced := replaceInterpolationString(`some ${a 23 b}after`)

	if replaced != `"some " ++ Debug.toString(a 23 b) ++ "after"` {
		t.Errorf("wrong replaced %v", replaced)
	}
}

func TestInterpolationSubstitution2(t *testing.T) {
	replaced := replaceInterpolationString(`some ${a 23 b}after ${b}`)

	const correct = `"some " ++ Debug.toString(a 23 b) ++ "after " ++ Debug.toString(b)`
	if replaced != correct {
		t.Errorf("wrong replaced. expected \n%v\nbut received\n%v\n", correct, replaced)
	}
}

func replaceInterpolationStringToExpressionHelper(t *testing.T, s string) *ast.StringInterpolation {
	expr, err := replaceInterpolationStringToExpression(token.NewStringToken(s, s, token.SourceFileReference{}))
	if err != nil {
		t.Fatal(err)
	}

	return expr
}

func TestInterpolationSubstitutionExpression(t *testing.T) {
	expr := replaceInterpolationStringToExpressionHelper(t, `some ${a 23 b}after ${b}`)
	if expr.String() != `((('some ' ++ [call Debug.$toString [[call $a [#23 $b]]]]) ++ 'after ') ++ [call Debug.$toString [$b]])` {
		t.Errorf("wrong expression")
	}
}

func TestInterpolationStringExpression2(t *testing.T) {
	expr := replaceInterpolationStringToExpressionHelper(t, `${a} ${b}`)
	if expr.String() != `((('' ++ [call Debug.$toString [$a]]) ++ ' ') ++ [call Debug.$toString [$b]])` {
		t.Errorf("wrong expression")
	}
}
