/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/dtype"
)

type TypeParameterContextDynamic struct {
	argumentNames     []*dtype.TypeArgumentName
	resolvedArguments []dtype.Type
}

func (t *TypeParameterContextDynamic) DeclareString() string {
	if len(t.resolvedArguments) == 0 {
		return ""
	}
	return fmt.Sprintf("%v", t.resolvedArguments)
}

func (t *TypeParameterContextDynamic) String() string {
	return fmt.Sprintf("[typeparamcontext %v (%v)]", t.resolvedArguments, t.argumentNames)
}

func (t *TypeParameterContextDynamic) DebugString() string {
	s := ""
	for index, name := range t.ArgumentNames() {
		argumentType := t.ArgumentTypes()[index]
		if index > 0 {
			s += ", "
		}
		s += fmt.Sprintf("%v = %v", name, argumentType)
	}
	return s
}

func (t *TypeParameterContextDynamic) DecoratedName() string {
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

func (t *TypeParameterContextDynamic) ArgumentNamesCount() int {
	return len(t.argumentNames)
}

func (t *TypeParameterContextDynamic) ArgumentNames() []*dtype.TypeArgumentName {
	return t.argumentNames
}

func (t *TypeParameterContextDynamic) ArgumentTypes() []dtype.Type {
	return t.resolvedArguments
}

func NewTypeParameterContextDynamic(names []*dtype.TypeArgumentName) *TypeParameterContextDynamic {
	return &TypeParameterContextDynamic{argumentNames: names, resolvedArguments: make([]dtype.Type, len(names))}
}

func (t *TypeParameterContextDynamic) ToParameterContext() (*TypeParameterContext, DecoratedTypeError) {
	if verifyErr := t.Verify(); verifyErr != nil {
		return nil, verifyErr
	}

	return NewTypeParameterContext(t.argumentNames, t.resolvedArguments)
}

func (t *TypeParameterContextDynamic) FillOutTheRestWithAny() {
	for index, x := range t.resolvedArguments {
		if x == nil {
			t.resolvedArguments[index] = NewAnyType()
		}
	}
}

func (t *TypeParameterContextDynamic) LookupTypeFromName(name string) dtype.Type {
	if len(t.argumentNames) == 0 {
		panic("strange, you can not lookup if it is empty")
	}
	for index, foundParam := range t.argumentNames {
		if foundParam.Name() == name {
			existing := t.resolvedArguments[index]
			if existing == nil {
				panic("how can existing be nil")
			}

			return existing
		}
	}
	return nil
}

func (t *TypeParameterContextDynamic) LookupType(name string) (dtype.Type, error) {
	foundType := t.LookupTypeFromName(name)
	if foundType == nil {
		fmt.Printf("couldn't find '%v'\n%v", name, t.DebugString())
		return nil, fmt.Errorf("couldn't find the name %v", name)
	}

	return foundType, nil
}

func (t *TypeParameterContextDynamic) SpecialSet(name string, resolved dtype.Type) error {
	for index, foundParam := range t.argumentNames {
		if foundParam.Name() == name {
			existing := t.resolvedArguments[index]
			if existing != nil {
				compatibleErr := CompatibleTypes(existing, resolved)
				if compatibleErr != nil {
					fmt.Printf("not compatible!!!\n")
					return fmt.Errorf("it didn't work %w", compatibleErr)
				}
			}

			t.resolvedArguments[index] = resolved

			return nil
		}
	}

	return fmt.Errorf("couldn't find %v", name)
}

func (t *TypeParameterContextDynamic) Verify() error {
	for index, x := range t.resolvedArguments {
		if x == nil {
			return fmt.Errorf("TypeParameterContextDynamic:Verify. Argument name %v has not been resolved",
				t.argumentNames[index])
		}
	}

	return nil
}

func (t *TypeParameterContextDynamic) LookupTypeFromArgument(param *dtype.TypeArgumentName) dtype.Type {
	return t.LookupTypeFromName(param.Name())
}
