/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"testing"

	parerr "github.com/swamp/compiler/src/parser/errors"
)

func TestCommentBeforeAnnotation(t *testing.T) {
	testParse(t,
		`
{- This should be ignored -}
hello : Int -> Int -> Int
hello first c =
    4 + first
`,
		`
[annotation: $hello [func-type [type-reference $Int] -> [type-reference $Int] -> [type-reference $Int]]]
[definition: $hello = [func ([$first $c]) -> (#4 + $first)]]
`)
}

func TestDisallowCommentAfterAnnotation(t *testing.T) {
	testParseError(t,
		`
hello : Int -> Int -> Int {- This should not be allowed -}
hello first c =
    4 + first
`, parerr.InternalError{})
}

func TestAllowCommentsAfterRecordTypeMembers(t *testing.T) {
	testParse(t,
		`
type alias Struct =
    { a : Int -- Explaining what this field is
    , b : Boolean
    }
`, `
[alias $Struct [record-type [[field: $a [type-reference $Int]] [field: $b [type-reference $Boolean]]]]]
`)
}

func TestAllowCommentsAfterRecordTypeMembersSecond(t *testing.T) {
	testParse(t,
		`
type alias Struct =
    { a : Int -- Explaining what this field is
    , b : Boolean -- And what this is
    }
`, `
[alias $Struct [record-type [[field: $a [type-reference $Int]] [field: $b [type-reference $Boolean]]]]]
`)
}

func xTestAllowCommentsEndOfFile(t *testing.T) {
	testParse(t,
		`
type alias Struct =
    { a : Int -- Explaining what this field is
    , b : Boolean -- And what this is
    }


{--
so cool
--}
`, `
[alias $Struct [record-type [[field: $a [type-reference $Int]] [field: $b [type-reference $Boolean]]]]]
`)
}
