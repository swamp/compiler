/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"github.com/swamp/compiler/src/ast"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"log"
)

func DecorateRecordType(info *ast.Record, t decorated.TypeAddAndReferenceMaker) (*dectype.RecordAtom, decorated.TypeError) {
	var convertedFields []*dectype.RecordField
	for _, field := range info.Fields() {
		convertedFieldType, convertedFieldTypeErr := ConvertFromAstToDecorated(field.Type(), t)
		if convertedFieldTypeErr != nil {
			return nil, convertedFieldTypeErr
		}
		log.Printf("encountered %T %v", convertedFieldType, convertedFieldType.FetchPositionLength().ToCompleteReferenceString())
		fieldName := dectype.NewRecordFieldName(field.VariableIdentifier())
		convertedField := dectype.NewRecordField(fieldName, convertedFieldType)
		convertedFields = append(convertedFields, convertedField)
	}

	record := dectype.NewRecordType(info, convertedFields)

	return record, nil
}
