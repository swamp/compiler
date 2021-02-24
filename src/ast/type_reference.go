/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type TypeReference struct {
	ident     *TypeIdentifier
	arguments []Type
}

func (i *TypeReference) String() string {
	if len(i.arguments) == 0 {
		return fmt.Sprintf("[type-reference %v]", i.ident)
	}
	return fmt.Sprintf("[type-reference %v %v]", i.ident, i.arguments)
}

func (i *TypeReference) TypeResolver() *TypeIdentifier {
	return i.ident
}

func (i *TypeReference) Arguments() []Type {
	return i.arguments
}

func (i *TypeReference) Name() string {
	s := ""
	if len(i.arguments) == 0 {
		return fmt.Sprintf("%v", i.ident.Name())
	}

	for index, argument := range i.arguments {
		if index > 0 {
			s += " "
		}
		s += argument.Name()
	}
	return fmt.Sprintf("%v<%v>", i.ident.Name(), s)
}

func (i *TypeReference) FetchPositionLength() token.Range {
	return i.ident.FetchPositionLength()
}
