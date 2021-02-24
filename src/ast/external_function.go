/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type ExternalFunction struct {
	tokens         string
	parameterCount uint
	posLen         token.SourceFileReference
}

func NewExternalFunction(tokens string, parameterCount uint, posLen token.SourceFileReference) *ExternalFunction {
	return &ExternalFunction{tokens: tokens, parameterCount: parameterCount, posLen: posLen}
}

func (d *ExternalFunction) ExternalFunction() string {
	return d.tokens
}

func (d *ExternalFunction) ParameterCount() uint {
	return d.parameterCount
}

func (d *ExternalFunction) FetchPositionLength() token.SourceFileReference {
	return d.posLen
}

func (d *ExternalFunction) String() string {
	return fmt.Sprintf("[external function: %v %d]", d.ExternalFunction(), d.ParameterCount())
}

func (d *ExternalFunction) DebugString() string {
	return d.String()
}
