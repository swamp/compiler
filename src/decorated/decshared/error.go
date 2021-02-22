/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decshared

import "github.com/swamp/compiler/src/token"

type DecoratedError interface {
	FetchPositionLength() token.Range
	Error() string
}
