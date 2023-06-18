/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"
	"log"

	"github.com/swamp/compiler/src/decorated/dtype"
)

type DynamicResolver interface {
	SetType(defRef *LocalTypeName, definedType dtype.Type) error
}

type DynamicLocalTypeResolver struct {
	argumentNames     []*dtype.LocalTypeName
	resolvedArguments []dtype.Type
}

func (t *DynamicLocalTypeResolver) DeclareString() string {
	if len(t.resolvedArguments) == 0 {
		return ""
	}
	return fmt.Sprintf("%v", t.resolvedArguments)
}

func (t *DynamicLocalTypeResolver) String() string {
	return fmt.Sprintf("[typeparamcontext %v (%v)]", t.resolvedArguments, t.argumentNames)
}

func (t *DynamicLocalTypeResolver) DebugString() string {
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

func (t *DynamicLocalTypeResolver) IsDefined() bool {
	for _, def := range t.resolvedArguments {
		if def == nil {
			return false
		}
	}

	return true
}

func (t *DynamicLocalTypeResolver) DebugAllNotDefined() string {
	s := ""
	for index, def := range t.resolvedArguments {
		if def == nil {
			if len(s) != 0 {
				s += "\n"
			}
			s += t.argumentNames[index].String()
		}
	}

	return s
}

func (t *DynamicLocalTypeResolver) ArgumentNamesCount() int {
	return len(t.argumentNames)
}

func (t *DynamicLocalTypeResolver) ArgumentNames() []*dtype.LocalTypeName {
	return t.argumentNames
}

func (t *DynamicLocalTypeResolver) ArgumentTypes() []dtype.Type {
	return t.resolvedArguments
}

func NewDynamicLocalTypeResolver(names []*dtype.LocalTypeName) *DynamicLocalTypeResolver {
	return &DynamicLocalTypeResolver{argumentNames: names, resolvedArguments: make([]dtype.Type, len(names))}
}

func (t *DynamicLocalTypeResolver) FillOutTheRestWithAny() {
	for index, x := range t.resolvedArguments {
		if x == nil {
			t.resolvedArguments[index] = NewAnyType()
		}
	}
}

func (t *DynamicLocalTypeResolver) specialSet(name string, resolved dtype.Type) error {
	log.Printf("setting '%s' to %v", name, resolved)

	for index, foundParam := range t.argumentNames {
		if foundParam.Name() == name {
			existing := t.resolvedArguments[index]
			if existing != nil { /*
					compatibleErr := CompatibleTypes(existing, resolved)
					if compatibleErr != nil {
						log.Printf("not compatible!!! %v %s\n%v\n%v", compatibleErr,
							resolved.FetchPositionLength().ToCompleteReferenceString(), existing, resolved)
						return fmt.Errorf("it didn't work %w", compatibleErr)
					}*/

				return nil
			}

			t.resolvedArguments[index] = resolved

			return nil
		}
	}

	return fmt.Errorf("couldn't find %v", name)
}

func (t *DynamicLocalTypeResolver) SetType(defRef *LocalTypeName, definedType dtype.Type) error {
	return t.specialSet(defRef.Name(), definedType)
}

func (t *DynamicLocalTypeResolver) Verify() error {
	for index, x := range t.resolvedArguments {
		if x == nil {
			return fmt.Errorf("DynamicLocalTypeResolver:Verify. Argument name %v has not been resolved",
				t.argumentNames[index])
		}
	}

	return nil
}

/*
func (t *DynamicLocalTypeResolver) LookupTypeFromArgument(param *dtype.LocalTypeName) dtype.Type {
	return t.LookupTypeFromName(param.Name())
}
*/
