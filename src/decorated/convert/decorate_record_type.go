/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func DecorateRecordType(info *ast.Record, t decorated.TypeAddAndReferenceMaker) (*dectype.RecordAtom, decorated.TypeError) {
	var convertedFields []*dectype.RecordField
	for _, field := range info.Fields() {
		convertedFieldType, convertedFieldTypeErr := ConvertFromAstToDecorated(field.Type(), t)
		if convertedFieldTypeErr != nil {
			return nil, convertedFieldTypeErr
		}

		concreteType, isConcreteType := convertedFieldType.(dtype.Type)
		if isConcreteType {
			fieldName := dectype.NewRecordFieldName(field.VariableIdentifier())
			convertedField := dectype.NewRecordField(fieldName, concreteType)
			convertedFields = append(convertedFields, convertedField)
		}
	}

	var convertedParameters []dtype.Type
	for _, a := range info.TypeParameters() {
		convertedParameter, convertedParameterErr := ConvertFromAstToDecorated(a, t)
		if convertedParameterErr != nil {
			return nil, convertedParameterErr
		}
		convertedParameters = append(convertedParameters, convertedParameter)
	}

	record := dectype.NewRecordType(info, convertedFields, convertedParameters)

	return record, nil
}
