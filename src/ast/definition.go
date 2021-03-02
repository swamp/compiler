/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type FunctionValueNamedDefinition struct {
	identifier    *VariableIdentifier
	functionValue *FunctionValue
}

func NewFunctionValueNamedDefinition(identifier *VariableIdentifier, functionValue *FunctionValue) *FunctionValueNamedDefinition {
	return &FunctionValueNamedDefinition{identifier: identifier, functionValue: functionValue}
}

func (i *FunctionValueNamedDefinition) Identifier() *VariableIdentifier {
	return i.identifier
}

func (i *FunctionValueNamedDefinition) FunctionValue() *FunctionValue {
	return i.functionValue
}

func (i *FunctionValueNamedDefinition) FetchPositionLength() token.SourceFileReference {
	return i.identifier.FetchPositionLength()
}

func (i *FunctionValueNamedDefinition) String() string {
	return fmt.Sprintf("[fndefinition: %v = %v]", i.identifier, i.functionValue)
}

func (i *FunctionValueNamedDefinition) DebugString() string {
	return "[fndefinition]"
}
