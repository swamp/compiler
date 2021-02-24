/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func decorateRecordLookups(d DecorateStream, lookups *ast.Lookups, context *VariableContext) (*decorated.RecordLookups, decshared.DecoratedError) {
	expressionToLookup := context.ResolveVariable(lookups.ContextIdentifier())
	if expressionToLookup == nil {
		return nil, decorated.NewCouldNotFindIdentifierInLookups(lookups)
	}

	var lookupFields []decorated.LookupField

	typeToLookup := dectype.Unalias(expressionToLookup.Type())
	for _, lookupIdentifier := range lookups.FieldNames() {
		recordTypeToCheck, lookupErr := dectype.ResolveToRecordType(typeToLookup)
		if lookupErr != nil {
			fmt.Printf("this is not a record!_!!? %T %v\n\n%v\n", typeToLookup, typeToLookup, lookups)
			return nil, decorated.NewUnMatchingTypes(nil, nil, nil, nil)
		}

		nameToLookup := lookupIdentifier.Name()
		foundRecordTypeField := recordTypeToCheck.FindField(nameToLookup)
		if foundRecordTypeField == nil {
			return nil, decorated.NewCouldNotFindFieldInLookup(lookups, lookupIdentifier, typeToLookup)
		}

		fakeLookupField := decorated.NewLookupField(foundRecordTypeField)
		lookupFields = append(lookupFields, fakeLookupField)
		typeToLookup = dectype.Unalias(foundRecordTypeField.Type())
	}

	if len(lookupFields) == 0 {
		panic("must have at least one lookup to be valid")
	}

	return decorated.NewRecordLookups(expressionToLookup, lookupFields), nil
}
