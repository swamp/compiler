/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"
	"github.com/swamp/compiler/src/decorated/dtype"
)

type TypeParameterContext struct {
	argumentNames     []*dtype.TypeArgumentName
	resolvedArguments []dtype.Type
}

func (t *TypeParameterContext) DeclareString() string {
	if len(t.resolvedArguments) == 0 {
		return ""
	}
	return fmt.Sprintf("%v", t.resolvedArguments)
}

func (t *TypeParameterContext) String() string {
	return fmt.Sprintf("[typeparamcontext %v (%v)]", t.resolvedArguments, t.argumentNames)
}

func (t *TypeParameterContext) DebugString() string {
	s := ""
	for index, name := range t.argumentNames {
		argumentType := t.resolvedArguments[index]
		if index > 0 {
			s += "\n"
		}
		s += fmt.Sprintf(" %v = %v", name, argumentType)
	}
	return s
}

func (t *TypeParameterContext) DecoratedName() string {
	s := ""
	for index, arg := range t.resolvedArguments {
		if index > 0 {
			s += ","
		}
		if arg == nil {
			panic("illegal state in " + t.String())
		}
		s += arg.DecoratedName()
	}
	return s
}

func (t *TypeParameterContext) ArgumentNamesCount() int {
	return len(t.argumentNames)
}

func (t *TypeParameterContext) ArgumentNames() []*dtype.TypeArgumentName {
	return t.argumentNames
}

func (t *TypeParameterContext) ArgumentTypes() []dtype.Type {
	return t.resolvedArguments
}

type DecoratedTypeError interface {
	Error() string
}

func NewTypeParameterContext(names []*dtype.TypeArgumentName, resolvedArguments []dtype.Type) (*TypeParameterContext, DecoratedTypeError) {
	if len(names) != len(resolvedArguments) {
		panic(fmt.Sprintf("mismatch type parameter context count"))
		return nil, fmt.Errorf("mismatch type parameter context count %v vs %v", names, resolvedArguments)
	}

	return &TypeParameterContext{argumentNames: names, resolvedArguments: resolvedArguments}, nil
}

func (t *TypeParameterContext) LookupTypeFromName(name string) dtype.Type {
	for index, foundParam := range t.argumentNames {
		if foundParam.Name() == name {
			return t.resolvedArguments[index]
		}
	}
	return nil
}

func (t *TypeParameterContext) LookupTypeFromArgument(param *dtype.TypeArgumentName) dtype.Type {
	return t.LookupTypeFromName(param.Name())
}


