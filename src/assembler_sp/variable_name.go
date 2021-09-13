/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import "fmt"

type VariableName struct {
	name string
}

func NewVariableName(name string) *VariableName {
	return &VariableName{name: name}
}

func (o *VariableName) Name() string {
	return o.name
}

func (o *VariableName) String() string {
	return fmt.Sprintf("[var %v]", o.name)
}
