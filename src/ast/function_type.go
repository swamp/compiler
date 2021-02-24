/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"github.com/swamp/compiler/src/token"
)

type FunctionType struct {
	functionParameters []Type
}

func (i *FunctionType) Name() string {
	s := ""
	for index, arg := range i.functionParameters {
		if index > 0 {
			s += " -> "
		}
		s += arg.Name()
	}

	return s
}

func (i *FunctionType) FunctionParameters() []Type {
	return i.functionParameters
}

func (i *FunctionType) String() string {
	s := "[func-type "
	for index, arg := range i.functionParameters {
		if index > 0 {
			s += " -> "
		}
		s += arg.String()
	}
	s += "]"
	return s
}

func (i *FunctionType) FetchPositionLength() token.SourceFileReference {
	return i.functionParameters[0].FetchPositionLength()
}

func NewFunctionType(functionParameters []Type) *FunctionType {
	return &FunctionType{functionParameters: functionParameters}
}
