/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"
	"reflect"

	"github.com/swamp/compiler/src/decorated/dtype"
)

type TypeReferenceScopedOrNormal interface {
	dtype.Type
	NameReference() *NamedDefinitionTypeReference
}

func compareAtoms(pureExpected dtype.Atom, pureActual dtype.Atom) error {
	expectedIsAny := IsAtomAny(pureExpected)
	actualIsAny := IsAtomAny(pureActual)

	if expectedIsAny || actualIsAny {
		return nil
	}

	if pureExpected == nil || pureActual == nil {
		return fmt.Errorf("can not have nil stuff here")
	}

	if reflect.TypeOf(pureExpected) == nil {
		panic(fmt.Errorf("pureExpected is nil"))
	}

	equalErr := pureExpected.IsEqual(pureActual)
	if equalErr != nil {
		return fmt.Errorf("*** NOT EQUAL:\n %v\nvs\n %v\n %w", pureExpected.AtomName(), pureActual.AtomName(), equalErr)
	}

	return nil
}

func CompatibleTypesCheckCustomType(expectedType dtype.Type, actualType dtype.Type) error {
	if expectedType == nil {
		panic(fmt.Sprintf("shouldn't happen. expected is nil, actualType is %v", actualType))
	}

	if actualType == nil {
		panic(fmt.Sprintf("shouldn't happen. actualType is nil, expectedType is %v", expectedType))
	}

	pureExpected, expectedErr := expectedType.Resolve()
	pureActual, actualErr := actualType.Resolve()

	if pureExpected == nil || pureActual == nil {
		panic("error")
	}

	if expectedErr == nil && actualErr == nil {
		expectedVariant, wasExpectedVariant := pureExpected.(*CustomTypeVariantAtom)
		if wasExpectedVariant {
			actualVariant, wasActualVariant := pureActual.(*CustomTypeVariantAtom)
			if wasActualVariant {
				return CompatibleTypes(actualVariant.inCustomType, expectedVariant.inCustomType)
			}
			actualCustom, wasActualCustom := pureActual.(*CustomTypeAtom)
			if wasActualCustom {
				actualFoundVariant := actualCustom.FindVariant(expectedVariant.Name().Name())
				return CompatibleTypes(actualFoundVariant, expectedVariant)
			}
		}
	}

	return CompatibleTypes(expectedType, actualType)
}

func CompatibleTypes(expectedType dtype.Type, actualType dtype.Type) error {
	if expectedType == nil {
		panic(fmt.Sprintf("shouldn't happen. expected is nil, actualType is %v", actualType))
	}

	if actualType == nil {
		panic(fmt.Sprintf("shouldn't happen. actualType is nil, expectedType is %v", expectedType))
	}

	unaliasExpectedType := Unalias(expectedType)

	expectedLocalTypeNameOnlyContext, wasLocalTypeNameContext := unaliasExpectedType.(*LocalTypeNameOnlyContext)
	if wasLocalTypeNameContext {
		return CompatibleTypes(expectedLocalTypeNameOnlyContext.Next(), actualType)
	}

	unaliasActualType := Unalias(actualType)
	actualLocalTypeNameOnlyContext, wasActualLocalTypeNameContext := unaliasActualType.(*LocalTypeNameOnlyContext)
	if wasActualLocalTypeNameContext {
		return CompatibleTypes(expectedType, actualLocalTypeNameOnlyContext.Next())
	}

	pureExpected, expectedErr := expectedType.Resolve()
	pureActual, actualErr := actualType.Resolve()

	if pureActual == nil {
		panic(fmt.Errorf("pureActual is nil"))
	}

	isAny := IsAtomAny(pureActual)
	if isAny {
		return nil
	}

	isExpectedAny := IsAtomAny(pureExpected)
	if isExpectedAny {
		return nil
	}

	if actualErr != nil {
		panic(fmt.Errorf("pureActual can not be resolved: %w", actualErr))
		return actualErr
	}
	if expectedErr != nil {
		panic(fmt.Errorf("pureExpected can not be resolved: %w", expectedErr))
		return expectedErr
	}

	err := compareAtoms(pureExpected, pureActual)

	return err
}

func ResolveToRecordType(expectedRecord dtype.Type) (*RecordAtom, error) {
	atom := UnaliasWithResolveInvoker(expectedRecord)

	recordAtom, wasRecord := atom.(*RecordAtom)
	if !wasRecord {
		return nil, fmt.Errorf("resolved to something else than a record %T", atom)
	}

	return recordAtom, nil
}
