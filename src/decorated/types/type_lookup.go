/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"

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

	equalErr := pureExpected.IsEqual(pureActual)
	if equalErr != nil {
		return fmt.Errorf("*** NOT EQUAL:\n %v\nvs\n %v\n %w", pureExpected.AtomName(), pureActual.AtomName(), equalErr)
	}

	return nil
}

func CompatibleTypes(expectedType dtype.Type, actualType dtype.Type) error {
	if expectedType == nil {
		panic(fmt.Sprintf("shouldn't happen. expected is nil, actualType is %v", actualType))
	}

	if actualType == nil {
		panic(fmt.Sprintf("shouldn't happen. actualType is nil, expectedType is %v", expectedType))
	}

	customType, wasCustomType := expectedType.(*CustomTypeAtom)
	if wasCustomType {
		otherVariant, wasVariant := actualType.(*CustomTypeVariant)
		if wasVariant {
			return customType.IsVariantEqual(otherVariant)
		}
	}

	pureExpected, expectedErr := expectedType.Resolve()
	pureActual, actualErr := actualType.Resolve()

	isAny := IsAtomAny(pureActual)
	if isAny {
		return nil
	}

	isExpectedAny := IsAtomAny(pureExpected)
	if isExpectedAny {
		return nil
	}

	if actualErr != nil {
		return actualErr
	}
	if expectedErr != nil {
		return expectedErr
	}

	return compareAtoms(pureExpected, pureActual)
}

func ResolveToRecordType(expectedRecord dtype.Type) (*RecordAtom, error) {
	atom, atomErr := expectedRecord.Resolve()
	if atomErr != nil {
		return nil, fmt.Errorf("couldn't resolve to record %w", atomErr)
	}

	recordAtom, wasRecord := atom.(*RecordAtom)
	if !wasRecord {
		return nil, fmt.Errorf("resolved to something else than a record %v", atom)
	}

	return recordAtom, nil
}
