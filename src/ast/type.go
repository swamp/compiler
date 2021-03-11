/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"github.com/swamp/compiler/src/token"
)

type Type interface {
	Name() string
	String() string
	FetchPositionLength() token.SourceFileReference
}
