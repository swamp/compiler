/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"testing"
)

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

func TestAllowCommentsEndOfFile(t *testing.T) {
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
