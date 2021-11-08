/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"log"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func decorateRecordLookups(d DecorateStream, lookups *ast.Lookups, context *VariableContext) (*decorated.RecordLookups, decshared.DecoratedError) {
	expressionToLookup, expressionLookUpErr := context.ResolveVariable(lookups.ContextIdentifier())
	if expressionLookUpErr != nil {
		return nil, expressionLookUpErr
	}
	if expressionToLookup == nil {
		return nil, decorated.NewCouldNotFindIdentifierInLookups(lookups)
	}

	var lookupFields []decorated.LookupField

	aliasedTypeToLookup := expressionToLookup.Type()
	typeToLookup := dectype.Unalias(aliasedTypeToLookup)
	for _, lookupIdentifier := range lookups.FieldNames() {
		recordTypeToCheck, lookupErr := dectype.ResolveToRecordType(typeToLookup)
		if lookupErr != nil {
			log.Panicf("can not resolve this to a record: %v\n%v", expressionToLookup.FetchPositionLength().ToCompleteReferenceString(), typeToLookup)
			return nil, decorated.NewUnMatchingTypes(nil, nil, nil, lookupErr)
		}

		nameToLookup := lookupIdentifier.Name()
		foundRecordTypeField := recordTypeToCheck.FindField(nameToLookup)
		if foundRecordTypeField == nil {
			return nil, decorated.NewCouldNotFindFieldInLookup(lookups, lookupIdentifier, typeToLookup)
		}

		recordFieldReference := decorated.NewRecordFieldReference(lookupIdentifier, aliasedTypeToLookup, recordTypeToCheck, foundRecordTypeField)
		fakeLookupField := decorated.NewLookupField(recordFieldReference)
		lookupFields = append(lookupFields, fakeLookupField)
		aliasedTypeToLookup = foundRecordTypeField.Type()
		typeToLookup = dectype.Unalias(aliasedTypeToLookup)
	}

	if len(lookupFields) == 0 {
		panic("must have at least one lookup to be valid")
	}

	return decorated.NewRecordLookups(expressionToLookup, lookupFields), nil
}
