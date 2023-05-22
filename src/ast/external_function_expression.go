/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type ExternalFunctionExpression struct {
	posLen token.SourceFileReference `debug:"true"`
}

func NewExternalFunctionExpression(posLen token.SourceFileReference) *ExternalFunctionExpression {
	return &ExternalFunctionExpression{posLen: posLen}
}

func (d *ExternalFunctionExpression) FetchPositionLength() token.SourceFileReference {
	return d.posLen
}

func (d *ExternalFunctionExpression) String() string {
	return fmt.Sprintf("[externalfunc: %v]", d.posLen)
}

func (d *ExternalFunctionExpression) DebugString() string {
	return d.String()
}
