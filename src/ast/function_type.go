/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"github.com/swamp/compiler/src/token"
)

type FunctionType struct {
	functionParameters  []Type
	sourceFileReference token.SourceFileReference
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
	s := "[FnType "
	for index, arg := range i.functionParameters {
		if index > 0 {
			if index == len(i.functionParameters)-1 {
				s += " -> "
			} else {
				s += ", "
			}
		}
		s += arg.String()
	}
	s += "]"
	return s
}

func (i *FunctionType) FetchPositionLength() token.SourceFileReference {
	return i.sourceFileReference
}

func NewFunctionType(functionParameters []Type) *FunctionType {
	first := functionParameters[0].FetchPositionLength()
	last := functionParameters[len(functionParameters)-1].FetchPositionLength()
	sourceFileReference := token.MakeInclusiveSourceFileReference(first, last)
	return &FunctionType{functionParameters: functionParameters, sourceFileReference: sourceFileReference}
}
