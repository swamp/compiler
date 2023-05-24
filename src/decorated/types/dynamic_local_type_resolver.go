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

/*

func (t *DynamicLocalTypeResolver) LookupTypeFromName(name string) dtype.Type {
	if len(t.argumentNames) == 0 {
		panic("strange, you can not lookup if it is empty")
	}
	for index, foundParam := range t.argumentNames {
		if foundParam.Name() == name {
			existing := t.localTypeNamesMap[index]
			if existing == nil {
				panic(fmt.Errorf("how can existing be nil '%v' %v", name, foundParam))
			}

			return existing
		}
	}
	return nil
}

func (t *DynamicLocalTypeResolver) LookupType(name string) (dtype.Type, error) {
	foundType := t.LookupTypeFromName(name)
	if foundType == nil {
		log.Printf("couldn't find '%v'\n%v", name, t.DebugString())
		return nil, fmt.Errorf("couldn't find the name %v", name)
	}

	return foundType, nil
}

*/

func (t *DynamicLocalTypeResolver) SpecialSet(name string, resolved dtype.Type) error {
	if IsLocalType(resolved) || IsAny(resolved) {
		panic(fmt.Errorf("must be set with a concrete type %T %v", resolved, resolved))
	}
	for index, foundParam := range t.argumentNames {
		if foundParam.Name() == name {
			existing := t.resolvedArguments[index]
			if existing != nil {
				compatibleErr := CompatibleTypes(existing, resolved)
				if compatibleErr != nil {
					log.Printf("not compatible!!!\n")
					return fmt.Errorf("it didn't work %w", compatibleErr)
				}
			}

			log.Printf("%s <- %T %v", name, resolved, resolved)
			t.resolvedArguments[index] = resolved

			return nil
		}
	}

	return fmt.Errorf("couldn't find %v", name)
}

func (t *DynamicLocalTypeResolver) SetType(defRef *LocalTypeName, definedType dtype.Type) error {
	return t.SpecialSet(defRef.Name(), definedType)
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
