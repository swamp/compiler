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

type TypeParameterContextOther struct {
	resolvedArguments map[string]dtype.Type
}

func (t *TypeParameterContextOther) DeclareString() string {
	if len(t.resolvedArguments) == 0 {
		return ""
	}
	return fmt.Sprintf("%v", t.resolvedArguments)
}

func (t *TypeParameterContextOther) String() string {
	return fmt.Sprintf("[typeparamcontext %v]", t.resolvedArguments)
}

func (t *TypeParameterContextOther) DebugString() string {
	s := ""
	for name, argumentType := range t.resolvedArguments {
		s += fmt.Sprintf("%v = %v", name, argumentType)
	}
	return s
}

func NewTypeParameterContextOther() *TypeParameterContextOther {
	t := &TypeParameterContextOther{resolvedArguments: make(map[string]dtype.Type)}

	return t
}

func (t *TypeParameterContextOther) LookupTypeFromName(name string) dtype.Type {
	return t.resolvedArguments[name]
}

func (t *TypeParameterContextOther) LookupType(name string) (dtype.Type, error) {
	found := t.LookupTypeFromName(name)
	if found == nil {
		log.Printf("%p couldn't find '%v' count:%d \n%v", t, name, len(t.resolvedArguments), t.DebugString())
		return nil, fmt.Errorf("couldn't find '%v' name", name)
	}

	return found, nil
}

func (t *TypeParameterContextOther) SpecialSet(name string, resolved dtype.Type) (dtype.Type, error) {
	existing := t.resolvedArguments[name]
	if existing != nil {
		compatibleErr := CompatibleTypes(existing, resolved)
		if compatibleErr != nil {
			return nil, fmt.Errorf("it didn't work %w", compatibleErr)
		}

		return existing, nil
	}

	t.resolvedArguments[name] = resolved

	return resolved, nil
}

func (t *TypeParameterContextOther) Verify() error {
	for _, x := range t.resolvedArguments {
		if x == nil {
			return fmt.Errorf("TypeParameterContextOther:Verify. Argument name %v has not been resolved", t.resolvedArguments)
		}
	}

	return nil
}

func (t *TypeParameterContextOther) LookupTypeFromArgument(param *dtype.TypeArgumentName) dtype.Type {
	return t.LookupTypeFromName(param.Name())
}
