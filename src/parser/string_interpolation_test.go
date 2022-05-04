/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"testing"
)

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
	replaced := replaceInterpolationStringToString(`some ${a 23 b}after`)

	if replaced != `"some " ++ Debug.toString(a 23 b) ++ "after"` {
		t.Errorf("wrong replaced %v", replaced)
	}
}

func TestInterpolationSubstitution2(t *testing.T) {
	replaced := replaceInterpolationStringToString(`some {a 23 b}after {b}`)

	const correct = `"some " ++ Debug.toString(a 23 b) ++ "after " ++ Debug.toString(b)`
	if replaced != correct {
		t.Errorf("wrong replaced. expected \n%v\nbut received\n%v\n", correct, replaced)
	}
}
