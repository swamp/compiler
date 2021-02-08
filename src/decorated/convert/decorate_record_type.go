/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func DecorateRecordType(info *ast.Record, t *dectype.TypeRepo) (*dectype.RecordAtom, dectype.DecoratedTypeError) {
	var convertedFields []*dectype.RecordField
	for _, field := range info.Fields() {
		convertedFieldType, convertedFieldTypeErr := ConvertFromAstToDecorated(field.Type(), t)
		if convertedFieldTypeErr != nil {
			return nil, convertedFieldTypeErr
		}

		concreteType, isConcreteType := convertedFieldType.(dtype.Type)
		if isConcreteType {
			convertedField := dectype.NewRecordField(field.VariableIdentifier(), concreteType)
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

	record := dectype.NewRecordType(convertedFields, convertedParameters)

	return t.DeclareRecordType(record), nil
}
