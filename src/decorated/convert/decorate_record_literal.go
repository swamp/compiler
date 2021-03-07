/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func findField(name string, fields []*dectype.RecordField) *dectype.RecordField {
	for _, field := range fields {
		if field.Name() == name {
			return field
		}
	}
	return nil
}

func decorateRecordLiteral(d DecorateStream, record *ast.RecordLiteral, context *VariableContext) (decorated.Expression, decshared.DecoratedError) {
	var sortedRecordAssignment []*decorated.RecordLiteralAssignment
	var recordTypeFields []*dectype.RecordField
	var decoratedRecordLiteralExpression decorated.Expression
	var decoratedRecordLiteralExpressionErr decshared.DecoratedError

	allowToAddFields := true

	var foundTemplateRecord *dectype.RecordAtom

	if record.TemplateExpression() != nil {
		allowToAddFields = false

		decoratedRecordLiteralExpression, decoratedRecordLiteralExpressionErr = DecorateExpression(d, record.TemplateExpression(), context)
		if decoratedRecordLiteralExpressionErr != nil {
			return nil, decoratedRecordLiteralExpressionErr
		}

		templateRecord, resolveErr := dectype.ResolveToRecordType(decoratedRecordLiteralExpression.Type())
		if resolveErr != nil {
			return nil, decorated.NewInternalError(resolveErr)
		}

		foundTemplateRecord = templateRecord

		recordTypeFields = append(recordTypeFields, templateRecord.SortedFields()...)
	}

	for _, assignment := range record.SortedAssignments() {
		decoratedExpression, decoratedExpressionErr := DecorateExpression(d, assignment.Expression(), context)
		if decoratedExpressionErr != nil {
			return nil, decoratedExpressionErr
		}
		encounteredFieldType := decoratedExpression.Type()
		name := assignment.Identifier().Name()
		existingField := findField(name, recordTypeFields)
		fieldExists := existingField != nil
		if !fieldExists {
			if allowToAddFields {
				fieldName := dectype.NewRecordFieldName(assignment.Identifier())
				recordTypeField := dectype.NewRecordField(fieldName, encounteredFieldType)
				recordTypeFields = append(recordTypeFields, recordTypeField)
			} else {
				return nil, decorated.NewNewRecordLiteralFieldNotInType(assignment, foundTemplateRecord)
			}
		} else {
			if compatibleErr := dectype.CompatibleTypes(encounteredFieldType, existingField.Type()); compatibleErr != nil {
				return nil, decorated.NewRecordLiteralFieldTypeMismatch(assignment, existingField, encounteredFieldType, compatibleErr)
			}
		}
	}

	recordType := dectype.NewRecordType(nil, recordTypeFields, nil) // TODO: FIX

	for _, assignment := range record.ParseOrderedAssignments() {
		decoratedExpression, decoratedExpressionErr := DecorateExpression(d, assignment.Expression(), context)
		if decoratedExpressionErr != nil {
			return nil, decoratedExpressionErr
		}
		name := assignment.Identifier().Name()
		field := recordType.FindField(name)
		recordAssignment := decorated.NewRecordLiteralAssignment(field.Index(), decorated.NewRecordLiteralField(assignment.Identifier()), decoratedExpression)
		sortedRecordAssignment = append(sortedRecordAssignment, recordAssignment)
	}

	return decorated.NewRecordLiteral(recordType, decoratedRecordLiteralExpression, sortedRecordAssignment), nil
}
