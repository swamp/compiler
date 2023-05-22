/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dtype

import (
	"github.com/swamp/compiler/src/token"
)

type Type interface {
	HumanReadable() string
	String() string
	Resolve() (Atom, error)
	Next() Type
	FetchPositionLength() token.SourceFileReference
	WasReferenced() bool
}
