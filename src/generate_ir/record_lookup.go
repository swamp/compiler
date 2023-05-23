/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_ir

import (
	"fmt"
	"log"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func handleRecordLookup(lookups *decorated.RecordLookups,
	genContext *generateContext) (value.Value, error) {
	structPtr, err := generateExpression(lookups.Expression(), false, genContext)
	if err != nil {
		return nil, err
	}

	if !types.IsPointer(structPtr.Type()) {
		panic(fmt.Errorf("something went wrong %v", structPtr))
	}

	log.Printf("structPtr: %v", structPtr)

	for lookupIndex, field := range lookups.LookupFields() {
		var indices []value.Value
		if lookupIndex == 0 {
			//indices = append(indices, constant.NewInt(types.I32, int64(0)))
		}
		indices = append(indices, constant.NewInt(types.I32, int64(field.Index())))
		irFieldType, err := genContext.irTypeRepo.GetTypeRef(field.RecordTypeFieldReference().Type())
		if err != nil {
			return nil, err
		}
		//irFieldTypePtr := types.NewPointer(irFieldType)
		log.Printf("expected type %v", irFieldType)

		log.Printf(
			"#%d lookup %v source:%v indices:%v raw:%v", lookupIndex, irFieldType, structPtr, indices,
			lookups.LookupFields(),
		)
		structPtr = ir.NewGetElementPtr(irFieldType, structPtr, indices...)
	}

	lastField := lookups.LookupFields()[len(lookups.LookupFields())-1]
	lastIrFieldType, lastErr := genContext.irTypeRepo.GetTypeRef(lastField.RecordTypeFieldReference().Type())
	if lastErr != nil {
		return nil, lastErr
	}

	loadInstruction := ir.NewLoad(lastIrFieldType, structPtr)

	return loadInstruction, nil
}

func generateLookups(lookups *decorated.RecordLookups,
	genContext *generateContext) (value.Value, error) {
	sourcePosRange, err := handleRecordLookup(lookups, genContext)
	if err != nil {
		return nil, err
	}

	return sourcePosRange, nil
}
