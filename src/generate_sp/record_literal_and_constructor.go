package generate_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func generateRecordConstructorSortedAssignmentsHelper(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	recordType *dectype.RecordAtom, sortedAssignments []*decorated.RecordLiteralAssignment, genContext *generateContext) error {
	if uint(target.Size) != recordType.MemorySize() {
		return fmt.Errorf("wrong target size for record constructor")
	}

	for index, assignment := range sortedAssignments {
		recordField := recordType.SortedFields()[index]
		fieldTarget := createTargetWithMemoryOffsetAndSize(target, recordField.MemoryOffset(), recordField.MemorySize())
		genErr := generateExpression(code, fieldTarget, assignment.Expression(), genContext)
		if genErr != nil {
			return genErr
		}
	}

	return nil
}

func generateRecordConstructorSortedAssignments(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	recordConstructor *decorated.RecordConstructorFromParameters, genContext *generateContext) error {
	recordType := recordConstructor.RecordType()
	return generateRecordConstructorSortedAssignmentsHelper(code, target, recordType, recordConstructor.SortedAssignments(), genContext)
}

func generateRecordLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	record *decorated.RecordLiteral, genContext *generateContext) error {
	recordType := record.RecordType()
	if record.RecordTemplate() != nil {
		if uint(target.Size) != recordType.MemorySize() {
			return fmt.Errorf("target for record literal is wrong")
		}

		structToCopyVar, genErr := generateExpressionWithSourceVar(code, record.RecordTemplate(), genContext, "gopher")
		if genErr != nil {
			return genErr
		}

		code.CopyMemory(target.Pos, structToCopyVar)

		recordFields := recordType.SortedFields()
		for index, assignment := range record.SortedAssignments() {
			recordField := recordFields[index]
			fieldTarget := createTargetWithMemoryOffsetAndSize(target, recordField.MemoryOffset(), recordField.MemorySize())
			genErr := generateExpression(code, fieldTarget, assignment.Expression(), genContext)
			if genErr != nil {
				return genErr
			}
		}
	} else {
		return generateRecordConstructorSortedAssignmentsHelper(code, target, recordType, record.SortedAssignments(), genContext)
	}
	return nil
}

func handleRecordLiteral(code *assembler_sp.Code,
	record *decorated.RecordLiteral, genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	recordType := record.RecordType()
	targetPosRange := allocMemoryForType(genContext.context.stackMemory, recordType, "record literal: "+recordType.String())
	if err := generateRecordLiteral(code, targetPosRange, record, genContext); err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}

	return targetToSourceStackPosRange(targetPosRange), nil
}
